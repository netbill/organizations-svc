package invite

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/models"
)

type CreateParams struct {
	AccountID      uuid.UUID
	OrganizationID uuid.UUID
	ExpiresAt      time.Time
}

func (m *Module) CreateInvite(
	ctx context.Context,
	accountID uuid.UUID,
	params CreateParams,
) (invite models.Invite, err error) {
	initiator, err := m.getInitiator(ctx, accountID, params.OrganizationID)
	if err != nil {
		return invite, err
	}

	exist, err := m.repo.MemberExists(ctx, params.AccountID, params.OrganizationID)
	if err != nil {
		return models.Invite{}, err
	}
	if exist == true {
		return models.Invite{}, errx.ErrorAccountAlreadyMember.Raise(
			fmt.Errorf("account '%s' is already a member of organization '%s'", params.AccountID, params.OrganizationID),
		)
	}

	if err = m.checkPermissionForManageInvites(ctx, initiator.ID); err != nil {
		return models.Invite{}, err
	}

	if _, err = m.checkOrganizationIsActiveAndExists(ctx, params.OrganizationID); err != nil {
		return models.Invite{}, err
	}

	err = m.repo.Transaction(ctx, func(ctx context.Context) error {
		invite, err = m.repo.CreateInvite(ctx, params)
		if err != nil {
			return err
		}

		err = m.messenger.WriteOrgInviteCreated(ctx, invite)
		if err != nil {
			return err
		}

		return nil
	})

	return invite, err
}
