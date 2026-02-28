package organization

import (
	"context"

	"github.com/netbill/organizations-svc/internal/core/models"
)

type CreateParams struct {
	Name string
}

func (m *Module) Create(
	ctx context.Context,
	actor models.AccountActor,
	params CreateParams,
) (org models.Organization, err error) {
	if err = m.repo.Transaction(ctx, func(ctx context.Context) error {
		org, err = m.repo.CreateOrganization(ctx, params)
		if err != nil {
			return err
		}

		if err = m.messenger.WriteOrganizationCreated(ctx, org); err != nil {
			return err
		}

		member, err := m.repo.CreateMemberHead(ctx, actor, org.ID)
		if err != nil {
			return err
		}

		return m.messenger.WriteOrgMemberCreated(ctx, member)
	}); err != nil {
		return models.Organization{}, err
	}

	return org, err
}
