package controller

import (
	"context"

	"github.com/google/uuid"
	"github.com/umisto/agglomerations-svc/internal/domain/models"
	"github.com/umisto/agglomerations-svc/internal/domain/modules/agglomeration"
	"github.com/umisto/agglomerations-svc/internal/domain/modules/invite"
	"github.com/umisto/agglomerations-svc/internal/domain/modules/member"
	"github.com/umisto/agglomerations-svc/internal/domain/modules/role"
	"github.com/umisto/logium"
	"github.com/umisto/pagi"
)

type aggloSvc interface {
	CreateAgglomeration(
		ctx context.Context,
		accountID uuid.UUID,
		params agglomeration.CreateParams,
	) (models.Agglomeration, error)

	GetAgglomeration(ctx context.Context, ID uuid.UUID) (models.Agglomeration, error)
	FilterAgglomerations(
		ctx context.Context,
		params agglomeration.FilterParams,
		offset, limit uint,
	) (pagi.Page[[]models.Agglomeration], error)
	GetAgglomerationForUser(
		ctx context.Context,
		accountID uuid.UUID,
		limit, offset uint,
	) (pagi.Page[[]models.Agglomeration], error)

	UpdateAgglomeration(ctx context.Context, ID uuid.UUID, params agglomeration.UpdateParams) (models.Agglomeration, error)
	UpdateAgglomerationByUser(
		ctx context.Context,
		accountID, agglomerationID uuid.UUID,
		params agglomeration.UpdateParams,
	) (models.Agglomeration, error)

	ActivateAgglomeration(ctx context.Context, ID uuid.UUID) (models.Agglomeration, error)
	ActivateAgglomerationByUser(
		ctx context.Context,
		accountID, agglomerationID uuid.UUID,
	) (models.Agglomeration, error)

	DeactivateAgglomeration(ctx context.Context, ID uuid.UUID) (models.Agglomeration, error)
	DeactivateAgglomerationByUser(
		ctx context.Context,
		accountID, agglomerationID uuid.UUID,
	) (models.Agglomeration, error)

	SuspendAgglomeration(ctx context.Context, ID uuid.UUID) (models.Agglomeration, error)

	DeleteAgglomeration(ctx context.Context, ID uuid.UUID) error
}

type inviteSvc interface {
	CreateInvite(ctx context.Context, params invite.CreateParams) (models.Invite, error)
	SentInviteByUser(
		ctx context.Context,
		accountID uuid.UUID,
		params invite.CreateParams,
	) (models.Invite, error)

	GetInvite(ctx context.Context, id uuid.UUID) (models.Invite, error)
	FilterInvites(
		ctx context.Context,
		filter invite.FilterInviteParams,
	) ([]models.Invite, error)

	DeclineInvite(
		ctx context.Context,
		accountID, inviteID uuid.UUID,
	) (models.Invite, error)
	AcceptInvite(
		ctx context.Context,
		accountID, inviteID uuid.UUID,
	) (models.Invite, error)

	DeleteInvite(
		ctx context.Context,
		accountID, inviteID uuid.UUID,
	) error
}

type memberSvc interface {
	GetMemberByID(ctx context.Context, ID uuid.UUID) (models.Member, error)
	GetMemberByAccountAndAgglomeration(ctx context.Context, accountID, agglomerationID uuid.UUID) (models.Member, error)
	GetInitiatorMember(ctx context.Context, accountID, agglomerationID uuid.UUID) (models.Member, error)

	FilterMembers(
		ctx context.Context,
		filter member.FilterParams,
		offset uint,
		limit uint,
	) (pagi.Page[[]models.Member], error)

	UpdateMemberByUser(
		ctx context.Context,
		accountID, memberID uuid.UUID,
		params member.UpdateParams,
	) (models.Member, error)

	DeleteMemberByUser(
		ctx context.Context,
		accountID, memberID uuid.UUID,
	) error
}

type roleSvc interface {
	CreateRole(ctx context.Context, params role.CreateParams) (models.Role, error)
	CreateRoleByUser(
		ctx context.Context,
		initiatorID uuid.UUID,
		params role.CreateParams,
	) (models.Role, error)

	GetRole(ctx context.Context, roleID uuid.UUID) (models.Role, error)
	GetRoleWithPermissions(ctx context.Context, roleID uuid.UUID) (models.Role, []models.Permission, error)
	FilterRoles(
		ctx context.Context,
		params role.FilterParams,
		offset uint,
		limit uint,
	) (pagi.Page[[]models.Role], error)

	UpdateRoleByUser(
		ctx context.Context,
		accountID uuid.UUID,
		roleID uuid.UUID,
		params role.UpdateParams,
	) (models.Role, error)

	UpdateRolesRanksByUser(
		ctx context.Context,
		accountID uuid.UUID,
		agglomerationID uuid.UUID,
		order map[uuid.UUID]uint,
	) error

	DeleteRole(ctx context.Context, roleID uuid.UUID) error
	DeleteRoleByUser(ctx context.Context, accountID, roleID uuid.UUID) error

	GetRolePermissions(ctx context.Context, roleID uuid.UUID) ([]models.Permission, error)
	SetRolePermissions(
		ctx context.Context,
		roleID uuid.UUID,
		permissions map[models.CodeRolePermission]bool,
	) ([]models.Permission, error)
	SetRolePermissionsByUser(
		ctx context.Context,
		accountID, roleID uuid.UUID,
		permissions map[models.CodeRolePermission]bool,
	) ([]models.Permission, error)

	GetAllPermissions(ctx context.Context) ([]models.Permission, error)
}

type core struct {
	aggloSvc
	inviteSvc
	memberSvc
	roleSvc
}

type Service struct {
	core core
	log  logium.Logger
}
