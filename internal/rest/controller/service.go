package controller

import (
	"context"

	"github.com/google/uuid"
	"github.com/paulmach/orb"
	"github.com/umisto/cities-svc/internal/domain/entity"
	"github.com/umisto/cities-svc/internal/domain/modules/agglomeration"
	"github.com/umisto/cities-svc/internal/domain/modules/city"
	"github.com/umisto/cities-svc/internal/domain/modules/invite"
	"github.com/umisto/cities-svc/internal/domain/modules/member"
	"github.com/umisto/cities-svc/internal/domain/modules/role"
	"github.com/umisto/pagi"
)

type Service struct {
	Agglomeration Agglomeration
	City          City
	Member        Member
	Role          Role
	Invite        Invite
}

type Agglomeration interface {
	CreateAgglomeration(ctx context.Context, name string) (entity.Agglomeration, error)

	GetAgglomeration(ctx context.Context, ID uuid.UUID) (entity.Agglomeration, error)
	FilterAgglomerations(
		ctx context.Context,
		params agglomeration.FilterParams,
		offset, limit uint,
	) (pagi.Page[[]entity.Agglomeration], error)

	UpdateAgglomeration(ctx context.Context, ID uuid.UUID, params agglomeration.UpdateParams) (entity.Agglomeration, error)
	UpdateAgglomerationByUser(
		ctx context.Context,
		accountID, agglomerationID uuid.UUID,
		params agglomeration.UpdateParams,
	) (entity.Agglomeration, error)

	ActivateAgglomeration(ctx context.Context, ID uuid.UUID) (entity.Agglomeration, error)
	ActivateAgglomerationByUser(
		ctx context.Context,
		accountID, agglomerationID uuid.UUID,
	) (entity.Agglomeration, error)

	DeactivateAgglomeration(ctx context.Context, ID uuid.UUID) (entity.Agglomeration, error)
	DeactivateAgglomerationByUser(
		ctx context.Context,
		accountID, agglomerationID uuid.UUID,
	) (entity.Agglomeration, error)

	SuspendAgglomeration(ctx context.Context, ID uuid.UUID) (entity.Agglomeration, error)

	DeleteAgglomeration(ctx context.Context, ID uuid.UUID) error
}

type City interface {
	CreateCity(ctx context.Context, params city.CreateParams) (entity.City, error)

	GetCity(ctx context.Context, id uuid.UUID) (entity.City, error)
	GetCityBySlug(ctx context.Context, slug string) (entity.City, error)
	FilterCities(
		ctx context.Context,
		params city.FilterParams,
		offset, limit uint,
	) (pagi.Page[[]entity.City], error)
	FilterCitiesNearest(
		ctx context.Context,
		filter city.FilterParams,
		point orb.Point,
		offset, limit uint,
	) (pagi.Page[map[float64]entity.City], error)

	UpdateCity(ctx context.Context, id uuid.UUID, params city.UpdateParams) (entity.City, error)
	UpdateCityByUser(
		ctx context.Context,
		accountID, cityID uuid.UUID,
		params city.UpdateParams,
	) (entity.City, error)

	ActivateCity(ctx context.Context, ID uuid.UUID) (entity.City, error)
	ActivateCityByUser(
		ctx context.Context,
		accountID, cityID uuid.UUID,
	) (entity.City, error)

	DeactivateCity(ctx context.Context, ID uuid.UUID) (entity.City, error)
	DeactivateCityByUser(
		ctx context.Context,
		accountID, cityID uuid.UUID,
	) (entity.City, error)

	ArchiveCity(ctx context.Context, ID uuid.UUID) (entity.City, error)
}

type Member interface {
	GetMemberByID(ctx context.Context, ID uuid.UUID) (entity.Member, error)
	GetMemberByAccountAndAgglomeration(ctx context.Context, accountID, agglomerationID uuid.UUID) (entity.Member, error)
	GetInitiatorMember(ctx context.Context, accountID, agglomerationID uuid.UUID) (entity.Member, error)

	FilterMembers(
		ctx context.Context,
		filter member.FilterParams,
		offset uint,
		limit uint,
	) (pagi.Page[[]entity.Member], error)

	UpdateMember(ctx context.Context, ID uuid.UUID, params member.UpdateParams) (entity.Member, error)
	UpdateMemberByUser(
		ctx context.Context,
		accountID, memberID uuid.UUID,
		params member.UpdateParams,
	) (entity.Member, error)

	DeleteMember(ctx context.Context, memberID uuid.UUID) error
	DeleteMemberByUser(
		ctx context.Context,
		accountID, memberID uuid.UUID,
	) error
}

type Invite interface {
	CreateInvite(ctx context.Context, params invite.CreateParams) (entity.Invite, error)
	CreateInviteByUser(
		ctx context.Context,
		accountID uuid.UUID,
		params invite.CreateParams,
	) (entity.Invite, error)

	GetInvite(ctx context.Context, id uuid.UUID) (entity.Invite, error)
	FilterInvites(
		ctx context.Context,
		filter invite.FilterInviteParams,
	) ([]entity.Invite, error)

	DeclineInvite(
		ctx context.Context,
		accountID, inviteID uuid.UUID,
	) (entity.Invite, error)
	AcceptInvite(
		ctx context.Context,
		accountID, inviteID uuid.UUID,
	) (entity.Invite, error)

	DeleteInvite(
		ctx context.Context,
		accountID, inviteID uuid.UUID,
	) error
}

type Role interface {
	CreateRole(ctx context.Context, params role.CreateParams) (entity.Role, error)
	CreateRoleByUser(
		ctx context.Context,
		initiatorID uuid.UUID,
		params role.CreateParams,
	) (entity.Role, error)

	GetRole(ctx context.Context, roleID uuid.UUID) (entity.Role, error)
	GetRoleWithPermissions(ctx context.Context, roleID uuid.UUID) (entity.Role, []entity.Permission, error)
	FilterRoles(
		ctx context.Context,
		params role.FilterParams,
		offset uint,
		limit uint,
	) (pagi.Page[[]entity.Role], error)

	UpdateRoleByUser(
		ctx context.Context,
		accountID uuid.UUID,
		roleID uuid.UUID,
		params role.UpdateParams,
	) (entity.Role, error)

	UpdateRolesRanksByUser(
		ctx context.Context,
		accountID uuid.UUID,
		agglomerationID uuid.UUID,
		order map[uint]uuid.UUID,
	) error
	UpdateRoleRank(
		ctx context.Context,
		roleID uuid.UUID,
		newRank uint,
	) (entity.Role, error)
	UpdateRoleRankByUser(
		ctx context.Context,
		accountID uuid.UUID,
		roleID uuid.UUID,
		newRank uint,
	) (entity.Role, error)

	DeleteRole(ctx context.Context, roleID uuid.UUID) error
	DeleteRoleByUser(ctx context.Context, accountID, roleID uuid.UUID) error

	GetRolePermissions(ctx context.Context, roleID uuid.UUID) ([]entity.Permission, error)
	SetRolePermissions(ctx context.Context, roleID uuid.UUID, permissions map[entity.CodeRolePermission]bool) ([]entity.Permission, error)

	GetAllPermissions(ctx context.Context) ([]entity.Permission, error)
}

func New(agglo Agglomeration, city City, member Member, role Role, invite Invite) Service {
	return Service{
		Agglomeration: agglo,
		City:          city,
		Member:        member,
		Role:          role,
		Invite:        invite,
	}
}
