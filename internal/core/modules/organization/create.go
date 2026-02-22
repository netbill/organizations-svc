package organization

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/domain"
)

type CreateParams struct {
	Name string
}

func (m *Module) Create(
	ctx context.Context,
	initiator domain.AccountActor,
	params CreateParams,
) (org domain.Organization, err error) {
	if err = m.repo.Transaction(ctx, func(ctx context.Context) error {
		org, err = m.repo.CreateOrganization(ctx, params)
		if err != nil {
			return err
		}

		err = m.messenger.WriteOrganizationCreated(ctx, org)
		if err != nil {
			return err
		}

		_, err = m.createMemberHead(ctx, initiator, org.ID)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return domain.Organization{}, err
	}

	return org, err
}

func (m *Module) createMemberHead(
	ctx context.Context,
	accountID uuid.UUID,
	organizationID uuid.UUID,
) (member domain.Member, err error) {
	member, err = m.repo.CreateMemberHead(ctx, accountID, organizationID)
	if err != nil {
		return domain.Member{}, err
	}

	err = m.messenger.WriteOrgMemberCreated(ctx, member)
	if err != nil {
		return domain.Member{}, err
	}

	return member, nil
}
