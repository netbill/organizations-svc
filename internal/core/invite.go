package core

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/models"
	"github.com/netbill/restkit/pagi"
)

type inviteRepo interface {
	Create(ctx context.Context, params InviteCreateParams) (models.Invite, error)

	Get(ctx context.Context, inviteID uuid.UUID) (models.Invite, error)
	GetList(
		ctx context.Context,
		params FilterInvitesParams,
		limit, offset uint,
	) (pagi.Page[[]models.Invite], error)
	ExistActiveForAccountInOrg(ctx context.Context, accountID, organizationID uuid.UUID) (bool, error)

	UpdateStatus(ctx context.Context, inviteID uuid.UUID, status string) (models.Invite, error)

	Delete(ctx context.Context, inviteID uuid.UUID) error
}

type inviteMemberRepo interface {
	Create(
		ctx context.Context,
		accountID, organizationID uuid.UUID,
		head bool,
	) (models.Member, error)

	ExistsForAccountAndOrg(ctx context.Context, accountID, organizationID uuid.UUID) (bool, error)
}

type inviteProfileRepo interface {
	ExistsByAccountID(ctx context.Context, accountID uuid.UUID) (bool, error)
}

type inviteMessenger interface {
	WriteOrgMemberCreated(ctx context.Context, member models.Member) error

	WriteOrgInviteCreated(ctx context.Context, invite models.Invite) error
	WriteOrgInviteAccepted(ctx context.Context, invite models.Invite) error
	WriteOrgInviteDeclined(ctx context.Context, invite models.Invite) error
	WriteOrgInviteCanceled(ctx context.Context, invite models.Invite) error
}

type InviteModule struct {
	auth      orgAuth
	repo      inviteRepo
	member    inviteMemberRepo
	profile   inviteProfileRepo
	tx        transactor
	messenger inviteMessenger
}

type InviteDeps struct {
	Auth      orgAuth
	Repo      inviteRepo
	Member    inviteMemberRepo
	Profile   inviteProfileRepo
	Tx        transactor
	Messenger inviteMessenger
}

func NewInviteModule(deps InviteDeps) *InviteModule {
	return &InviteModule{
		auth:      deps.Auth,
		repo:      deps.Repo,
		member:    deps.Member,
		profile:   deps.Profile,
		tx:        deps.Tx,
		messenger: deps.Messenger,
	}
}

type InviteCreateParams struct {
	AccountID      uuid.UUID
	OrganizationID uuid.UUID
	ExpiresAt      time.Time
}

func (m *InviteModule) Create(
	ctx context.Context,
	actor models.AccountActor,
	params InviteCreateParams,
) (invite models.Invite, err error) {
	_, err = m.auth.authorizeOrgHead(ctx, actor, params.OrganizationID)
	if err != nil {
		return models.Invite{}, err
	}

	_, err = m.auth.validateOrg(ctx, params.OrganizationID)
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

type FilterInvitesParams struct {
	OrganizationID *uuid.UUID
	AccountID      *uuid.UUID
}

func (m *InviteModule) GetList(
	ctx context.Context,
	actor models.AccountActor,
	params FilterInvitesParams,
	limit, offset uint,
) (pagi.Page[[]models.Invite], error) {
	if params.OrganizationID != nil {
		_, err := m.auth.authorizeOrgMember(ctx, actor, *params.OrganizationID)
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

func (m *InviteModule) GetForAccount(
	ctx context.Context,
	accountID uuid.UUID,
	inviteID uuid.UUID,
) (models.Invite, error) {
	invite, err := m.repo.Get(ctx, inviteID)
	if err != nil {
		return models.Invite{}, err
	}

	if invite.AccountID != accountID {
		_, err = m.auth.authorizeOrgMember(ctx, accountID, invite.OrganizationID)
		if err != nil {
			return models.Invite{}, err
		}
	}

	return invite, nil
}
