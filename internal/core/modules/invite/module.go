package invite

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/domain"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/orgperm"
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
	CreateInvite(ctx context.Context, params CreateParams) (domain.Invite, error)

	GetInvite(
		ctx context.Context,
		id uuid.UUID,
	) (domain.Invite, error)
	GetOrganizationInvites(
		ctx context.Context,
		organizationID uuid.UUID,
		limit, offset uint,
	) (pagi.Page[[]domain.Invite], error)
	GetAccountInvites(
		ctx context.Context,
		accountID uuid.UUID,
		limit, offset uint,
	) (pagi.Page[[]domain.Invite], error)

	UpdateInviteStatus(
		ctx context.Context,
		id uuid.UUID,
		status string,
	) (domain.Invite, error)

	DeleteInvite(
		ctx context.Context,
		id uuid.UUID,
	) error

	CheckMemberHavePermission(
		ctx context.Context,
		memberID uuid.UUID,
		permissionID uuid.UUID,
	) (bool, error)

	CreateMember(ctx context.Context, accountID, organizationID uuid.UUID) (domain.Member, error)

	GetOrganizationByID(ctx context.Context, ID uuid.UUID) (domain.Organization, error)
	MemberExists(ctx context.Context, accountID, organizationID uuid.UUID) (bool, error)
	GetMemberByAccountAndOrganization(
		ctx context.Context,
		accountID, organizationID uuid.UUID,
	) (domain.Member, error)

	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type messenger interface {
	WriteOrgMemberCreated(ctx context.Context, member domain.Member) error

	WriteOrgInviteCreated(ctx context.Context, invite domain.Invite) error
	WriteOrgInviteAccepted(ctx context.Context, invite domain.Invite) error
	WriteOrgInviteDeclined(ctx context.Context, invite domain.Invite) error
	WriteOrgInviteDeleted(ctx context.Context, invite domain.Invite) error
}

func (m *Module) checkPermissionForManageInvites(
	ctx context.Context,
	initiator domain.AccountActor,
	organizationID uuid.UUID,
) error {
	member, err := m.getInitiator(ctx, initiator, organizationID)
	if err != nil {
		return err
	}

	if !member.Head {
		access, err := m.repo.CheckMemberHavePermission(ctx, member.ID, orgperm.InvitesManageID)
		if err != nil {
			return err
		}

		if !access {
			return errx.ErrorNotEnoughRights.Raise(
				fmt.Errorf("initiator has no access to activate organization"),
			)
		}
	}

	return nil
}

func (m *Module) checkOrganizationIsActiveAndExists(
	ctx context.Context,
	organizationID uuid.UUID,
) (domain.Organization, error) {
	org, err := m.repo.GetOrganizationByID(ctx, organizationID)
	if err != nil {
		return domain.Organization{}, err
	}

	if org.Status != domain.OrganizationStatusActive {
		return domain.Organization{}, errx.ErrorOrganizationIsNotActive.Raise(
			fmt.Errorf("organization with id %s is not active", organizationID),
		)
	}

	return org, nil
}

func (m *Module) getInitiator(
	ctx context.Context,
	initiator domain.AccountActor,
	organizationID uuid.UUID,
) (domain.Member, error) {
	row, err := m.repo.GetMemberByAccountAndOrganization(ctx, initiator, organizationID)
	if err != nil {
		if errors.Is(err, errx.ErrorMemberNotFound) {
			return domain.Member{}, errx.ErrorNotEnoughRights.Raise(
				fmt.Errorf(
					"initiator with account id %s is not a member of organization %s",
					initiator, organizationID,
				),
			)
		}
		return domain.Member{}, err
	}

	return row, nil
}
