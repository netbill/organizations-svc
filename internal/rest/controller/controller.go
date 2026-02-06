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

type organizationSvc interface {
	Create(
		ctx context.Context,
		initiator models.Initiator,
		params organization.CreateParams,
	) (models.Organization, error)

	GetByID(
		ctx context.Context,
		organizationID uuid.UUID,
	) (models.Organization, error)
	GetList(
		ctx context.Context,
		params organization.FilterParams,
		limit, offset uint,
	) (pagi.Page[[]models.Organization], error)
	GetForUser(
		ctx context.Context,
		initiator models.Initiator,
		limit, offset uint,
	) (pagi.Page[[]models.Organization], error)

	OpenUpdateSession(
		ctx context.Context,
		initiator models.Initiator,
		organizationID uuid.UUID,
	) (models.Organization, models.UpdateOrganizationMedia, error)
	UpdateWithSession(
		ctx context.Context,
		initiator models.Initiator,
		organizationID uuid.UUID,
		params organization.UpdateParams,
	) (models.Organization, error)
	DeleteUpdateIconInSession(
		ctx context.Context,
		initiator models.Initiator, organizationID, uploadSessionID uuid.UUID,
	) error
	DeleteUpdateBannerInSession(
		ctx context.Context,
		initiator models.Initiator,
		organizationID, uploadSessionID uuid.UUID,
	) error

	Activate(
		ctx context.Context,
		initiator models.Initiator,
		organizationID uuid.UUID,
	) (models.Organization, error)
	Deactivate(
		ctx context.Context,
		initiator models.Initiator,
		organizationID uuid.UUID,
	) (models.Organization, error)

	Delete(
		ctx context.Context,
		initiator models.Initiator,
		organizationID uuid.UUID,
	) error
}

type inviteSvc interface {
	Create(
		ctx context.Context,
		initiator models.Initiator,
		params invite.CreateParams,
	) (models.Invite, error)

	Decline(
		ctx context.Context,
		initiator models.Initiator,
		inviteID uuid.UUID,
	) (models.Invite, error)
	Accept(
		ctx context.Context,
		initiator models.Initiator,
		inviteID uuid.UUID,
	) (models.Invite, error)

	Delete(
		ctx context.Context,
		initiator models.Initiator,
		inviteID uuid.UUID,
	) error

	GetForOrganizations(
		ctx context.Context,
		initiator models.Initiator,
		organizationID uuid.UUID,
		limit, offset uint,
	) (pagi.Page[[]models.Invite], error)
	GetForAccount(
		ctx context.Context,
		accountId uuid.UUID,
		inviteID uuid.UUID,
	) (models.Invite, error)
}

type memberSvc interface {
	GetByID(
		ctx context.Context,
		memberID uuid.UUID,
	) (models.Member, error)

	GetByAccountAndOrganization(
		ctx context.Context,
		initiator models.Initiator,
		organizationID uuid.UUID,
	) (models.Member, error)

	GetInitiator(
		ctx context.Context,
		initiator models.Initiator,
		organizationID uuid.UUID,
	) (models.Member, error)

	GetList(
		ctx context.Context,
		filter member.FilterParams,
		limit, offset uint,
	) (pagi.Page[[]models.Member], error)

	Update(
		ctx context.Context,
		initiator models.Initiator,
		memberID uuid.UUID,
		params member.UpdateParams,
	) (models.Member, error)

	Delete(
		ctx context.Context,
		initiator models.Initiator,
		memberID uuid.UUID,
	) error
}

type roleSvc interface {
	Create(
		ctx context.Context,
		initiator models.Initiator,
		params role.CreateParams,
	) (models.Role, error)

	GetByID(
		ctx context.Context,
		roleID uuid.UUID,
	) (models.Role, error)

	GetWithPermissions(
		ctx context.Context,
		initiator models.Initiator,
		roleID uuid.UUID,
	) (models.Role, models.OrgRolePermissionDictWithDetails, error)

	GetList(
		ctx context.Context,
		params role.FilterParams,
		limit, offset uint,
	) (pagi.Page[[]models.Role], error)

	Update(
		ctx context.Context,
		initiator models.Initiator,
		roleID uuid.UUID,
		params role.UpdateParams,
	) (models.Role, error)

	UpdateRanks(
		ctx context.Context,
		initiator models.Initiator,
		organizationID uuid.UUID,
		order map[uuid.UUID]uint,
	) error

	Delete(
		ctx context.Context,
		initiator models.Initiator,
		roleID uuid.UUID,
	) error

	//SetForRole(
	//	ctx context.Context,
	//	initiator models.Initiator,
	//	roleID uuid.UUID,
	//	permissions models.SetPermissionParams,
	//) (models.Role, models.OrgRolePermissionDictWithDetails, error)

	AddForMember(
		ctx context.Context,
		initiator models.Initiator,
		memberID, roleID uuid.UUID,
	) error

	RemoveFromMember(
		ctx context.Context,
		initiator models.Initiator,
		memberID, roleID uuid.UUID,
	) error

	//GetAll(ctx context.Context) ([]models.OrgRolePermission, error)
}

type rolePermissionSvc interface {
	GetAll(ctx context.Context) ([]models.OrgRolePermission, error)

	SetForRole(
		ctx context.Context,
		initiator models.Initiator,
		roleID uuid.UUID,
		permissions models.OrgRolePermissionEnable,
	) (models.Role, models.OrgRolePermissionDictWithDetails, error)
}

type responser interface {
	Render(w http.ResponseWriter, status int, res ...interface{})
	RenderErr(w http.ResponseWriter, errs ...error)
}

type core struct {
	organization organizationSvc
	member       memberSvc
	role         roleSvc
	permissions  rolePermissionSvc
	invite       inviteSvc
}

type Controller struct {
	log       *logium.Logger
	core      *core
	responser responser
}

func New(
	log *logium.Logger,
	responser responser,
	orgSvc organizationSvc,
	memberSvc memberSvc,
	roleSvc roleSvc,
	permSvc rolePermissionSvc,
	inviteSvc inviteSvc,
) *Controller {
	return &Controller{
		log: log,
		core: &core{
			organization: orgSvc,
			member:       memberSvc,
			role:         roleSvc,
			invite:       inviteSvc,
			permissions:  permSvc,
		},
		responser: responser,
	}
}
