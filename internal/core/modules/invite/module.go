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

type Module struct {
	repo      repo
	messenger messenger
}

func New(repo repo, messenger messenger) *Module {
	return &Module{
		repo:      repo,
		messenger: messenger,
	}
}

type repo interface {
	GetOrganizationByID(ctx context.Context, ID uuid.UUID) (models.Organization, error)

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

	UpdateInviteStatus(
		ctx context.Context,
		id uuid.UUID,
		status string,
	) (models.Invite, error)
	DeleteInvite(
		ctx context.Context,
		id uuid.UUID,
	) error

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

	ExistsProfileByAccountID(ctx context.Context, accountID uuid.UUID) (bool, error)

	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type messenger interface {
	WriteOrgMemberCreated(ctx context.Context, member models.Member) error

	WriteOrgInviteCreated(ctx context.Context, invite models.Invite) error
	WriteOrgInviteAccepted(ctx context.Context, invite models.Invite) error
	WriteOrgInviteDeclined(ctx context.Context, invite models.Invite) error
	WriteOrgInviteDeleted(ctx context.Context, invite models.Invite) error
}

func (m *Module) getInitiator(
	ctx context.Context,
	actor models.AccountActor,
	organizationID uuid.UUID,
) (models.Member, error) {
	row, err := m.repo.GetMemberByAccountAndOrganization(ctx, actor, organizationID)
	if errors.Is(err, errx.ErrorMemberNotFound) {
		return models.Member{}, errx.ErrorNotEnoughRights.Raise(
			fmt.Errorf("initiator with account id %s is not a member of organization %s", actor, organizationID),
		)
	} else if err != nil {
		return models.Member{}, err
	}

	return row, nil
}
