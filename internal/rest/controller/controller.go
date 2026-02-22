package controller

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/domain"
	"github.com/netbill/organizations-svc/internal/core/modules/invite"
	"github.com/netbill/organizations-svc/internal/core/modules/member"
	"github.com/netbill/organizations-svc/internal/core/modules/organization"
	"github.com/netbill/organizations-svc/internal/core/modules/role"
	"github.com/netbill/restkit/pagi"
)

type organizationSvc interface {
	Create(
		ctx context.Context,
		initiator domain.AccountActor,
		params organization.CreateParams,
	) (domain.Organization, error)

	GetByID(
		ctx context.Context,
		organizationID uuid.UUID,
	) (domain.Organization, error)
	GetList(
		ctx context.Context,
		params organization.FilterParams,
		limit, offset uint,
	) (pagi.Page[[]domain.Organization], error)
	GetForUser(
		ctx context.Context,
		initiator domain.AccountActor,
		limit, offset uint,
	) (pagi.Page[[]domain.Organization], error)

	CreateOrgUploadMediaLinks(
		ctx context.Context,
		actor domain.AccountActor,
		organizationID uuid.UUID,
	) (domain.Organization, domain.UploadOrgMediaLinks, error)
	Update(
		ctx context.Context,
		initiator domain.AccountActor,
		organizationID uuid.UUID,
		params organization.UpdateParams,
	) (domain.Organization, error)
	DeleteOrgUploadIcon(
		ctx context.Context,
		actor domain.AccountActor,
		organizationID uuid.UUID,
		key string,
	) error
	DeleteOrgUploadBanner(
		ctx context.Context,
		actor domain.AccountActor,
		organizationID uuid.UUID,
		key string,
	) error

	Activate(
		ctx context.Context,
		initiator domain.AccountActor,
		organizationID uuid.UUID,
	) (domain.Organization, error)
	Deactivate(
		ctx context.Context,
		initiator domain.AccountActor,
		organizationID uuid.UUID,
	) (domain.Organization, error)

	Delete(
		ctx context.Context,
		initiator domain.AccountActor,
		organizationID uuid.UUID,
	) error
}

type inviteSvc interface {
	Create(
		ctx context.Context,
		initiator domain.AccountActor,
		params invite.CreateParams,
	) (domain.Invite, error)

	Decline(
		ctx context.Context,
		initiator domain.AccountActor,
		inviteID uuid.UUID,
	) (domain.Invite, error)
	Accept(
		ctx context.Context,
		initiator domain.AccountActor,
		inviteID uuid.UUID,
	) (domain.Invite, error)

	Delete(
		ctx context.Context,
		initiator domain.AccountActor,
		inviteID uuid.UUID,
	) error

	GetForOrganizations(
		ctx context.Context,
		initiator domain.AccountActor,
		organizationID uuid.UUID,
		limit, offset uint,
	) (pagi.Page[[]domain.Invite], error)
	GetForAccount(
		ctx context.Context,
		accountId uuid.UUID,
		inviteID uuid.UUID,
	) (domain.Invite, error)
}

type memberSvc interface {
	GetByID(
		ctx context.Context,
		memberID uuid.UUID,
	) (domain.Member, error)

	GetByAccountAndOrganization(
		ctx context.Context,
		initiator domain.AccountActor,
		organizationID uuid.UUID,
	) (domain.Member, error)

	GetInitiator(
		ctx context.Context,
		initiator domain.AccountActor,
		organizationID uuid.UUID,
	) (domain.Member, error)

	GetList(
		ctx context.Context,
		filter member.FilterParams,
		limit, offset uint,
	) (pagi.Page[[]domain.Member], error)

	Update(
		ctx context.Context,
		initiator domain.AccountActor,
		memberID uuid.UUID,
		params member.UpdateParams,
	) (domain.Member, error)

	Delete(
		ctx context.Context,
		initiator domain.AccountActor,
		memberID uuid.UUID,
	) error
}

type roleSvc interface {
	Create(
		ctx context.Context,
		initiator domain.AccountActor,
		params role.CreateParams,
	) (domain.Role, error)

	GetByID(
		ctx context.Context,
		roleID uuid.UUID,
	) (domain.Role, error)

	GetWithPermissions(
		ctx context.Context,
		initiator domain.AccountActor,
		roleID uuid.UUID,
	) (domain.Role, domain.OrgRolePermissionsWithDetailsForRole, error)

	GetList(
		ctx context.Context,
		params role.FilterParams,
		limit, offset uint,
	) (pagi.Page[[]domain.Role], error)

	Update(
		ctx context.Context,
		initiator domain.AccountActor,
		roleID uuid.UUID,
		params role.UpdateParams,
	) (domain.Role, error)

	UpdateRanks(
		ctx context.Context,
		initiator domain.AccountActor,
		organizationID uuid.UUID,
		order map[uuid.UUID]uint,
	) error

	Delete(
		ctx context.Context,
		initiator domain.AccountActor,
		roleID uuid.UUID,
	) error

	AddForMember(
		ctx context.Context,
		initiator domain.AccountActor,
		memberID, roleID uuid.UUID,
	) error

	RemoveFromMember(
		ctx context.Context,
		initiator domain.AccountActor,
		memberID, roleID uuid.UUID,
	) error

	GetAllPermissions(ctx context.Context) ([]domain.OrgRolePermission, error)

	UpdatePermissions(
		ctx context.Context,
		initiator domain.AccountActor,
		roleID uuid.UUID,
		permissions role.SetPermissions,
	) (domain.Role, domain.OrgRolePermissionsWithDetailsForRole, error)
}

type responser interface {
	Status(w http.ResponseWriter, status int)
	Render(w http.ResponseWriter, status int, res interface{})
	RenderErr(w http.ResponseWriter, errs ...error)
}

type Modules struct {
	Organization organizationSvc
	Member       memberSvc
	Role         roleSvc
	Invite       inviteSvc
}

type Controller struct {
	modules   *Modules
	responser responser
}

func New(modules *Modules, responser responser) *Controller {
	return &Controller{
		modules:   modules,
		responser: responser,
	}
}
