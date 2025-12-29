package controller

import (
	"context"

	"github.com/google/uuid"
	"github.com/paulmach/orb"
	"github.com/umisto/cities-svc/internal/domain/models"
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
	CreateAgglomeration(ctx context.Context, name string) (models.Agglomeration, error)

	GetAgglomeration(ctx context.Context, ID uuid.UUID) (models.Agglomeration, error)
	FilterAgglomerations(
		ctx context.Context,
		params agglomeration.FilterParams,
		offset, limit uint,
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

type City interface {
	CreateCity(ctx context.Context, params city.CreateParams) (models.City, error)

	GetCityByID(ctx context.Context, id uuid.UUID) (models.City, error)
	GetCityBySlug(ctx context.Context, slug string) (models.City, error)
	FilterCities(
		ctx context.Context,
		params city.FilterParams,
		offset, limit uint,
	) (pagi.Page[[]models.City], error)
	FilterCitiesNearest(
		ctx context.Context,
		filter city.FilterParams,
		point orb.Point,
		offset, limit uint,
	) (pagi.Page[map[float64]models.City], error)

	UpdateCity(ctx context.Context, id uuid.UUID, params city.UpdateParams) (models.City, error)
	UpdateCityByUser(
		ctx context.Context,
		accountID, cityID uuid.UUID,
		params city.UpdateParams,
	) (models.City, error)

	UpdateCitySlug(
		ctx context.Context,
		id uuid.UUID,
		newSlug *string,
	) (city models.City, err error)
	UpdateCitySlugByUser(
		ctx context.Context,
		accountID, cityID uuid.UUID,
		newSlug *string,
	) (models.City, error)

	UpdateCityAgglomeration(
		ctx context.Context,
		id uuid.UUID,
		newAggloID *uuid.UUID,
	) (city models.City, err error)
	UpdateCityAgglomerationByUser(
		ctx context.Context,
		accountID, cityID uuid.UUID,
		newAggloID *uuid.UUID,
	) (models.City, error)

	ActivateCity(ctx context.Context, ID uuid.UUID) (models.City, error)
	ActivateCityByUser(
		ctx context.Context,
		accountID, cityID uuid.UUID,
	) (models.City, error)

	DeactivateCity(ctx context.Context, ID uuid.UUID) (models.City, error)
	DeactivateCityByUser(
		ctx context.Context,
		accountID, cityID uuid.UUID,
	) (models.City, error)

	ArchiveCity(ctx context.Context, ID uuid.UUID) (models.City, error)
}

type Member interface {
	GetMemberByID(ctx context.Context, ID uuid.UUID) (models.Member, error)
	GetMemberByAccountAndAgglomeration(ctx context.Context, accountID, agglomerationID uuid.UUID) (models.Member, error)
	GetInitiatorMember(ctx context.Context, accountID, agglomerationID uuid.UUID) (models.Member, error)

	FilterMembers(
		ctx context.Context,
		filter member.FilterParams,
		offset uint,
		limit uint,
	) (pagi.Page[[]models.Member], error)

	UpdateMember(ctx context.Context, ID uuid.UUID, params member.UpdateParams) (models.Member, error)
	UpdateMemberByUser(
		ctx context.Context,
		accountID, memberID uuid.UUID,
		params member.UpdateParams,
	) (models.Member, error)

	DeleteMember(ctx context.Context, memberID uuid.UUID) error
	DeleteMemberByUser(
		ctx context.Context,
		accountID, memberID uuid.UUID,
	) error
}

type Invite interface {
	CreateInvite(ctx context.Context, params invite.CreateParams) (models.Invite, error)
	CreateInviteByUser(
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

type Role interface {
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
		order map[uint]uuid.UUID,
	) error
	UpdateRoleRank(
		ctx context.Context,
		roleID uuid.UUID,
		newRank uint,
	) (models.Role, error)
	UpdateRoleRankByUser(
		ctx context.Context,
		accountID uuid.UUID,
		roleID uuid.UUID,
		newRank uint,
	) (models.Role, error)

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

func New(agglo Agglomeration, city City, member Member, role Role, invite Invite) Service {
	return Service{
		Agglomeration: agglo,
		City:          city,
		Member:        member,
		Role:          role,
		Invite:        invite,
	}
}
