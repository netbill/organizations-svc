package organization

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/models"
)

type CreateParams struct {
	Name string
}

func (m *Module) Create(
	ctx context.Context,
	initiator models.AccountActor,
	params CreateParams,
) (org models.Organization, err error) {
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
		return models.Organization{}, err
	}

	return org, err
}

func (m *Module) createMemberHead(
	ctx context.Context,
	accountID uuid.UUID,
	organizationID uuid.UUID,
) (member models.Member, err error) {
	member, err = m.repo.CreateMember(ctx, accountID, organizationID)
	if err != nil {
		return models.Member{}, err
	}

	err = m.messenger.WriteOrgMemberCreated(ctx, member)
	if err != nil {
		return models.Member{}, err
	}

	return member, nil
}
