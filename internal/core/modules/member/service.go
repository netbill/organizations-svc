package member

import (
	"context"

	"github.com/google/uuid"
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
	CreateMember(ctx context.Context, accountID, organizationID uuid.UUID) (models.Member, error)

	UpdateMember(ctx context.Context, ID uuid.UUID, params UpdateParams) (models.Member, error)

	GetMember(ctx context.Context, memberID uuid.UUID) (models.Member, error)
	GetMemberByAccountAndOrganization(
		ctx context.Context,
		accountID, organizationID uuid.UUID,
	) (models.Member, error)
	GetMembers(
		ctx context.Context,
		filter FilterParams,
		limit uint,
		offset uint,
	) (pagi.Page[[]models.Member], error)

	DeleteMember(ctx context.Context, memberID uuid.UUID) error

	CheckMemberHavePermission(
		ctx context.Context,
		memberID uuid.UUID,
		permissionCode string,
	) (bool, error)
	GetMemberMaxRole(ctx context.Context, memberID uuid.UUID) (models.Role, error)

	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type messenger interface {
	WriteMemberCreated(ctx context.Context, member models.Member) error
	WriteMemberUpdated(ctx context.Context, member models.Member) error
	WriteMemberDeleted(ctx context.Context, member models.Member) error
}
