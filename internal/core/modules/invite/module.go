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
	CreateInvite(ctx context.Context, params CreateParams) (models.Invite, error)

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

	CheckMemberHavePermission(
		ctx context.Context,
		memberID uuid.UUID,
		permissionCode string,
	) (bool, error)

	CreateMember(ctx context.Context, accountID, organizationID uuid.UUID) (models.Member, error)

	GetOrganizationByID(ctx context.Context, ID uuid.UUID) (models.Organization, error)
	MemberExists(ctx context.Context, accountID, organizationID uuid.UUID) (bool, error)
	GetMemberByAccountAndOrganization(
		ctx context.Context,
		accountID, organizationID uuid.UUID,
	) (models.Member, error)

	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type messenger interface {
	WriteOrgMemberCreated(ctx context.Context, member models.Member) error

	WriteOrgInviteCreated(ctx context.Context, invite models.Invite) error
	WriteOrgInviteAccepted(ctx context.Context, invite models.Invite) error
	WriteOrgInviteDeclined(ctx context.Context, invite models.Invite) error
	WriteOrgInviteDeleted(ctx context.Context, invite models.Invite) error
}

func (m *Module) checkPermissionForManageInvites(
	ctx context.Context,
	initiator models.InitiatorData,
	organizationID uuid.UUID,
) error {
	member, err := m.getInitiator(ctx, initiator, organizationID)
	if err != nil {
		return err
	}

	if !member.Head {
		access, err := m.repo.CheckMemberHavePermission(
			ctx,
			member.ID,
			models.RolePermissionManageInvites,
		)
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
) (models.Organization, error) {
	org, err := m.repo.GetOrganizationByID(ctx, organizationID)
	if err != nil {
		return models.Organization{}, err
	}

	if org.Status != models.OrganizationStatusActive {
		return models.Organization{}, errx.ErrorOrganizationIsNotActive.Raise(
			fmt.Errorf("organization with id %s is not active", organizationID),
		)
	}

	return org, nil
}

func (m *Module) getInitiator(
	ctx context.Context,
	initiator models.InitiatorData,
	organizationID uuid.UUID,
) (models.Member, error) {
	row, err := m.repo.GetMemberByAccountAndOrganization(ctx, initiator.AccountID, organizationID)
	if err != nil {
		if errors.Is(err, errx.ErrorMemberNotFound) {
			return models.Member{}, errx.ErrorNotEnoughRights.Raise(
				fmt.Errorf(
					"initiator with account id %s is not a member of organization %s",
					initiator.AccountID, organizationID,
				),
			)
		}
		return models.Member{}, err
	}

	return row, nil
}
