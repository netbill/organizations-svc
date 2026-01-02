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

	GetAgglomerations(
		ctx context.Context,
		params agglomeration.FilterParams,
		offset, limit uint,
	) (pagi.Page[[]models.Agglomeration], error)
	GetAgglomerationForUser(
		ctx context.Context,
		accountID uuid.UUID,
		offset, limit uint,
	) (pagi.Page[[]models.Agglomeration], error)
	GetAgglomeration(
		ctx context.Context,
		agglomerationID uuid.UUID,
	) (models.Agglomeration, error)

	UpdateAgglomeration(
		ctx context.Context,
		accountID, agglomerationID uuid.UUID,
		params agglomeration.UpdateParams,
	) (models.Agglomeration, error)

	ActivateAgglomeration(
		ctx context.Context,
		accountID, agglomerationID uuid.UUID,
	) (models.Agglomeration, error)
	DeactivateAgglomeration(
		ctx context.Context,
		accountID, agglomerationID uuid.UUID,
	) (models.Agglomeration, error)
	SuspendAgglomeration(ctx context.Context, ID uuid.UUID) (models.Agglomeration, error)

	DeleteAgglomeration(
		ctx context.Context,
		accountID, agglomerationID uuid.UUID,
	) error
}

type inviteSvc interface {
	CreateInvite(
		ctx context.Context,
		accountID uuid.UUID,
		params invite.CreateParams,
	) (models.Invite, error)

	GetInvite(
		ctx context.Context,
		accountID, inviteID uuid.UUID,
	) (models.Invite, error)

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

	GetAgglomerationInvites(
		ctx context.Context,
		accountID, agglomerationID uuid.UUID,
		limit, offset uint,
	) (pagi.Page[[]models.Invite], error)
	GetAccountInvites(
		ctx context.Context,
		accountID uuid.UUID,
		limit, offset uint,
	) (pagi.Page[[]models.Invite], error)
}

type memberSvc interface {
	GetMemberByID(ctx context.Context, ID uuid.UUID) (models.Member, error)
	GetMemberByAccountAndAgglomeration(ctx context.Context, accountID, agglomerationID uuid.UUID) (models.Member, error)
	GetInitiatorMember(ctx context.Context, accountID, agglomerationID uuid.UUID) (models.Member, error)

	GetMembers(
		ctx context.Context,
		filter member.FilterParams,
		offset uint,
		limit uint,
	) (pagi.Page[[]models.Member], error)

	UpdateMember(
		ctx context.Context,
		accountID, memberID uuid.UUID,
		params member.UpdateParams,
	) (models.Member, error)

	DeleteMember(
		ctx context.Context,
		accountID, memberID uuid.UUID,
	) error
}

type roleSvc interface {
	CreateRole(
		ctx context.Context,
		initiatorID uuid.UUID,
		params role.CreateParams,
	) (models.Role, error)

	GetRole(ctx context.Context, roleID uuid.UUID) (models.Role, error)
	GetRoleWithPermissions(ctx context.Context, accountID, roleID uuid.UUID) (models.Role, map[models.Permission]bool, error)
	GetRoles(
		ctx context.Context,
		params role.FilterParams,
		offset uint,
		limit uint,
	) (pagi.Page[[]models.Role], error)

	UpdateRole(
		ctx context.Context,
		accountID uuid.UUID,
		roleID uuid.UUID,
		params role.UpdateParams,
	) (models.Role, error)

	UpdateRolesRanks(
		ctx context.Context,
		accountID uuid.UUID,
		agglomerationID uuid.UUID,
		order map[uuid.UUID]uint,
	) error

	DeleteRole(ctx context.Context, accountID, roleID uuid.UUID) error

	SetRolePermissions(
		ctx context.Context,
		accountID uuid.UUID,
		roleID uuid.UUID,
		permissions map[string]bool,
	) (models.Role, map[models.Permission]bool, error)

	MemberAddRole(
		ctx context.Context,
		accountID, memberID, roleID uuid.UUID,
	) error
	RemoveMemberRole(
		ctx context.Context,
		accountID, memberID, roleID uuid.UUID,
	) error

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

func New(
	aggloSvc aggloSvc,
	memberSvc memberSvc,
	roleSvc roleSvc,
	inviteSvc inviteSvc,
	log logium.Logger,
) Service {
	return Service{
		core: core{
			aggloSvc:  aggloSvc,
			inviteSvc: inviteSvc,
			memberSvc: memberSvc,
			roleSvc:   roleSvc,
		},
		log: log,
	}
}
