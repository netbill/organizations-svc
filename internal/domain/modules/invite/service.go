package invite

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/umisto/cities-svc/internal/domain/errx"
	"github.com/umisto/cities-svc/internal/domain/models"
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

	GetInviteByID(
		ctx context.Context,
		id uuid.UUID,
	) (models.Invite, error)
	FilterInvites(
		ctx context.Context,
		filter FilterInviteParams,
	) ([]models.Invite, error)

	UpdateInviteStatus(
		ctx context.Context,
		id uuid.UUID,
		status string,
	) (models.Invite, error)

	DeleteInvite(
		ctx context.Context,
		id uuid.UUID,
	) error

	CheckAccountHavePermissionByCode(
		ctx context.Context,
		accountID, agglomerationID uuid.UUID,
		permissionKey string,
	) (bool, error)

	CreateMember(ctx context.Context, accountID, agglomerationID uuid.UUID) (models.Member, error)

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
	accountID uuid.UUID,
	agglomerationID uuid.UUID,
) error {
	access, err := s.repo.CheckAccountHavePermissionByCode(
		ctx,
		accountID,
		agglomerationID,
		models.RolePermissionManageInvites.String(),
	)
	if err != nil {
		return errx.ErrorInternal.Raise(
			fmt.Errorf("failed to check initiator permissions: %w", err))
	}
	if !access {
		return errx.ErrorNotEnoughRights.Raise(
			fmt.Errorf("initiator has no access to activate agglomeration"),
		)
	}

	return nil
}
