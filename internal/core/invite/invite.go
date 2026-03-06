package invite

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/errx"
	"github.com/netbill/organizations-svc/internal/models"
	"github.com/netbill/restkit/pagi"
)

type messenger interface {
	WriteOrgMemberCreated(ctx context.Context, member models.Member) error

	WriteOrgInviteCreated(ctx context.Context, invite models.Invite) error
	WriteOrgInviteAccepted(ctx context.Context, invite models.Invite) error
	WriteOrgInviteDeclined(ctx context.Context, invite models.Invite) error
	WriteOrgInviteCanceled(ctx context.Context, invite models.Invite) error
}

type Service struct {
	auth      org
	repo      repo
	member    member
	profile   profile
	tx        transactor
	messenger messenger
}

type InviteDeps struct {
	Auth      org
	Repo      repo
	Member    member
	Profile   profile
	Tx        transactor
	Messenger messenger
}

func NewInviteModule(deps InviteDeps) *Service {
	return &Service{
		auth:      deps.Auth,
		repo:      deps.Repo,
		member:    deps.Member,
		profile:   deps.Profile,
		tx:        deps.Tx,
		messenger: deps.Messenger,
	}
}

type CreateParams struct {
	AccountID      uuid.UUID
	OrganizationID uuid.UUID
	ExpiresAt      time.Time
}

func (m *Service) Create(
	ctx context.Context,
	actor models.AccountActor,
	params CreateParams,
) (invite models.Invite, err error) {
	_, err = m.auth.AuthorizeOrgHead(ctx, actor, params.OrganizationID)
	if err != nil {
		return models.Invite{}, err
	}

	_, err = m.auth.ValidateOrg(ctx, params.OrganizationID)
	if err != nil {
		return models.Invite{}, err
	}

	exist, err := m.profile.ExistsByAccountID(ctx, params.AccountID)
	if err != nil {
		return models.Invite{}, err
	}
	if !exist {
		return models.Invite{}, errx.ErrorProfileNotExists.Raise(
			fmt.Errorf("profile for '%s' not found", params.AccountID),
		)
	}

	memberExists, err := m.member.ExistsForAccountAndOrg(ctx, params.AccountID, params.OrganizationID)
	if err != nil {
		return models.Invite{}, err
	}
	if memberExists {
		return models.Invite{}, errx.ErrorAccountAlreadyMember.Raise(
			fmt.Errorf("account with id %s is already a member of organization %s",
				params.AccountID, params.OrganizationID,
			),
		)
	}

	exist, err = m.repo.ExistActiveForAccountInOrg(ctx, params.AccountID, params.OrganizationID)
	if err != nil {
		return models.Invite{}, err
	}
	if exist {
		return models.Invite{}, errx.ErrorActiveInviteAlreadyExists.Raise(
			fmt.Errorf("active invite for account %s already exists in organization %s",
				params.AccountID, params.OrganizationID,
			),
		)
	}

	err = m.tx.Transaction(ctx, func(ctx context.Context) error {
		invite, err = m.repo.Create(ctx, params)
		if err != nil {
			return err
		}

		return m.messenger.WriteOrgInviteCreated(ctx, invite)
	})

	return invite, err
}

type FilterParams struct {
	OrganizationID *uuid.UUID
	AccountID      *uuid.UUID
}

func (m *Service) GetList(
	ctx context.Context,
	actor models.AccountActor,
	params FilterParams,
	limit, offset uint,
) (pagi.Page[[]models.Invite], error) {
	if params.OrganizationID != nil {
		_, err := m.auth.AuthorizeOrgMember(ctx, actor, *params.OrganizationID)
		if err != nil {
			return pagi.Page[[]models.Invite]{}, err
		}
	}
	if params.AccountID != nil && params.OrganizationID == nil {
		if actor != *params.AccountID {
			return pagi.Page[[]models.Invite]{}, errx.ErrorCannotGetInvitesForOtherAccount.Raise(
				fmt.Errorf("cannot get invites for other account"),
			)
		}
	}
	if params.OrganizationID == nil && params.AccountID == nil {
		params.AccountID = &actor
	}

	return m.repo.GetList(ctx, params, limit, offset)
}

func (m *Service) GetForAccount(
	ctx context.Context,
	accountID uuid.UUID,
	inviteID uuid.UUID,
) (models.Invite, error) {
	invite, err := m.repo.Get(ctx, inviteID)
	if err != nil {
		return models.Invite{}, err
	}

	if invite.AccountID != accountID {
		_, err = m.auth.AuthorizeOrgMember(ctx, accountID, invite.OrganizationID)
		if err != nil {
			return models.Invite{}, err
		}
	}

	return invite, nil
}
