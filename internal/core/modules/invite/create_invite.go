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

func (s Service) CreateInvite(
	ctx context.Context,
	accountID uuid.UUID,
	params CreateParams,
) (invite models.Invite, err error) {
	initiator, err := s.getInitiator(ctx, accountID, params.OrganizationID)
	if err != nil {
		return invite, err
	}

	exist, err := s.repo.MemberExists(ctx, params.AccountID, params.OrganizationID)
	if err != nil {
		return models.Invite{}, err
	}
	if exist == true {
		return models.Invite{}, errx.ErrorAccountAlreadyMember.Raise(
			fmt.Errorf("account '%s' is already a member of organization '%s'", params.AccountID, params.OrganizationID),
		)
	}

	if err = s.checkPermissionForManageInvites(ctx, initiator.ID); err != nil {
		return models.Invite{}, err
	}

	if _, err = s.checkOrganizationIsActiveAndExists(ctx, params.OrganizationID); err != nil {
		return models.Invite{}, err
	}

	err = s.repo.Transaction(ctx, func(ctx context.Context) error {
		invite, err = s.repo.CreateInvite(ctx, params)
		if err != nil {
			return err
		}

		err = s.messenger.WriteOrgInviteCreated(ctx, invite)
		if err != nil {
			return err
		}

		return nil
	})

	return invite, err
}
