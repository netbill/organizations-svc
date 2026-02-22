package invite

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/domain"
	"github.com/netbill/organizations-svc/internal/core/errx"
)

type CreateParams struct {
	AccountID      uuid.UUID
	OrganizationID uuid.UUID
	ExpiresAt      time.Time
}

func (m *Module) Create(
	ctx context.Context,
	initiator domain.AccountActor,
	params CreateParams,
) (invite domain.Invite, err error) {
	exist, err := m.repo.MemberExists(ctx, params.AccountID, params.OrganizationID)
	if err != nil {
		return domain.Invite{}, err
	}
	if exist {
		return domain.Invite{}, errx.ErrorAccountAlreadyMember.Raise(
			fmt.Errorf("account '%s' is already a member of organization '%s'", params.AccountID, params.OrganizationID),
		)
	}

	if err = m.checkPermissionForManageInvites(ctx, initiator, params.OrganizationID); err != nil {
		return domain.Invite{}, err
	}

	if _, err = m.checkOrganizationIsActiveAndExists(ctx, params.OrganizationID); err != nil {
		return domain.Invite{}, err
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
