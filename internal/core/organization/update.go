package organization

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/models"
)

type UpdateParams struct {
	Name      *string
	IconKey   *string
	BannerKey *string
}

func (p UpdateParams) HasChanges(current models.Organization) bool {
	if p.Name != nil && *p.Name != current.Name {
		return true
	}
	if p.IconKey != nil && !ptrEqual(p.IconKey, current.IconKey) {
		return true
	}
	if p.BannerKey != nil && !ptrEqual(p.BannerKey, current.BannerKey) {
		return true
	}
	return false
}

func ptrEqual[T comparable](a, b *T) bool {
	if a == nil || b == nil {
		return a == b
	}
	return *a == *b
}

func (s *Service) Update(
	ctx context.Context,
	actor models.AccountActor,
	organizationID uuid.UUID,
	params UpdateParams,
) (models.Organization, error) {
	_, err := s.AuthorizeOrgHead(ctx, actor, organizationID)
	if err != nil {
		return models.Organization{}, err
	}

	org, err := s.ValidateOrg(ctx, organizationID)
	if err != nil {
		return models.Organization{}, err
	}

	if !params.HasChanges(org) {
		return org, nil
	}

	switch {
	case params.IconKey != nil && *params.IconKey == "" && org.IconKey != nil:
		if err := s.media.DeleteOrganizationIcon(ctx, org.ID, *org.IconKey); err != nil {
			return models.Organization{}, err
		}
		params.IconKey = nil
	case params.IconKey != nil:
		iconKey, err := s.media.UpdateOrganizationIcon(ctx, org.ID, *params.IconKey)
		if err != nil {
			return models.Organization{}, err
		}

		params.IconKey = &iconKey
	}

	switch {
	case params.BannerKey != nil && *params.BannerKey == "" && org.BannerKey != nil:
		if err := s.media.DeleteOrganizationBanner(ctx, org.ID, *org.BannerKey); err != nil {
			return models.Organization{}, err
		}
		params.BannerKey = nil
	case params.BannerKey != nil:
		bannerKey, err := s.media.UpdateOrganizationBanner(ctx, org.ID, *params.BannerKey)
		if err != nil {
			return models.Organization{}, err
		}

		params.BannerKey = &bannerKey
	}

	err = s.tx.Transaction(ctx, func(ctx context.Context) error {
		org, err = s.repo.Update(ctx, organizationID, params)
		if err != nil {
			return err
		}

		return s.messenger.WriteOrganizationUpdated(ctx, org)
	})

	return org, err
}

func (s *Service) Activate(
	ctx context.Context,
	actor models.AccountActor,
	organizationID uuid.UUID,
	value bool,
) (models.Organization, error) {
	_, err := s.AuthorizeOrgHead(ctx, actor, organizationID)
	if err != nil {
		return models.Organization{}, err
	}

	org, err := s.ValidateOrg(ctx, organizationID)
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

	return s.updateStatus(ctx, organizationID, newStatus)
}

func (s *Service) Suspend(
	ctx context.Context,
	organizationID uuid.UUID,
	value bool,
) (models.Organization, error) {
	org, err := s.GetByID(ctx, organizationID)
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

	return s.updateStatus(ctx, organizationID, newStatus)
}

func (s *Service) updateStatus(
	ctx context.Context,
	organizationID uuid.UUID,
	status string,
) (org models.Organization, err error) {
	err = s.tx.Transaction(ctx, func(ctx context.Context) error {
		org, err = s.repo.UpdateStatus(ctx, organizationID, status)
		if err != nil {
			return err
		}

		return s.messenger.WriteOrganizationUpdated(ctx, org)
	})

	return org, err
}
