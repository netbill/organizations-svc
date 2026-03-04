package core

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/models"
	"github.com/netbill/restkit/pagi"
)

type orgBucket interface {
	CreateOrganizationIconUploadMediaLinks(ctx context.Context, orgID uuid.UUID) (models.UploadMediaLink, error)
	CreateOrganizationBannerUploadMediaLinks(ctx context.Context, orgID uuid.UUID) (models.UploadMediaLink, error)

	UpdateOrganizationIcon(
		ctx context.Context,
		orgID uuid.UUID,
		oldKey *string,
		tempKey *string,
	) (newKey *string, err error)

	UpdateOrganizationBanner(
		ctx context.Context,
		orgID uuid.UUID,
		oldKey *string,
		tempKey *string,
	) (newKey *string, err error)

	DeleteUploadOrganizationIcon(ctx context.Context, orgID uuid.UUID, key string) error
	DeleteUploadOrganizationBanner(ctx context.Context, orgID uuid.UUID, key string) error
}

type orgRepository interface {
	CreateOrganization(ctx context.Context, params OrganizationCreateParams) (models.Organization, error)

	UpdateOrganization(ctx context.Context, ID uuid.UUID, params OrganizationUpdateParams) (models.Organization, error)
	UpdateOrganizationStatus(ctx context.Context, ID uuid.UUID, status string) (models.Organization, error)

	GetOrganizationByID(ctx context.Context, ID uuid.UUID) (models.Organization, error)
	GetOrganizationsByIDs(ctx context.Context, IDs []uuid.UUID) ([]models.Organization, error)
	GetOrganizations(ctx context.Context, filter OrganizationFilterParams, limit, offset uint) (pagi.Page[[]models.Organization], error)
	GetOrganizationsForUser(ctx context.Context, accountID uuid.UUID, limit, offset uint) (pagi.Page[[]models.Organization], error)

	DeleteOrganization(ctx context.Context, ID uuid.UUID) error
}

type orgMemberRepository interface {
	CreateMemberHead(ctx context.Context, accountID, organizationID uuid.UUID) (models.Member, error)
	GetMemberByAccountAndOrganization(ctx context.Context, accountID, organizationID uuid.UUID) (models.Member, error)
}

type orgPlaceRepository interface {
	GetPlaceExistsForOrganization(ctx context.Context, organizationID uuid.UUID) (bool, error)
}

type orgTombstoneRepository interface {
	BuryOrganization(ctx context.Context, organizationID uuid.UUID) error
	OrganizationIsBuried(ctx context.Context, organizationID uuid.UUID) (bool, error)
}

type transactor interface {
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type orgMessenger interface {
	WriteOrganizationCreated(ctx context.Context, organization models.Organization) error
	WriteOrganizationUpdated(ctx context.Context, organization models.Organization) error
	WriteOrganizationDeleted(ctx context.Context, organization models.Organization) error

	WriteOrgMemberCreated(ctx context.Context, member models.Member) error
}

type OrganizationModule struct {
	orgRepo       orgRepository
	memberRepo    orgMemberRepository
	placeRepo     orgPlaceRepository
	tombstoneRepo orgTombstoneRepository
	tx            transactor
	messenger     orgMessenger
	bucket        orgBucket
}

func (m *OrganizationModule) authorizeOrgHead(
	ctx context.Context,
	actor models.AccountActor,
	organizationID uuid.UUID,
) (models.Organization, models.Member, error) {
	org, err := m.GetByID(ctx, organizationID)
	if err != nil {
		return models.Organization{}, models.Member{}, err
	}

	if org.Status == models.OrganizationStatusSuspended {
		return models.Organization{}, models.Member{}, errx.ErrorOrganizationIsSuspended.Raise(
			fmt.Errorf("organization with id %s is suspended", organizationID),
		)
	}

	member, err := m.memberRepo.GetMemberByAccountAndOrganization(ctx, actor, organizationID)
	if errors.Is(err, errx.ErrorMemberNotExists) {
		return models.Organization{}, models.Member{}, errx.ErrorInitiatorNotMemberOfOrganization.Raise(
			fmt.Errorf("initiator with account id %s is not a member of organization %s", actor, organizationID),
		)
	}
	if err != nil {
		return models.Organization{}, models.Member{}, err
	}

	if !member.Head {
		return models.Organization{}, models.Member{}, errx.ErrorNotOrganizationHead.Raise(
			fmt.Errorf(
				"only organization head member can manage organization, but member %s is not head", member.ID,
			),
		)
	}

	return org, member, nil
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
		org, err = m.orgRepo.CreateOrganization(ctx, params)
		if err != nil {
			return err
		}

		if err = m.messenger.WriteOrganizationCreated(ctx, org); err != nil {
			return err
		}

		member, err := m.memberRepo.CreateMemberHead(ctx, actor, org.ID)
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
	org, err := m.orgRepo.GetOrganizationByID(ctx, organizationID)
	if errors.Is(err, errx.ErrorOrganizationNotExists) {
		buried, err := m.tombstoneRepo.OrganizationIsBuried(ctx, organizationID)
		if err != nil {
			return models.Organization{}, fmt.Errorf("check organization tombstone: %w", err)
		}

		if buried {
			return models.Organization{}, errx.ErrorOrganizationNotExists.Raise(
				fmt.Errorf("organization with id %s is buried", organizationID),
			)
		}
	}

	return org, err
}

func (m *OrganizationModule) GetByIDs(
	ctx context.Context,
	organizationIDs []uuid.UUID,
) ([]models.Organization, error) {
	return m.orgRepo.GetOrganizationsByIDs(ctx, organizationIDs)
}

type OrganizationFilterParams struct {
	Text   *string `json:"name,omitempty"`
	Status *string `json:"status,omitempty"`
}

func (m *OrganizationModule) GetList(
	ctx context.Context,
	params OrganizationFilterParams,
	limit, offset uint,
) (pagi.Page[[]models.Organization], error) {
	return m.orgRepo.GetOrganizations(ctx, params, limit, offset)
}

func (m *OrganizationModule) GetForUser(
	ctx context.Context,
	actor models.AccountActor,
	limit, offset uint,
) (pagi.Page[[]models.Organization], error) {
	return m.orgRepo.GetOrganizationsForUser(ctx, actor, limit, offset)
}

type OrganizationUpdateParams struct {
	Name      string
	IconKey   *string
	BannerKey *string
}

func (m *OrganizationUpdateParams) isEmpty(model models.Organization) (upd bool) {
	if m.Name != model.Name {
		return false
	}

	if m.IconKey != nil && (model.IconKey == nil || *m.IconKey != *model.IconKey) {
		return false
	}

	if m.BannerKey != nil && (model.BannerKey == nil || *m.BannerKey != *model.BannerKey) {
		return false
	}

	return true
}

func (m *OrganizationModule) Update(
	ctx context.Context,
	actor models.AccountActor,
	organizationID uuid.UUID,
	params OrganizationUpdateParams,
) (models.Organization, error) {
	org, _, err := m.authorizeOrgHead(ctx, actor, organizationID)
	if err != nil {
		return models.Organization{}, err
	}

	if params.isEmpty(org) {
		return org, nil
	}

	if params.IconKey != nil {
		iconKey, err := m.bucket.UpdateOrganizationIcon(ctx, org.ID, org.IconKey, params.IconKey)
		if err != nil {
			return models.Organization{}, err
		}
		params.IconKey = iconKey
	}

	if params.BannerKey != nil {
		bannerKey, err := m.bucket.UpdateOrganizationBanner(ctx, org.ID, org.BannerKey, params.BannerKey)
		if err != nil {
			return models.Organization{}, err
		}
		params.BannerKey = bannerKey
	}

	err = m.tx.Transaction(ctx, func(ctx context.Context) error {
		org, err = m.orgRepo.UpdateOrganization(ctx, organizationID, params)
		if err != nil {
			return err
		}

		return m.messenger.WriteOrganizationUpdated(ctx, org)
	})

	return org, err
}

func (m *OrganizationModule) Delete(
	ctx context.Context,
	actor models.AccountActor,
	organizationID uuid.UUID,
) error {
	org, _, err := m.authorizeOrgHead(ctx, actor, organizationID)
	if err != nil {
		return err
	}

	places, err := m.placeRepo.GetPlaceExistsForOrganization(ctx, org.ID)
	if err != nil {
		return err
	}
	if places {
		return errx.ErrorOrganizationHavePlace.Raise(
			fmt.Errorf("organization %s has places, cannot be deleted", org.ID),
		)
	}

	return m.tx.Transaction(ctx, func(ctx context.Context) error {
		if err = m.tombstoneRepo.BuryOrganization(ctx, organizationID); err != nil {
			return err
		}

		if err = m.orgRepo.DeleteOrganization(ctx, organizationID); err != nil {
			return err
		}

		return m.messenger.WriteOrganizationDeleted(ctx, org)
	})
}

func (m *OrganizationModule) Activate(
	ctx context.Context,
	actor models.AccountActor,
	organizationID uuid.UUID,
	value bool,
) (models.Organization, error) {
	org, _, err := m.authorizeOrgHead(ctx, actor, organizationID)
	if err != nil {
		return models.Organization{}, err
	}

	if org.Status == models.OrganizationStatusActive && value {
		return org, nil
	}
	if org.Status != models.OrganizationStatusActive && !value {
		return org, nil
	}

	var newStatus string
	if value {
		newStatus = models.OrganizationStatusActive
	} else {
		newStatus = models.OrganizationStatusInactive
	}

	return m.updateStatus(ctx, organizationID, newStatus)
}

func (m *OrganizationModule) Suspend(
	ctx context.Context,
	organizationID uuid.UUID,
	value bool,
) (models.Organization, error) {
	org, err := m.GetByID(ctx, organizationID)
	if err != nil {
		return models.Organization{}, err
	}

	if org.Status == models.OrganizationStatusSuspended && value {
		return org, nil
	}
	if org.Status != models.OrganizationStatusSuspended && !value {
		return org, nil
	}

	var newStatus string
	if value {
		newStatus = models.OrganizationStatusSuspended
	} else {
		newStatus = models.OrganizationStatusInactive
	}

	return m.updateStatus(ctx, organizationID, newStatus)
}

func (m *OrganizationModule) updateStatus(
	ctx context.Context,
	organizationID uuid.UUID,
	status string,
) (org models.Organization, err error) {
	err = m.tx.Transaction(ctx, func(ctx context.Context) error {
		org, err = m.orgRepo.UpdateOrganizationStatus(ctx, organizationID, status)
		if err != nil {
			return err
		}

		return m.messenger.WriteOrganizationUpdated(ctx, org)
	})

	return org, err
}

func (m *OrganizationModule) CreateOrgUploadMediaLinks(
	ctx context.Context,
	actor models.AccountActor,
	organizationID uuid.UUID,
) (models.Organization, models.UploadOrgMediaLinks, error) {
	org, _, err := m.authorizeOrgHead(ctx, actor, organizationID)
	if err != nil {
		return models.Organization{}, models.UploadOrgMediaLinks{}, err
	}

	iconLinks, err := m.bucket.CreateOrganizationIconUploadMediaLinks(ctx, organizationID)
	if err != nil {
		return models.Organization{}, models.UploadOrgMediaLinks{}, err
	}

	bannerLinks, err := m.bucket.CreateOrganizationBannerUploadMediaLinks(ctx, organizationID)
	if err != nil {
		return models.Organization{}, models.UploadOrgMediaLinks{}, err
	}

	return org, models.UploadOrgMediaLinks{
		Icon:   iconLinks,
		Banner: bannerLinks,
	}, nil
}

func (m *OrganizationModule) DeleteOrgUploadIcon(
	ctx context.Context,
	actor models.AccountActor,
	organizationID uuid.UUID,
	key string,
) error {
	_, _, err := m.authorizeOrgHead(ctx, actor, organizationID)
	if err != nil {
		return err
	}

	return m.bucket.DeleteUploadOrganizationIcon(ctx, organizationID, key)
}

func (m *OrganizationModule) DeleteOrgUploadBanner(
	ctx context.Context,
	actor models.AccountActor,
	organizationID uuid.UUID,
	key string,
) error {
	_, _, err := m.authorizeOrgHead(ctx, actor, organizationID)
	if err != nil {
		return err
	}

	return m.bucket.DeleteUploadOrganizationBanner(ctx, organizationID, key)
}
