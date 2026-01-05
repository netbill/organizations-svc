package invite

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/models"
	"github.com/netbill/pagi"
)

type Service struct {
	repo      repo
	messenger messenger
}

func New(repo repo, messenger messenger) Service {
	return Service{
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
	WriteMemberCreated(ctx context.Context, member models.Member) error

	WriteInviteCreated(ctx context.Context, invite models.Invite) error

	WriteInviteAccepted(ctx context.Context, invite models.Invite) error
	WriteInviteDeclined(ctx context.Context, invite models.Invite) error

	WriteInviteDeleted(ctx context.Context, invite models.Invite) error
}

func (s Service) checkPermissionForManageInvite(
	ctx context.Context,
	memberID uuid.UUID,
) error {
	access, err := s.repo.CheckMemberHavePermission(
		ctx,
		memberID,
		models.RolePermissionManageInvites,
	)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to check initiator permissions: %w", err))
	}
	if !access {
		return errx.ErrorNotEnoughRights.Raise(
			fmt.Errorf("initiator has no access to activate organization"),
		)
	}

	return nil
}

func (s Service) checkOrganizationIsActiveAndExists(ctx context.Context, organizationID uuid.UUID) (models.Organization, error) {
	org, err := s.repo.GetOrganizationByID(ctx, organizationID)
	if err != nil {
		return models.Organization{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get organization by id: %w", err),
		)
	}
	if org.IsNil() {
		return models.Organization{}, errx.ErrorOrganizationNotFound.Raise(
			fmt.Errorf("organization with id %s not found", organizationID),
		)
	}

	if org.Status != models.OrganizationStatusActive {
		return models.Organization{}, errx.ErrorOrganizationIsNotActive.Raise(
			fmt.Errorf("organization with id %s is not active", organizationID),
		)
	}

	return org, nil
}

func (s Service) getInitiator(ctx context.Context, accountID, organizationID uuid.UUID) (models.Member, error) {
	row, err := s.repo.GetMemberByAccountAndOrganization(ctx, accountID, organizationID)
	if err != nil {
		return models.Member{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to get member with account id %s and organization id %s: %w",
				accountID, organizationID, err),
		)
	}
	if row.IsNil() {
		return models.Member{}, errx.ErrorNotEnoughRights.Raise(
			fmt.Errorf("member with account id %s and organization id %s not found", accountID, organizationID),
		)
	}

	return row, nil
}
