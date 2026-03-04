package core

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/models"
	"github.com/netbill/restkit/pagi"
)

type inviteRepository interface {
	CreateInvite(ctx context.Context, params InviteCreateParams) (models.Invite, error)
	GetInvite(ctx context.Context, inviteID uuid.UUID) (models.Invite, error)
	GetOrganizationInvites(ctx context.Context, organizationID uuid.UUID, limit, offset uint) (pagi.Page[[]models.Invite], error)
	GetAccountInvites(ctx context.Context, accountID uuid.UUID, limit, offset uint) (pagi.Page[[]models.Invite], error)
	ExistsActiveInviteByAccountID(ctx context.Context, accountID, organizationID uuid.UUID) (bool, error)
	UpdateInviteStatus(ctx context.Context, inviteID uuid.UUID, status string) (models.Invite, error)
	DeleteInvite(ctx context.Context, inviteID uuid.UUID) error
}

type inviteMemberRepo interface {
	GetMemberByAccountAndOrganization(ctx context.Context, accountID, organizationID uuid.UUID) (models.Member, error)
	MemberExists(ctx context.Context, accountID, organizationID uuid.UUID) (bool, error)
	CreateMember(ctx context.Context, accountID, organizationID uuid.UUID) (models.Member, error)
}

type inviteOrgRepo interface {
	GetOrganizationByID(ctx context.Context, ID uuid.UUID) (models.Organization, error)
}

type inviteProfileRepo interface {
	ExistsProfileByAccountID(ctx context.Context, accountID uuid.UUID) (bool, error)
}

type inviteMessenger interface {
	WriteOrgMemberCreated(ctx context.Context, member models.Member) error

	WriteOrgInviteCreated(ctx context.Context, invite models.Invite) error
	WriteOrgInviteAccepted(ctx context.Context, invite models.Invite) error
	WriteOrgInviteDeclined(ctx context.Context, invite models.Invite) error
	WriteOrgInviteCanceled(ctx context.Context, invite models.Invite) error
}

type InviteModule struct {
	inviteRepo  inviteRepository
	memberRepo  inviteMemberRepo
	orgRepo     inviteOrgRepo
	profileRepo inviteProfileRepo
	tx          transactor
	messenger   inviteMessenger
}

func (m *InviteModule) authorizeOrgHead(
	ctx context.Context,
	actor models.AccountActor,
	organizationID uuid.UUID,
) (models.Organization, models.Member, error) {
	org, err := m.orgRepo.GetOrganizationByID(ctx, organizationID)
	if err != nil {
		return models.Organization{}, models.Member{}, err
	}

	if org.Status == models.OrganizationStatusSuspended {
		return models.Organization{}, models.Member{}, errx.ErrorOrganizationIsSuspended.Raise(
			fmt.Errorf("organization with id %s is suspended", organizationID),
		)
	}

	member, err := m.getInitiator(ctx, actor, organizationID)
	if err != nil {
		return models.Organization{}, models.Member{}, err
	}

	if !member.Head {
		return models.Organization{}, models.Member{}, errx.ErrorNotOrganizationHead.Raise(
			fmt.Errorf(
				"only organization head member can manage members, but member %s is not head", member.ID,
			),
		)
	}

	return org, member, nil
}

func (m *InviteModule) getInitiator(
	ctx context.Context,
	accountID uuid.UUID,
	organizationID uuid.UUID,
) (models.Member, error) {
	initiator, err := m.memberRepo.GetMemberByAccountAndOrganization(ctx, accountID, organizationID)
	if errors.Is(err, errx.ErrorMemberNotExists) {
		return models.Member{}, errx.ErrorInitiatorNotMemberOfOrganization.Raise(
			fmt.Errorf("initiator with account id %s is not a member of organization %s", accountID, organizationID),
		)
	}
	if err != nil {
		return models.Member{}, err
	}

	return initiator, nil
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
	_, _, err = m.authorizeOrgHead(ctx, actor, params.OrganizationID)
	if err != nil {
		return models.Invite{}, err
	}

	exist, err := m.profileRepo.ExistsProfileByAccountID(ctx, params.AccountID)
	if err != nil {
		return models.Invite{}, err
	}
	if !exist {
		return models.Invite{}, errx.ErrorProfileNotExists.Raise(
			fmt.Errorf("profile for '%s' not found", params.AccountID),
		)
	}

	memberExists, err := m.memberRepo.MemberExists(ctx, params.AccountID, params.OrganizationID)
	if err != nil {
		return models.Invite{}, err
	}
	if memberExists {
		return models.Invite{}, errx.ErrorAccountAlreadyMember.Raise(
			fmt.Errorf("account with id %s is already a member of organization %s", params.AccountID, params.OrganizationID),
		)
	}

	exist, err = m.inviteRepo.ExistsActiveInviteByAccountID(ctx, params.AccountID, params.OrganizationID)
	if err != nil {
		return models.Invite{}, err
	}
	if exist {
		return models.Invite{}, errx.ErrorActiveInviteAlreadyExists.Raise(
			fmt.Errorf("active invite for account %s already exists in organization %s", params.AccountID, params.OrganizationID),
		)
	}

	err = m.tx.Transaction(ctx, func(ctx context.Context) error {
		invite, err = m.inviteRepo.CreateInvite(ctx, params)
		if err != nil {
			return err
		}

		return m.messenger.WriteOrgInviteCreated(ctx, invite)
	})

	return invite, err
}

// GetListForOrganization - получить список инвайтов организации
func (m *InviteModule) GetListForOrganization(
	ctx context.Context,
	actor models.AccountActor,
	organizationID uuid.UUID,
	limit, offset uint,
) (pagi.Page[[]models.Invite], error) {
	_, err := m.getInitiator(ctx, actor, organizationID)
	if err != nil {
		return pagi.Page[[]models.Invite]{}, err
	}

	return m.inviteRepo.GetOrganizationInvites(ctx, organizationID, limit, offset)
}

// GetListForAccount - получить список инвайтов аккаунта
func (m *InviteModule) GetListForAccount(
	ctx context.Context,
	actor models.AccountActor,
	limit, offset uint,
) (pagi.Page[[]models.Invite], error) {
	return m.inviteRepo.GetAccountInvites(ctx, actor, limit, offset)
}

// GetForAccount - получить инвайт для аккаунта
func (m *InviteModule) GetForAccount(
	ctx context.Context,
	accountID uuid.UUID,
	inviteID uuid.UUID,
) (models.Invite, error) {
	res, err := m.inviteRepo.GetInvite(ctx, inviteID)
	if err != nil {
		return models.Invite{}, err
	}

	if res.AccountID != accountID {
		_, err = m.getInitiator(ctx, accountID, res.OrganizationID)
		if err != nil {
			return models.Invite{}, err
		}
	}

	return res, nil
}

// Accept - принять инвайт
func (m *InviteModule) Accept(
	ctx context.Context,
	actor models.AccountActor,
	inviteID uuid.UUID,
) (invite models.Invite, err error) {
	invite, err = m.inviteRepo.GetInvite(ctx, inviteID)
	if err != nil {
		return models.Invite{}, err
	}

	if invite.AccountID != actor {
		return models.Invite{}, errx.ErrorInviteNotForInitiator.Raise(
			fmt.Errorf("account has no rights to accept this invite"),
		)
	}
	if invite.Status == models.InviteStatusAccepted {
		return invite, nil
	}
	if invite.Status != models.InviteStatusSent {
		return models.Invite{}, errx.ErrorInviteAlreadyAnswered.Raise(
			fmt.Errorf("invite status is %s", invite.Status),
		)
	}
	if invite.ExpiresAt.Before(time.Now().UTC()) {
		return models.Invite{}, errx.ErrorInviteExpired.Raise(
			fmt.Errorf("invite expired at %s", invite.ExpiresAt),
		)
	}

	org, err := m.orgRepo.GetOrganizationByID(ctx, invite.OrganizationID)
	if err != nil {
		return models.Invite{}, err
	}
	if org.Status != models.OrganizationStatusActive {
		return models.Invite{}, errx.ErrorOrganizationIsNotActive.Raise(
			fmt.Errorf("organization with id %s is not active", invite.OrganizationID),
		)
	}

	if err = m.tx.Transaction(ctx, func(ctx context.Context) error {
		invite, err = m.inviteRepo.UpdateInviteStatus(ctx, inviteID, models.InviteStatusAccepted)
		if err != nil {
			return err
		}

		mem, err := m.memberRepo.CreateMember(ctx, actor, invite.OrganizationID)
		if err != nil {
			return err
		}

		if err = m.messenger.WriteOrgInviteAccepted(ctx, invite); err != nil {
			return err
		}

		return m.messenger.WriteOrgMemberCreated(ctx, mem)
	}); err != nil {
		return models.Invite{}, err
	}

	return invite, err
}

// Decline - отклонить инвайт
func (m *InviteModule) Decline(
	ctx context.Context,
	actor models.AccountActor,
	inviteID uuid.UUID,
) (invite models.Invite, err error) {
	invite, err = m.inviteRepo.GetInvite(ctx, inviteID)
	if err != nil {
		return models.Invite{}, err
	}
	if invite.AccountID != actor {
		return models.Invite{}, errx.ErrorInviteNotForInitiator.Raise(
			fmt.Errorf("account has no rights to decline this invite"),
		)
	}
	if invite.Status == models.InviteStatusDeclined {
		return invite, nil
	}
	if invite.Status != models.InviteStatusSent {
		return models.Invite{}, errx.ErrorInviteAlreadyAnswered.Raise(
			fmt.Errorf("invite status is %s", invite.Status),
		)
	}
	if invite.ExpiresAt.Before(time.Now().UTC()) {
		return models.Invite{}, errx.ErrorInviteExpired.Raise(
			fmt.Errorf("invite expired at %s", invite.ExpiresAt),
		)
	}

	org, err := m.orgRepo.GetOrganizationByID(ctx, invite.OrganizationID)
	if err != nil {
		return models.Invite{}, err
	}
	if org.Status != models.OrganizationStatusActive {
		return models.Invite{}, errx.ErrorOrganizationIsNotActive.Raise(
			fmt.Errorf("organization with id %s is not active", invite.OrganizationID),
		)
	}

	if err = m.tx.Transaction(ctx, func(ctx context.Context) error {
		invite, err = m.inviteRepo.UpdateInviteStatus(ctx, inviteID, models.InviteStatusDeclined)
		if err != nil {
			return err
		}

		return m.messenger.WriteOrgInviteDeclined(ctx, invite)
	}); err != nil {
		return models.Invite{}, err
	}

	return invite, err
}

// Cancelled - отменить инвайт
func (m *InviteModule) Cancelled(
	ctx context.Context,
	actor models.AccountActor,
	inviteID uuid.UUID,
) (models.Invite, error) {
	invite, err := m.inviteRepo.GetInvite(ctx, inviteID)
	if err != nil {
		return models.Invite{}, err
	}

	_, _, err = m.authorizeOrgHead(ctx, actor, invite.OrganizationID)
	if err != nil {
		return models.Invite{}, err
	}

	if invite.Status == models.InviteStatusCancelled {
		return models.Invite{}, nil
	}
	if invite.Status != models.InviteStatusSent {
		return models.Invite{}, errx.ErrorInviteAlreadyAnswered.Raise(
			fmt.Errorf("invite status is %s", invite.Status),
		)
	}

	if err = m.tx.Transaction(ctx, func(ctx context.Context) error {
		invite, err = m.inviteRepo.UpdateInviteStatus(ctx, inviteID, models.InviteStatusCancelled)
		if err != nil {
			return err
		}

		return m.messenger.WriteOrgInviteCanceled(ctx, invite)
	}); err != nil {
		return models.Invite{}, err
	}

	return invite, nil
}
