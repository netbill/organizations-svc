package controller

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/models"
	"github.com/netbill/organizations-svc/internal/core/modules/invite"
	"github.com/netbill/organizations-svc/internal/core/modules/member"
	"github.com/netbill/organizations-svc/internal/core/modules/organization"

	"github.com/netbill/restkit/pagi"
)

type organizationSvc interface {
	Create(
		ctx context.Context,
		actor models.AccountActor,
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
		actor models.AccountActor,
		limit, offset uint,
	) (pagi.Page[[]models.Organization], error)

	CreateOrgUploadMediaLinks(
		ctx context.Context,
		actor models.AccountActor,
		organizationID uuid.UUID,
	) (models.Organization, models.UploadOrgMediaLinks, error)
	Update(
		ctx context.Context,
		actor models.AccountActor,
		organizationID uuid.UUID,
		params organization.UpdateParams,
	) (models.Organization, error)
	DeleteOrgUploadIcon(
		ctx context.Context,
		actor models.AccountActor,
		organizationID uuid.UUID,
		key string,
	) error
	DeleteOrgUploadBanner(
		ctx context.Context,
		actor models.AccountActor,
		organizationID uuid.UUID,
		key string,
	) error

	Activate(
		ctx context.Context,
		actor models.AccountActor,
		organizationID uuid.UUID,
	) (models.Organization, error)
	Deactivate(
		ctx context.Context,
		actor models.AccountActor,
		organizationID uuid.UUID,
	) (models.Organization, error)
	Suspend(
		ctx context.Context,
		organizationID uuid.UUID,
		value bool,
	) (models.Organization, error)

	Delete(
		ctx context.Context,
		actor models.AccountActor,
		organizationID uuid.UUID,
	) error
}

type inviteSvc interface {
	Create(
		ctx context.Context,
		actor models.AccountActor,
		params invite.CreateParams,
	) (models.Invite, error)

	Decline(
		ctx context.Context,
		actor models.AccountActor,
		inviteID uuid.UUID,
	) (models.Invite, error)
	Accept(
		ctx context.Context,
		actor models.AccountActor,
		inviteID uuid.UUID,
	) (models.Invite, error)

	Cancelled(
		ctx context.Context,
		actor models.AccountActor,
		inviteID uuid.UUID,
	) (models.Invite, error)

	GetListForOrganization(
		ctx context.Context,
		actor models.AccountActor,
		organizationID uuid.UUID,
		limit, offset uint,
	) (pagi.Page[[]models.Invite], error)
	GetListForAccount(
		ctx context.Context,
		actor models.AccountActor,
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
		actor models.AccountActor,
		organizationID uuid.UUID,
	) (models.Member, error)

	GetByAccountAndOrgs(
		ctx context.Context,
		actor models.AccountActor,
		organizationIDs []uuid.UUID,
	) ([]models.Member, error)

	GetList(
		ctx context.Context,
		filter member.FilterParams,
		limit, offset uint,
	) (pagi.Page[[]models.Member], error)

	Update(
		ctx context.Context,
		actor models.AccountActor,
		memberID uuid.UUID,
		params member.UpdateParams,
	) (models.Member, error)

	Delete(
		ctx context.Context,
		actor models.AccountActor,
		memberID uuid.UUID,
	) error

	DeleteSelf(
		ctx context.Context,
		actor models.AccountActor,
		orgID uuid.UUID,
	) error
}

type profileSvc interface {
	GetByID(
		ctx context.Context,
		accountID uuid.UUID,
	) (models.Profile, error)

	GetByIDs(
		ctx context.Context,
		accountIDs []uuid.UUID,
	) ([]models.Profile, error)
}

type Modules struct {
	Organization organizationSvc
	Member       memberSvc
	Invite       inviteSvc
	Profile      profileSvc
}

type Controller struct {
	modules *Modules
}

func New(modules *Modules) *Controller {
	return &Controller{
		modules: modules,
	}
}
