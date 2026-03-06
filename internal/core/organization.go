package core

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/models"
	"github.com/netbill/restkit/pagi"
)

type organizationRepo interface {
	Create(
		ctx context.Context,
		params OrganizationCreateParams,
	) (models.Organization, error)

	Get(
		ctx context.Context,
		ID uuid.UUID,
	) (models.Organization, error)
	GetListByIds(
		ctx context.Context,
		IDs []uuid.UUID,
	) ([]models.Organization, error)
	GetList(
		ctx context.Context,
		filter OrganizationFilterParams,
		limit, offset uint,
	) (pagi.Page[[]models.Organization], error)
	GetForAccountAndOrg(
		ctx context.Context,
		accountID uuid.UUID,
		limit, offset uint,
	) (pagi.Page[[]models.Organization], error)

	Update(
		ctx context.Context,
		organizationID uuid.UUID,
		params OrganizationUpdateParams,
	) (models.Organization, error)
	UpdateStatus(
		ctx context.Context,
		ID uuid.UUID,
		status string,
	) (models.Organization, error)

	Delete(ctx context.Context, ID uuid.UUID) error
}

type orgTombstoneRepo interface {
	BuryOrganization(ctx context.Context, organizationID uuid.UUID) error
	OrganizationIsBuried(ctx context.Context, organizationID uuid.UUID) (bool, error)
}

type orgMemberRepo interface {
	Create(
		ctx context.Context,
		accountID, organizationID uuid.UUID,
		head bool,
	) (models.Member, error)

	GetForAccountAndOrg(
		ctx context.Context,
		accountID models.AccountActor,
		organizationID uuid.UUID,
	) (models.Member, error)
}

type orgPlaceRepo interface {
	ExistsForOrg(ctx context.Context, organizationID uuid.UUID) (bool, error)
}

type transactor interface {
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type organizationMedia interface {
	CreateOrganizationIconUploadMediaLinks(
		ctx context.Context,
		organizationID uuid.UUID,
	) (models.UploadMediaLink, error)

	UpdateOrganizationIcon(
		ctx context.Context,
		orgID uuid.UUID,
		tempKey string,
	) (newKey string, err error)

	DeleteOrganizationIcon(
		ctx context.Context,
		organizationID uuid.UUID,
		key string,
	) error

	DeleteUploadOrganizationIcon(
		ctx context.Context,
		orgID uuid.UUID,
		key string,
	) error

	CreateOrganizationBannerUploadMediaLinks(
		ctx context.Context,
		organizationID uuid.UUID,
	) (models.UploadMediaLink, error)

	UpdateOrganizationBanner(
		ctx context.Context,
		organizationID uuid.UUID,
		key string,
	) (string, error)

	DeleteOrganizationBanner(
		ctx context.Context,
		organizationID uuid.UUID,
		key string,
	) error

	DeleteUploadOrganizationBanner(
		ctx context.Context,
		orgID uuid.UUID,
		key string,
	) error
}

type orgMessenger interface {
	WriteOrganizationCreated(ctx context.Context, organization models.Organization) error
	WriteOrganizationUpdated(ctx context.Context, organization models.Organization) error
	WriteOrganizationDeleted(ctx context.Context, organization models.Organization) error

	WriteOrgMemberCreated(ctx context.Context, member models.Member) error
}

type OrganizationModule struct {
	org       organizationRepo
	member    orgMemberRepo
	place     orgPlaceRepo
	tombstone orgTombstoneRepo
	tx        transactor
	messenger orgMessenger
	media     organizationMedia
}

type OrganizationDeps struct {
	Repo      organizationRepo
	Member    orgMemberRepo
	Place     orgPlaceRepo
	Tombstone orgTombstoneRepo
	Tx        transactor
	Messenger orgMessenger
	Media     organizationMedia
}

func NewOrganizationModule(deps OrganizationDeps) *OrganizationModule {
	return &OrganizationModule{
		org:       deps.Repo,
		member:    deps.Member,
		place:     deps.Place,
		tombstone: deps.Tombstone,
		tx:        deps.Tx,
		messenger: deps.Messenger,
		media:     deps.Media,
	}
}

type OrganizationCreateParams struct {
	Name string
}

func (m *OrganizationModule) Create(
	ctx context.Context,
	actor models.AccountActor,
	params OrganizationCreateParams,
) (org models.Organization, err error) {
	err = m.tx.Transaction(ctx, func(ctx context.Context) error {
		org, err = m.org.Create(ctx, params)
		if err != nil {
			return err
		}

		err = m.messenger.WriteOrganizationCreated(ctx, org)
		if err != nil {
			return err
		}

		member, err := m.member.Create(ctx, actor, org.ID, true)
		if err != nil {
			return err
		}

		return m.messenger.WriteOrgMemberCreated(ctx, member)
	})

	return org, err
}

func (m *OrganizationModule) GetByID(
	ctx context.Context,
	organizationID uuid.UUID,
) (models.Organization, error) {
	return m.org.Get(ctx, organizationID)
}

func (m *OrganizationModule) GetByIDs(
	ctx context.Context,
	organizationIDs []uuid.UUID,
) ([]models.Organization, error) {
	return m.org.GetListByIds(ctx, organizationIDs)
}

type OrganizationFilterParams struct {
	Text   *string
	Status *string
}

func (m *OrganizationModule) GetList(
	ctx context.Context,
	params OrganizationFilterParams,
	limit, offset uint,
) (pagi.Page[[]models.Organization], error) {
	return m.org.GetList(ctx, params, limit, offset)
}

func (m *OrganizationModule) Delete(
	ctx context.Context,
	actor models.AccountActor,
	organizationID uuid.UUID,
) error {
	_, err := m.authorizeOrgHead(ctx, actor, organizationID)
	if err != nil {
		return err
	}

	buried, err := m.tombstone.OrganizationIsBuried(ctx, organizationID)
	if err != nil {
		return err
	}
	if buried {
		return errx.ErrorOrganizationDeleted.Raise(
			fmt.Errorf("organization with id %s is already deleted", organizationID),
		)
	}

	org, err := m.validateOrg(ctx, organizationID)
	if err != nil {
		return err
	}

	places, err := m.place.ExistsForOrg(ctx, org.ID)
	if err != nil {
		return err
	}
	if places {
		return errx.ErrorOrganizationHavePlace.Raise(
			fmt.Errorf("organization %s has places, cannot be deleted", org.ID),
		)
	}

	return m.tx.Transaction(ctx, func(ctx context.Context) error {
		if err = m.tombstone.BuryOrganization(ctx, organizationID); err != nil {
			return err
		}

		if err = m.org.Delete(ctx, organizationID); err != nil {
			return err
		}

		return m.messenger.WriteOrganizationDeleted(ctx, org)
	})
}
