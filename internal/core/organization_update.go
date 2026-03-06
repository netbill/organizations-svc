package core

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/models"
)

type OrganizationUpdateParams struct {
	Name      *string
	IconKey   *string
	BannerKey *string
}

func (p OrganizationUpdateParams) HasChanges(current models.Organization) bool {
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

func (m *OrganizationModule) Update(
	ctx context.Context,
	actor models.AccountActor,
	organizationID uuid.UUID,
	params OrganizationUpdateParams,
) (models.Organization, error) {
	_, err := m.authorizeOrgHead(ctx, actor, organizationID)
	if err != nil {
		return models.Organization{}, err
	}

	org, err := m.validateOrg(ctx, organizationID)
	if err != nil {
		return models.Organization{}, err
	}

	if !params.HasChanges(org) {
		return org, nil
	}

	switch {
	case params.IconKey != nil && *params.IconKey == "" && org.IconKey != nil:
		if err := m.media.DeleteOrganizationIcon(ctx, org.ID, *org.IconKey); err != nil {
			return models.Organization{}, err
		}
		params.IconKey = nil
	case params.IconKey != nil:
		iconKey, err := m.media.UpdateOrganizationIcon(ctx, org.ID, *params.IconKey)
		if err != nil {
			return models.Organization{}, err
		}

		params.IconKey = &iconKey
	}

	switch {
	case params.BannerKey != nil && *params.BannerKey == "" && org.BannerKey != nil:
		if err := m.media.DeleteOrganizationBanner(ctx, org.ID, *org.BannerKey); err != nil {
			return models.Organization{}, err
		}
		params.BannerKey = nil
	case params.BannerKey != nil:
		bannerKey, err := m.media.UpdateOrganizationBanner(ctx, org.ID, *params.BannerKey)
		if err != nil {
			return models.Organization{}, err
		}

		params.BannerKey = &bannerKey
	}

	err = m.tx.Transaction(ctx, func(ctx context.Context) error {
		org, err = m.org.Update(ctx, organizationID, params)
		if err != nil {
			return err
		}

		return m.messenger.WriteOrganizationUpdated(ctx, org)
	})

	return org, err
}

func (m *OrganizationModule) Activate(
	ctx context.Context,
	actor models.AccountActor,
	organizationID uuid.UUID,
	value bool,
) (models.Organization, error) {
	_, err := m.authorizeOrgHead(ctx, actor, organizationID)
	if err != nil {
		return models.Organization{}, err
	}

	org, err := m.validateOrg(ctx, organizationID)
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
		org, err = m.org.UpdateStatus(ctx, organizationID, status)
		if err != nil {
			return err
		}

		return m.messenger.WriteOrganizationUpdated(ctx, org)
	})

	return org, err
}
