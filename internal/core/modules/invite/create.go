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

func (m *Module) Create(
	ctx context.Context,
	actor models.AccountActor,
	params CreateParams,
) (invite models.Invite, err error) {
	initiator, err := m.getInitiator(ctx, actor, params.OrganizationID)
	if err != nil {
		return models.Invite{}, err
	}
	if !initiator.Head {
		return models.Invite{}, errx.ErrorNotEnoughRights.Raise(
			fmt.Errorf("account has no rights to create invite for this organization"),
		)
	}

	exist, err := m.repo.ExistsProfileByAccountID(ctx, params.AccountID)
	if err != nil {
		return models.Invite{}, err
	}
	if !exist {
		return models.Invite{}, errx.ErrorProfileNotFound.Raise(
			fmt.Errorf("profile for '%s' not found", params.AccountID),
		)
	}

	org, err := m.repo.GetOrganizationByID(ctx, params.OrganizationID)
	if err != nil {
		return models.Invite{}, err
	}
	if org.Status != models.OrganizationStatusActive {
		return models.Invite{}, errx.ErrorOrganizationIsNotActive.Raise(
			fmt.Errorf("organization with id %s is not active", params.OrganizationID),
		)
	}

	member, err := m.repo.MemberExists(ctx, params.AccountID, params.OrganizationID)
	if err != nil {
		return models.Invite{}, err
	}
	if member {
		return models.Invite{}, errx.ErrorAccountAlreadyMember.Raise(
			fmt.Errorf("account with id %s is already a member of organization %s", params.AccountID, params.OrganizationID),
		)
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
