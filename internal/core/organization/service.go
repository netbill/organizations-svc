package organization

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/errx"
	"github.com/netbill/organizations-svc/internal/models"
	"github.com/netbill/restkit/pagi"
)

type messenger interface {
	WriteOrganizationCreated(ctx context.Context, organization models.Organization) error
	WriteOrganizationUpdated(ctx context.Context, organization models.Organization) error
	WriteOrganizationDeleted(ctx context.Context, organization models.Organization) error

	WriteOrgMemberCreated(ctx context.Context, member models.Member) error
}

type Service struct {
	repo      repo
	member    memberRepo
	place     placeRepo
	tombstone tombstone
	tx        transactor
	messenger messenger
	media     media
}

type ServiceDeps struct {
	Repo      repo
	Member    memberRepo
	Place     placeRepo
	Tombstone tombstone
	Tx        transactor
	Messenger messenger
	Media     media
}

func NewOrganizationModule(deps ServiceDeps) *Service {
	return &Service{
		repo:      deps.Repo,
		member:    deps.Member,
		place:     deps.Place,
		tombstone: deps.Tombstone,
		tx:        deps.Tx,
		messenger: deps.Messenger,
		media:     deps.Media,
	}
}

type CreateParams struct {
	Name string
}

func (s *Service) Create(
	ctx context.Context,
	actor models.AccountActor,
	params CreateParams,
) (org models.Organization, err error) {
	err = s.tx.Transaction(ctx, func(ctx context.Context) error {
		org, err = s.repo.Create(ctx, params)
		if err != nil {
			return err
		}

		err = s.messenger.WriteOrganizationCreated(ctx, org)
		if err != nil {
			return err
		}

		member, err := s.member.Create(ctx, actor, org.ID, true)
		if err != nil {
			return err
		}

		return s.messenger.WriteOrgMemberCreated(ctx, member)
	})

	return org, err
}

func (s *Service) GetByID(
	ctx context.Context,
	organizationID uuid.UUID,
) (models.Organization, error) {
	org, err := s.repo.Get(ctx, organizationID)
	if errors.Is(err, errx.ErrorOrganizationNotExists) {
		buried, err := s.tombstone.OrganizationIsBuried(ctx, organizationID)
		if err != nil {
			return models.Organization{}, err
		}
		if buried {
			return models.Organization{}, errx.ErrorOrganizationDeleted.Raise(
				fmt.Errorf("organization with id %s is deleted", organizationID),
			)
		}
	}
	if err != nil {
		return models.Organization{}, err
	}

	return org, nil
}

func (s *Service) GetByIDs(
	ctx context.Context,
	organizationIDs []uuid.UUID,
) ([]models.Organization, error) {
	return s.repo.GetListByIds(ctx, organizationIDs)
}

// TODO make options for filtering, sorting and pagination
type FilterParams struct {
	Text   *string
	Status *string
}

func (s *Service) GetList(
	ctx context.Context,
	params FilterParams,
	limit, offset uint,
) (pagi.Page[[]models.Organization], error) {
	return s.repo.GetList(ctx, params, limit, offset)
}

func (s *Service) Delete(
	ctx context.Context,
	actor models.AccountActor,
	organizationID uuid.UUID,
) error {
	_, err := s.AuthorizeOrgHead(ctx, actor, organizationID)
	if err != nil {
		return err
	}

	buried, err := s.tombstone.OrganizationIsBuried(ctx, organizationID)
	if err != nil {
		return err
	}
	if buried {
		return errx.ErrorOrganizationDeleted.Raise(
			fmt.Errorf("organization with id %s is already deleted", organizationID),
		)
	}

	org, err := s.ValidateOrg(ctx, organizationID)
	if err != nil {
		return err
	}

	places, err := s.place.ExistsForOrg(ctx, org.ID)
	if err != nil {
		return err
	}
	if places {
		return errx.ErrorOrganizationHavePlace.Raise(
			fmt.Errorf("organization %s has places, cannot be deleted", org.ID),
		)
	}

	return s.tx.Transaction(ctx, func(ctx context.Context) error {
		if err = s.tombstone.BuryOrganization(ctx, organizationID); err != nil {
			return err
		}

		if err = s.repo.Delete(ctx, organizationID); err != nil {
			return err
		}

		return s.messenger.WriteOrganizationDeleted(ctx, org)
	})
}
