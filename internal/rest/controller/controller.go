package controller

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/logium"
	"github.com/netbill/organizations-svc/internal/core/models"
	"github.com/netbill/organizations-svc/internal/core/modules/invite"
	"github.com/netbill/organizations-svc/internal/core/modules/member"
	"github.com/netbill/organizations-svc/internal/core/modules/organization"
	"github.com/netbill/organizations-svc/internal/core/modules/role"
	"github.com/netbill/pagi"
)

type aggloSvc interface {
	CreateOrganization(
		ctx context.Context,
		accountID uuid.UUID,
		params organization.CreateParams,
	) (models.Organization, error)

	GetOrganizations(
		ctx context.Context,
		params organization.FilterParams,
		offset, limit uint,
	) (pagi.Page[[]models.Organization], error)
	GetOrganizationForUser(
		ctx context.Context,
		accountID uuid.UUID,
		offset, limit uint,
	) (pagi.Page[[]models.Organization], error)
	GetOrganization(
		ctx context.Context,
		organizationID uuid.UUID,
	) (models.Organization, error)

	UpdateOrganization(
		ctx context.Context,
		accountID, organizationID uuid.UUID,
		params organization.UpdateParams,
	) (models.Organization, error)

	ActivateOrganization(
		ctx context.Context,
		accountID, organizationID uuid.UUID,
	) (models.Organization, error)
	DeactivateOrganization(
		ctx context.Context,
		accountID, organizationID uuid.UUID,
	) (models.Organization, error)
	SuspendOrganization(ctx context.Context, ID uuid.UUID) (models.Organization, error)

	DeleteOrganization(
		ctx context.Context,
		accountID, organizationID uuid.UUID,
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

	GetOrganizationInvites(
		ctx context.Context,
		accountID, organizationID uuid.UUID,
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
	GetMemberByAccountAndOrganization(ctx context.Context, accountID, organizationID uuid.UUID) (models.Member, error)
	GetInitiatorMember(ctx context.Context, accountID, organizationID uuid.UUID) (models.Member, error)

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
		organizationID uuid.UUID,
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

type Controller struct {
	core core
	log  logium.Logger
}

func New(
	aggloSvc aggloSvc,
	memberSvc memberSvc,
	roleSvc roleSvc,
	inviteSvc inviteSvc,
	log logium.Logger,
) Controller {
	return Controller{
		core: core{
			aggloSvc:  aggloSvc,
			inviteSvc: inviteSvc,
			memberSvc: memberSvc,
			roleSvc:   roleSvc,
		},
		log: log,
	}
}
