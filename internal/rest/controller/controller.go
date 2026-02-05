package controller

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/netbill/logium"
	"github.com/netbill/organizations-svc/internal/core/models"
	"github.com/netbill/organizations-svc/internal/core/modules/invite"
	"github.com/netbill/organizations-svc/internal/core/modules/member"
	"github.com/netbill/organizations-svc/internal/core/modules/organization"
	"github.com/netbill/organizations-svc/internal/core/modules/role"
	"github.com/netbill/restkit/pagi"
)

type orgSvc interface {
	CreateOrganization(
		ctx context.Context,
		initiator models.InitiatorData,
		params organization.CreateParams,
	) (models.Organization, error)

	GetOrganizations(
		ctx context.Context,
		params organization.FilterParams,
		limit, offset uint,
	) (pagi.Page[[]models.Organization], error)
	GetOrganizationForUser(
		ctx context.Context,
		initiator models.InitiatorData,
		limit, offset uint,
	) (pagi.Page[[]models.Organization], error)
	GetOrganization(
		ctx context.Context,
		organizationID uuid.UUID,
	) (models.Organization, error)

	OpenUpdateOrganizationSession(
		ctx context.Context,
		initiator models.InitiatorData,
		organizationID uuid.UUID,
	) (models.Organization, models.UpdateOrganizationMedia, error)
	UpdateOrganization(
		ctx context.Context,
		initiator models.InitiatorData,
		organizationID uuid.UUID,
		params organization.UpdateParams,
	) (models.Organization, error)
	DeleteUpdateOrganizationIconInSession(
		ctx context.Context,
		initiator models.InitiatorData, organizationID, uploadSessionID uuid.UUID,
	) error
	DeleteUpdateOrganizationBannerInSession(
		ctx context.Context,
		initiator models.InitiatorData,
		organizationID, uploadSessionID uuid.UUID,
	) error

	ActivateOrganization(
		ctx context.Context,
		initiator models.InitiatorData,
		organizationID uuid.UUID,
	) (models.Organization, error)
	DeactivateOrganization(
		ctx context.Context,
		initiator models.InitiatorData,
		organizationID uuid.UUID,
	) (models.Organization, error)

	DeleteOrganization(
		ctx context.Context,
		initiator models.InitiatorData,
		organizationID uuid.UUID,
	) error
}

type inviteSvc interface {
	CreateInvite(
		ctx context.Context,
		initiator models.InitiatorData,
		params invite.CreateParams,
	) (models.Invite, error)

	GetInviteForAccount(
		ctx context.Context,
		initiator models.InitiatorData,
		inviteID uuid.UUID,
	) (models.Invite, error)

	DeclineInvite(
		ctx context.Context,
		initiator models.InitiatorData,
		inviteID uuid.UUID,
	) (models.Invite, error)
	AcceptInvite(
		ctx context.Context,
		initiator models.InitiatorData,
		inviteID uuid.UUID,
	) (models.Invite, error)

	DeleteInvite(
		ctx context.Context,
		initiator models.InitiatorData,
		inviteID uuid.UUID,
	) error

	GetOrganizationInvites(
		ctx context.Context,
		initiator models.InitiatorData,
		organizationID uuid.UUID,
		limit, offset uint,
	) (pagi.Page[[]models.Invite], error)
	GetAccountInvites(
		ctx context.Context,
		initiator models.InitiatorData,
		limit, offset uint,
	) (pagi.Page[[]models.Invite], error)
}

type memberSvc interface {
	GetMemberByID(
		ctx context.Context,
		memberID uuid.UUID,
	) (models.Member, error)

	GetMemberByAccountAndOrganization(
		ctx context.Context,
		initiator models.InitiatorData,
		organizationID uuid.UUID,
	) (models.Member, error)

	GetInitiatorMember(
		ctx context.Context,
		initiator models.InitiatorData,
		organizationID uuid.UUID,
	) (models.Member, error)

	GetMembers(
		ctx context.Context,
		filter member.FilterParams,
		limit, offset uint,
	) (pagi.Page[[]models.Member], error)

	UpdateMember(
		ctx context.Context,
		initiator models.InitiatorData,
		memberID uuid.UUID,
		params member.UpdateParams,
	) (models.Member, error)

	DeleteMember(
		ctx context.Context,
		initiator models.InitiatorData,
		memberID uuid.UUID,
	) error
}

type roleSvc interface {
	CreateRole(
		ctx context.Context,
		initiator models.InitiatorData,
		params role.CreateParams,
	) (models.Role, error)

	GetRole(
		ctx context.Context,
		roleID uuid.UUID,
	) (models.Role, error)

	GetRoleWithPermissions(
		ctx context.Context,
		initiator models.InitiatorData,
		roleID uuid.UUID,
	) (models.Role, models.OrgRolePermissionLinks, error)

	GetRoles(
		ctx context.Context,
		params role.FilterParams,
		limit, offset uint,
	) (pagi.Page[[]models.Role], error)

	UpdateRole(
		ctx context.Context,
		initiator models.InitiatorData,
		roleID uuid.UUID,
		params role.UpdateParams,
	) (models.Role, error)

	UpdateRolesRanks(
		ctx context.Context,
		initiator models.InitiatorData,
		organizationID uuid.UUID,
		order map[uuid.UUID]uint,
	) error

	DeleteRole(
		ctx context.Context,
		initiator models.InitiatorData,
		roleID uuid.UUID,
	) error

	SetRolePermissions(
		ctx context.Context,
		initiator models.InitiatorData,
		roleID uuid.UUID,
		permissions models.OrgRolePermissionDict,
	) (models.Role, models.OrgRolePermissionLinks, error)

	MemberAddRole(
		ctx context.Context,
		initiator models.InitiatorData,
		memberID, roleID uuid.UUID,
	) error

	MemberRemoveRole(
		ctx context.Context,
		initiator models.InitiatorData,
		memberID, roleID uuid.UUID,
	) error

	GetAllPermissions(ctx context.Context) ([]models.OrgRolePermission, error)
}

type responser interface {
	Render(w http.ResponseWriter, status int, res ...interface{})
	RenderErr(w http.ResponseWriter, errs ...error)
}

type core struct {
	orgSvc
	inviteSvc
	memberSvc
	roleSvc
}

type Controller struct {
	core      core
	responser responser
	log       *logium.Logger
}

func New(
	log *logium.Logger,
	responser responser,
	orgSvc orgSvc,
	memberSvc memberSvc,
	roleSvc roleSvc,
	inviteSvc inviteSvc,
) *Controller {
	return &Controller{
		core: core{
			orgSvc:    orgSvc,
			inviteSvc: inviteSvc,
			memberSvc: memberSvc,
			roleSvc:   roleSvc,
		},
		log:       log,
		responser: responser,
	}
}
