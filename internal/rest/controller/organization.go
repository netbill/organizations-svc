package controller

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/models"
	"github.com/netbill/organizations-svc/internal/rest/requests"
	"github.com/netbill/organizations-svc/internal/rest/responses"
	"github.com/netbill/organizations-svc/internal/rest/scope"
	"github.com/netbill/restkit/pagi"
	"github.com/netbill/restkit/problems"
	"github.com/netbill/restkit/render"
)

type organizationModule interface {
	Create(
		ctx context.Context,
		actor models.AccountActor,
		params core.OrganizationCreateParams,
	) (models.Organization, error)

	GetByID(ctx context.Context, organizationID uuid.UUID) (models.Organization, error)
	GetByIDs(ctx context.Context, organizationIDs []uuid.UUID) ([]models.Organization, error)
	GetList(
		ctx context.Context,
		params core.OrganizationFilterParams,
		limit, offset uint,
	) (pagi.Page[[]models.Organization], error)

	Update(
		ctx context.Context,
		actor models.AccountActor,
		organizationID uuid.UUID,
		params core.OrganizationUpdateParams,
	) (models.Organization, error)

	Activate(
		ctx context.Context,
		actor models.AccountActor,
		organizationID uuid.UUID,
		value bool,
	) (models.Organization, error)
	Suspend(
		ctx context.Context,
		organizationID uuid.UUID,
		value bool,
	) (models.Organization, error)

	Delete(ctx context.Context, actor models.AccountActor, organizationID uuid.UUID) error

	CreateUploadMediaLinks(
		ctx context.Context,
		actor models.AccountActor,
		organizationID uuid.UUID,
	) (models.Organization, models.UploadOrgMediaLinks, error)

	DeleteUploadMedia(
		ctx context.Context,
		actor models.AccountActor,
		organizationID uuid.UUID,
		params core.DeleteUploadOrgMediaParams,
	) error
}

type organizationMemberModule interface {
	GetByAccountAndOrgs(
		ctx context.Context,
		actor models.AccountActor,
		organizationIDs []uuid.UUID,
	) ([]models.Member, error)
	GetList(
		ctx context.Context,
		filter core.MemberFilterParams,
		limit, offset uint,
	) (pagi.Page[[]models.Member], error)
}

type organizationProfileModule interface {
	GetByIDs(ctx context.Context, accountIDs []uuid.UUID) ([]models.Profile, error)
}

type OrganizationController struct {
	organizations organizationModule
	members       organizationMemberModule
	profiles      organizationProfileModule
}

type OrganizationControllerDeps struct {
	Organizations organizationModule
	Members       organizationMemberModule
	Profiles      organizationProfileModule
}

func NewOrganizationController(deps OrganizationControllerDeps) *OrganizationController {
	return &OrganizationController{
		organizations: deps.Organizations,
		members:       deps.Members,
		profiles:      deps.Profiles,
	}
}

const operationCreateOrganization = "create_organization"

func (c *OrganizationController) Create(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationCreateOrganization)

	req, err := requests.CreateOrganization(r)
	if err != nil {
		log.WithError(err).Warn("invalid create organization requests")
		render.ResponseError(w, problems.BadRequest(err)...)
		return
	}

	res, err := c.organizations.Create(
		r.Context(),
		scope.AccountActor(r),
		core.OrganizationCreateParams{
			Name: req.Data.Attributes.Name,
		},
	)
	switch {
	case err != nil:
		log.WithError(err).Error("failed to create organization")
		render.ResponseError(w, problems.InternalError())
	default:
		log.Info("organization created successfully")
		render.Response(w, http.StatusCreated, responses.Organization(r, res))
	}
}

const operationGetOrganization = "get_organization"

func (c *OrganizationController) Get(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationGetOrganization)

	organizationID, err := uuid.Parse(chi.URLParam(r, "organization_id"))
	if err != nil {
		log.WithError(err).Warn("invalid organization id")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"path/organization_id": fmt.Errorf(
				"invalid organization id: %s", chi.URLParam(r, "organization_id"),
			),
		})...)
		return
	}

	log = log.WithField("organization_id", organizationID)

	org, err := c.organizations.GetByID(r.Context(), organizationID)
	switch {
	case errors.Is(err, errx.ErrorOrganizationNotExists):
		log.WithError(err).Warn("organization not found")
		render.ResponseError(w, problems.NotFound("organization not found"))
	case err != nil:
		log.WithError(err).Error("failed to get organization")
		render.ResponseError(w, problems.InternalError())
	default:
		render.Response(w, http.StatusOK, responses.Organization(r, org))
	}
}

const operationGetOrganizations = "get_organizations"

func (c *OrganizationController) GetList(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationGetOrganizations)

	q := r.URL.Query()
	params := core.OrganizationFilterParams{}

	if v := strings.TrimSpace(q.Get("text")); v != "" {
		params.Text = &v
	}
	if v := strings.TrimSpace(q.Get("status")); v != "" {
		params.Status = &v
	}

	limit, offset := pagi.GetPagination(r)
	if limit > 100 {
		log.Info("invalid pagination limit")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"query/size": fmt.Errorf("pagination limit must be between 1 and 100"),
		})...)
		return
	}

	log = log.WithField("limit", limit).WithField("offset", offset)

	organizations, err := c.organizations.GetList(r.Context(), params, limit, offset)
	switch {
	case err != nil:
		log.WithError(err).Error("failed to get organizations")
		render.ResponseError(w, problems.InternalError())
	default:
		render.Response(w, http.StatusOK, responses.Organizations(r, organizations))
	}
}

const operationUpdateOrganization = "update_organization"

func (c *OrganizationController) Update(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationUpdateOrganization)

	req, err := requests.UpdateOrganization(r)
	if err != nil {
		log.WithError(err).Warn("invalid update organization requests")
		render.ResponseError(w, problems.BadRequest(err)...)
		return
	}

	log = log.WithField("organization_id", req.Data.Id)

	res, err := c.organizations.Update(
		r.Context(),
		scope.AccountActor(r),
		req.Data.Id,
		core.OrganizationUpdateParams{
			Name:      req.Data.Attributes.Name,
			IconKey:   req.Data.Attributes.IconKey,
			BannerKey: req.Data.Attributes.BannerKey,
		},
	)
	switch {
	case errors.Is(err, errx.ErrorOrganizationNotExists):
		log.WithError(err).Warn("organization not found")
		render.ResponseError(w, problems.NotFound("organization not found"))
	case errors.Is(err, errx.ErrorNotOrganizationHead):
		log.WithError(err).Warn("not enough rights to update organization")
		render.ResponseError(w, problems.Forbidden("not enough rights to update organization"))
	case errors.Is(err, errx.ErrorOrganizationIsSuspended):
		log.WithError(err).Warn("organization is suspended")
		render.ResponseError(w, problems.Forbidden("organization is suspended"))
	case errors.Is(err, errx.ErrorOrganizationUploadedIconInvalid):
		log.WithError(err).Warn("upload icon is invalid")
		render.ResponseError(w, problems.BadRequest(validation.Errors{"icon": err})...)
	case errors.Is(err, errx.ErrorOrganizationUploadedBannerInvalid):
		log.WithError(err).Warn("upload banner is invalid")
		render.ResponseError(w, problems.BadRequest(validation.Errors{"banner": err})...)
	case err != nil:
		log.WithError(err).Error("failed to update organization")
		render.ResponseError(w, problems.InternalError())
	default:
		log.Info("organization updated successfully")
		render.Response(w, http.StatusOK, responses.Organization(r, res))
	}
}

const operationDeleteOrganization = "delete_organization"

func (c *OrganizationController) Delete(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationDeleteOrganization)

	organizationID, err := uuid.Parse(chi.URLParam(r, "organization_id"))
	if err != nil {
		log.WithError(err).Warn("invalid organization id")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"path/organization_id": fmt.Errorf("invalid organization id: %s", chi.URLParam(r, "organization_id")),
		})...)
		return
	}

	log = log.WithField("organization_id", organizationID)

	err = c.organizations.Delete(r.Context(), scope.AccountActor(r), organizationID)
	switch {
	case errors.Is(err, errx.ErrorOrganizationNotExists):
		log.WithError(err).Warn("organization not found")
		render.ResponseError(w, problems.NotFound("organization not found"))
	case errors.Is(err, errx.ErrorOrganizationDeleted):
		log.WithError(err).Warn("organization already deleted")
		render.ResponseError(w, problems.NotFound("organization not found"))
	case errors.Is(err, errx.ErrorOrganizationIsSuspended):
		log.WithError(err).Warn("organization is suspended")
		render.ResponseError(w, problems.Forbidden("organization is suspended"))
	case errors.Is(err, errx.ErrorInitiatorNotMemberOfOrganization):
		log.WithError(err).Warn("membership in organization not found")
		render.ResponseError(w, problems.Forbidden("not a member of the organization"))
	case errors.Is(err, errx.ErrorNotOrganizationHead):
		log.WithError(err).Warn("not enough rights to delete organization")
		render.ResponseError(w, problems.Forbidden("not enough rights to delete organization"))
	case errors.Is(err, errx.ErrorOrganizationHavePlace):
		log.WithError(err).Warn("organization have place and cannot be deleted")
		render.ResponseError(w, problems.Forbidden("organization have place and cannot be deleted"))
	case err != nil:
		log.WithError(err).Error("failed to delete organization")
		render.ResponseError(w, problems.InternalError())
	default:
		log.Info("organization deleted")
		render.Response(w, http.StatusNoContent, nil)
	}
}
