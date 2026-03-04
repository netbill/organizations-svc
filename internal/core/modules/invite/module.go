package invite

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/models"
	"github.com/netbill/restkit/pagi"
)

type Invite struct {
	repo      repo
	messenger messenger
}

func New(repo repo, messenger messenger) *Invite {
	return &Invite{
		repo:      repo,
		messenger: messenger,
	}
}

type inviteRepo interface {
	//organization
	GetOrganizationByID(ctx context.Context, ID uuid.UUID) (models.Organization, error)

	//invite
	CreateInvite(
		ctx context.Context,
		params CreateParams,
	) (models.Invite, error)
	GetInvite(
		ctx context.Context,
		id uuid.UUID,
	) (models.Invite, error)
	GetOrganizationInvites(
		ctx context.Context,
		organizationID uuid.UUID,
		limit, offset uint,
	) (pagi.Page[[]models.Invite], error)
	GetAccountInvites(
		ctx context.Context,
		accountID uuid.UUID,
		limit, offset uint,
	) (pagi.Page[[]models.Invite], error)
	ExistsActiveInviteByAccountID(
		ctx context.Context,
		accountID, organizationID uuid.UUID,
	) (bool, error)
	UpdateInviteStatus(
		ctx context.Context,
		id uuid.UUID,
		status string,
	) (models.Invite, error)
	DeleteInvite(
		ctx context.Context,
		id uuid.UUID,
	) error

	//member
	CreateMember(
		ctx context.Context,
		accountID, organizationID uuid.UUID,
	) (models.Member, error)
	GetMemberByAccountAndOrganization(
		ctx context.Context,
		accountID, organizationID uuid.UUID,
	) (models.Member, error)
	MemberExists(
		ctx context.Context,
		accountID, organizationID uuid.UUID,
	) (bool, error)

	//profile
	ExistsProfileByAccountID(ctx context.Context, accountID uuid.UUID) (bool, error)

	//tombstone
	BuryInvite(ctx context.Context, inviteID uuid.UUID) error
	InviteIsBuried(ctx context.Context, inviteID uuid.UUID) (bool, error)

	//transaction
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type inviteMessenger interface {
	WriteOrgMemberCreated(ctx context.Context, member models.Member) error

	WriteOrgInviteCreated(ctx context.Context, invite models.Invite) error
	WriteOrgInviteAccepted(ctx context.Context, invite models.Invite) error
	WriteOrgInviteDeclined(ctx context.Context, invite models.Invite) error
	WriteOrgInviteCanceled(ctx context.Context, invite models.Invite) error
}

func (m *Invite) getInitiator(
	ctx context.Context,
	actor models.AccountActor,
	organizationID uuid.UUID,
) (models.Member, error) {
	row, err := m.repo.GetMemberByAccountAndOrganization(ctx, actor, organizationID)
	if errors.Is(err, errx.ErrorMemberNotExists) {
		return models.Member{}, errx.ErrorInitiatorNotMemberOfOrganization.Raise(
			fmt.Errorf("initiator with account id %s is not a member of organization %s", actor, organizationID),
		)
	}

	return row, err
}
