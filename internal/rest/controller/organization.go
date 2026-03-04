package controller

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"slices"
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
	Create(ctx context.Context, actor models.AccountActor, params core.OrganizationCreateParams) (models.Organization, error)
	GetByID(ctx context.Context, organizationID uuid.UUID) (models.Organization, error)
	GetByIDs(ctx context.Context, organizationIDs []uuid.UUID) ([]models.Organization, error)
	GetList(ctx context.Context, params core.OrganizationFilterParams, limit, offset uint) (pagi.Page[[]models.Organization], error)
	GetForUser(ctx context.Context, actor models.AccountActor, limit, offset uint) (pagi.Page[[]models.Organization], error)
	Update(ctx context.Context, actor models.AccountActor, organizationID uuid.UUID, params core.OrganizationUpdateParams) (models.Organization, error)
	Delete(ctx context.Context, actor models.AccountActor, organizationID uuid.UUID) error
	Activate(ctx context.Context, actor models.AccountActor, organizationID uuid.UUID) (models.Organization, error)
	Deactivate(ctx context.Context, actor models.AccountActor, organizationID uuid.UUID) (models.Organization, error)
	Suspend(ctx context.Context, organizationID uuid.UUID, value bool) (models.Organization, error)
	CreateOrgUploadMediaLinks(ctx context.Context, actor models.AccountActor, organizationID uuid.UUID) (models.Organization, models.UploadOrgMediaLinks, error)
	DeleteOrgUploadIcon(ctx context.Context, actor models.AccountActor, organizationID uuid.UUID, key string) error
	DeleteOrgUploadBanner(ctx context.Context, actor models.AccountActor, organizationID uuid.UUID, key string) error
}

type organizationMemberModule interface {
	GetByAccountAndOrgs(ctx context.Context, actor models.AccountActor, organizationIDs []uuid.UUID) ([]models.Member, error)
	DeleteSelf(ctx context.Context, actor models.AccountActor, orgID uuid.UUID) error
	GetList(ctx context.Context, filter core.MemberFilterParams, limit, offset uint) (pagi.Page[[]models.Member], error)
}

type organizationInviteModule interface {
	GetListForOrganization(ctx context.Context, actor models.AccountActor, organizationID uuid.UUID, limit, offset uint) (pagi.Page[[]models.Invite], error)
}

type organizationProfileModule interface {
	GetByIDs(ctx context.Context, accountIDs []uuid.UUID) ([]models.Profile, error)
}

type Organization struct {
	organizations organizationModule
	members       organizationMemberModule
	invites       organizationInviteModule
	profiles      organizationProfileModule
}

func NewOrganizationRouter(
	organizations organizationModule,
	members organizationMemberModule,
	invites organizationInviteModule,
	profiles organizationProfileModule,
) *Organization {
	return &Organization{
		organizations: organizations,
		members:       members,
		invites:       invites,
		profiles:      profiles,
	}
}

const operationCreateOrganization = "create_organization"

func (c *Organization) Create(w http.ResponseWriter, r *http.Request) {
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
		render.Response(w, http.StatusCreated, responses.Organization(res))
	}
}

const operationGetMyOrganizations = "get_my_organizations"

func (c *Organization) GetMyList(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationGetMyOrganizations)

	limit, offset := pagi.GetPagination(r)
	if limit > 100 {
		log.Warn("invalid pagination limit")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"query": fmt.Errorf("pagination limit cannot be greater than 100"),
		})...)
		return
	}

	accountID := scope.AccountActor(r)
	log = log.WithField("account_id", accountID).WithField("limit", limit).WithField("offset", offset)

	res, err := c.organizations.GetForUser(r.Context(), accountID, limit, offset)
	switch {
	case err != nil:
		log.WithError(err).Error("failed to get organizations")
		render.ResponseError(w, problems.InternalError())
		return
	}

	opts := make([]responses.OrgCollectionOption, 0)
	includesRaw := r.URL.Query()["include"]
	includes := make([]string, 0, 1)

	for _, v := range includesRaw {
		for _, part := range strings.Split(v, ",") {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}
			if !slices.Contains(includes, part) {
				includes = append(includes, part)
			}
		}
	}

	if slices.Contains(includes, "members") {
		organizationIDs := make([]uuid.UUID, res.Size)
		for i, organization := range res.Data {
			organizationIDs[i] = organization.ID
		}

		members, err := c.members.GetByAccountAndOrgs(r.Context(), accountID, organizationIDs)
		if err != nil {
			log.WithError(err).Error("failed to get members for organizations")
			render.ResponseError(w, problems.InternalError())
			return
		}

		opts = append(opts, responses.WithOrganizationMembers(members))
	}

	render.Response(w, http.StatusOK, responses.Organizations(r, res, opts...))
}

const operationGetOrganization = "get_organization"

func (c *Organization) Get(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationGetOrganization)

	organizationID, err := uuid.Parse(chi.URLParam(r, "organization_id"))
	if err != nil {
		log.WithError(err).Warn("invalid organization id")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"query": fmt.Errorf("invalid organization id: %s", chi.URLParam(r, "organization_id")),
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
		render.Response(w, http.StatusOK, responses.Organization(org))
	}
}

const operationGetOrganizations = "get_organizations"

func (c *Organization) GetList(w http.ResponseWriter, r *http.Request) {
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
	if limit < 1 || limit > 100 {
		log.Info("invalid pagination limit")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"query": fmt.Errorf("pagination limit must be between 1 and 100"),
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

func (c *Organization) Update(w http.ResponseWriter, r *http.Request) {
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
	case errors.Is(err, errx.ErrorOrganizationIconKeyIsInvalid):
		log.WithError(err).Warn("icon key is invalid")
		render.ResponseError(w, problems.BadRequest(validation.Errors{"icon": err})...)
	case errors.Is(err, errx.ErrorOrganizationIconFormatIsNotAllowed):
		log.WithError(err).Warn("icon format is not allowed")
		render.ResponseError(w, problems.BadRequest(validation.Errors{"icon": err})...)
	case errors.Is(err, errx.ErrorOrganizationIconContentIsExceedsMax):
		log.WithError(err).Warn("icon content is exceeds max")
		render.ResponseError(w, problems.BadRequest(validation.Errors{"icon": err})...)
	case errors.Is(err, errx.ErrorOrganizationIconResolutionIsInvalid):
		log.WithError(err).Warn("icon resolution is invalid")
		render.ResponseError(w, problems.BadRequest(validation.Errors{"icon": err})...)
	case errors.Is(err, errx.ErrorOrganizationBannerKeyIsInvalid):
		log.WithError(err).Warn("banner key is invalid")
		render.ResponseError(w, problems.BadRequest(validation.Errors{"banner": err})...)
	case errors.Is(err, errx.ErrorOrganizationBannerFormatIsNotAllowed):
		log.WithError(err).Warn("banner format is not allowed")
		render.ResponseError(w, problems.BadRequest(validation.Errors{"banner": err})...)
	case errors.Is(err, errx.ErrorOrganizationBannerContentIsExceedsMax):
		log.WithError(err).Warn("banner content is exceeds max")
		render.ResponseError(w, problems.BadRequest(validation.Errors{"banner": err})...)
	case errors.Is(err, errx.ErrorOrganizationBannerResolutionIsInvalid):
		log.WithError(err).Warn("banner resolution is invalid")
		render.ResponseError(w, problems.BadRequest(validation.Errors{"banner": err})...)
	case err != nil:
		log.WithError(err).Error("failed to update organization")
		render.ResponseError(w, problems.InternalError())
	default:
		log.Info("organization updated successfully")
		render.Response(w, http.StatusOK, responses.Organization(res))
	}
}

const operationActivateOrganization = "activate_organization"

func (c *Organization) Activate(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationActivateOrganization)

	organizationID, err := uuid.Parse(chi.URLParam(r, "organization_id"))
	if err != nil {
		log.WithError(err).Warn("invalid organization id")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"query": fmt.Errorf("invalid organization id: %s", chi.URLParam(r, "organization_id")),
		})...)
		return
	}

	log = log.WithField("organization_id", organizationID)

	res, err := c.organizations.Activate(r.Context(), scope.AccountActor(r), organizationID)
	switch {
	case errors.Is(err, errx.ErrorOrganizationNotExists):
		log.WithError(err).Warn("organization not found")
		render.ResponseError(w, problems.NotFound("organization not found"))
	case errors.Is(err, errx.ErrorOrganizationIsSuspended):
		log.WithError(err).Warn("organization is suspended")
		render.ResponseError(w, problems.Forbidden("organization is suspended"))
	case errors.Is(err, errx.ErrorInitiatorNotMemberOfOrganization):
		log.WithError(err).Warn("membership in organization not found")
		render.ResponseError(w, problems.Forbidden("not a member of the organization"))
	case errors.Is(err, errx.ErrorNotOrganizationHead):
		log.WithError(err).Warn("only organization head can activate organization")
		render.ResponseError(w, problems.Forbidden("only organization head can activate organization"))
	case err != nil:
		log.WithError(err).Error("failed to activate organization")
		render.ResponseError(w, problems.InternalError())
	default:
		render.Response(w, http.StatusOK, responses.Organization(res))
	}
}

const operationDeactivateOrganization = "deactivate_organization"

func (c *Organization) Deactivate(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationDeactivateOrganization)

	organizationID, err := uuid.Parse(chi.URLParam(r, "organization_id"))
	if err != nil {
		log.WithError(err).Warn("invalid organization id")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"query": fmt.Errorf("invalid organization id: %s", chi.URLParam(r, "organization_id")),
		})...)
		return
	}

	log = log.WithField("organization_id", organizationID)

	res, err := c.organizations.Deactivate(r.Context(), scope.AccountActor(r), organizationID)
	switch {
	case errors.Is(err, errx.ErrorOrganizationNotExists):
		log.WithError(err).Warn("organization not found")
		render.ResponseError(w, problems.NotFound("organization not found"))
	case errors.Is(err, errx.ErrorInitiatorNotMemberOfOrganization):
		log.WithError(err).Warn("membership in organization not found")
		render.ResponseError(w, problems.Forbidden("not a member of the organization"))
	case errors.Is(err, errx.ErrorNotOrganizationHead):
		log.WithError(err).Warn("not enough rights to update organization")
		render.ResponseError(w, problems.Forbidden("not enough rights to update organization"))
	case errors.Is(err, errx.ErrorOrganizationIsSuspended):
		log.WithError(err).Warn("organization is suspended")
		render.ResponseError(w, problems.Forbidden("organization is suspended"))
	case err != nil:
		log.WithError(err).Error("failed to deactivate organization")
		render.ResponseError(w, problems.InternalError())
	default:
		render.Response(w, http.StatusOK, responses.Organization(res))
	}
}

const operationSuspendOrganization = "suspend_organization"

func (c *Organization) Suspend(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationSuspendOrganization)

	organizationID, err := uuid.Parse(chi.URLParam(r, "organization_id"))
	if err != nil {
		log.WithError(err).Warn("invalid organization id")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"query": fmt.Errorf("invalid organization id: %s", chi.URLParam(r, "organization_id")),
		})...)
		return
	}

	log = log.WithField("organization_id", organizationID)

	org, err := c.organizations.Suspend(r.Context(), organizationID, true)
	switch {
	case errors.Is(err, errx.ErrorOrganizationNotExists):
		log.WithError(err).Warn("organization not found")
		render.ResponseError(w, problems.NotFound("organization not found"))
	case err != nil:
		log.WithError(err).Error("failed to suspend organization")
		render.ResponseError(w, problems.InternalError())
	default:
		render.Response(w, http.StatusNoContent, responses.Organization(org))
	}
}

const operationUnsuspendOrganization = "unsuspend_organization"

func (c *Organization) Unsuspend(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationUnsuspendOrganization)

	organizationID, err := uuid.Parse(chi.URLParam(r, "organization_id"))
	if err != nil {
		log.WithError(err).Warn("invalid organization id")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"query": fmt.Errorf("invalid organization id: %s", chi.URLParam(r, "organization_id")),
		})...)
		return
	}

	log = log.WithField("organization_id", organizationID)

	org, err := c.organizations.Suspend(r.Context(), organizationID, false)
	switch {
	case errors.Is(err, errx.ErrorOrganizationNotExists):
		log.WithError(err).Warn("organization not found")
		render.ResponseError(w, problems.NotFound("organization not found"))
	case err != nil:
		log.WithError(err).Error("failed to unsuspend organization")
		render.ResponseError(w, problems.InternalError())
	default:
		render.Response(w, http.StatusNoContent, responses.Organization(org))
	}
}

const operationDeleteOrganization = "delete_organization"

func (c *Organization) Delete(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationDeleteOrganization)

	orgID, err := uuid.Parse(chi.URLParam(r, "organization_id"))
	if err != nil {
		log.WithError(err).Warn("invalid organization id")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"query": fmt.Errorf("invalid organization id: %s", chi.URLParam(r, "organization_id")),
		})...)
		return
	}

	log = log.WithField("organization_id", orgID)

	err = c.organizations.Delete(r.Context(), scope.AccountActor(r), orgID)
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

const operationCreateOrganizationUploadMediaLink = "create_organization_upload_media_link"

func (c *Organization) CreateUploadMediaLink(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationCreateOrganizationUploadMediaLink)

	organizationID, err := uuid.Parse(chi.URLParam(r, "organization_id"))
	if err != nil {
		log.WithError(err).Warn("invalid organization id")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"query": fmt.Errorf("invalid organization id: %s", chi.URLParam(r, "organization_id")),
		})...)
		return
	}

	org, media, err := c.organizations.CreateOrgUploadMediaLinks(r.Context(), scope.AccountActor(r), organizationID)
	switch {
	case errors.Is(err, errx.ErrorOrganizationNotExists):
		log.WithError(err).Warn("organization does not exist")
		render.ResponseError(w, problems.NotFound("organization does not exist"))
	case errors.Is(err, errx.ErrorNotOrganizationHead):
		log.WithError(err).Warn("only organization head can create organization upload media link")
		render.ResponseError(w, problems.Forbidden("not enough rights to create organization upload media link"))
	case errors.Is(err, errx.ErrorInitiatorNotMemberOfOrganization):
		log.WithError(err).Warn("membership in organization not found")
		render.ResponseError(w, problems.Forbidden("not a member of the organization"))
	case errors.Is(err, errx.ErrorOrganizationIsSuspended):
		log.WithError(err).Warn("organization is suspended")
		render.ResponseError(w, problems.Forbidden("organization is suspended"))
	case err != nil:
		log.WithError(err).Error("failed to create organization upload media link")
		render.ResponseError(w, problems.InternalError())
	default:
		render.Response(w, http.StatusOK, responses.UploadOrganizationMediaLinks(org, media))
	}
}

const operationDeleteUploadOrganizationIcon = "delete_upload_organization_icon"

func (c *Organization) DeleteUploadIcon(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationDeleteUploadOrganizationIcon)

	req, err := requests.DeleteUploadOrgIcon(r)
	if err != nil {
		log.WithError(err).Warn("invalid delete upload organization icon requests")
		render.ResponseError(w, problems.BadRequest(err)...)
		return
	}

	log = log.WithField("organization_id", req.Data.Id)

	err = c.organizations.DeleteOrgUploadIcon(
		r.Context(),
		scope.AccountActor(r),
		req.Data.Id,
		req.Data.Attributes.IconKey,
	)
	switch {
	case errors.Is(err, errx.ErrorOrganizationNotExists):
		log.WithError(err).Warn("organization does not exist")
		render.ResponseError(w, problems.NotFound("organization does not exist"))
	case errors.Is(err, errx.ErrorOrganizationIsSuspended):
		log.WithError(err).Warn("organization is suspended")
		render.ResponseError(w, problems.Forbidden("organization is suspended"))
	case errors.Is(err, errx.ErrorInitiatorNotMemberOfOrganization):
		log.WithError(err).Warn("member not found")
		render.ResponseError(w, problems.Forbidden("member not found"))
	case errors.Is(err, errx.ErrorNotOrganizationHead):
		log.WithError(err).Warn("not enough rights to update organization")
		render.ResponseError(w, problems.Forbidden("not enough rights to update organization"))
	case err != nil:
		log.WithError(err).Error("failed to delete organization icon in upload session")
		render.ResponseError(w, problems.InternalError())
	default:
		render.Response(w, http.StatusNoContent, nil)
	}
}

const operationDeleteUploadOrganizationBanner = "delete_upload_organization_banner"

func (c *Organization) DeleteUploadBanner(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationDeleteUploadOrganizationBanner)

	req, err := requests.DeleteUploadOrgBanner(r)
	if err != nil {
		log.WithError(err).Warn("invalid delete upload organization banner requests")
		render.ResponseError(w, problems.BadRequest(err)...)
		return
	}

	log = log.WithField("organization_id", req.Data.Id)

	err = c.organizations.DeleteOrgUploadBanner(
		r.Context(),
		scope.AccountActor(r),
		req.Data.Id,
		req.Data.Attributes.BannerKey,
	)
	switch {
	case errors.Is(err, errx.ErrorOrganizationNotExists):
		log.WithError(err).Warn("organization does not exist")
		render.ResponseError(w, problems.NotFound("organization does not exist"))
	case errors.Is(err, errx.ErrorOrganizationIsSuspended):
		log.WithError(err).Warn("organization is suspended")
		render.ResponseError(w, problems.Forbidden("organization is suspended"))
	case errors.Is(err, errx.ErrorInitiatorNotMemberOfOrganization):
		log.WithError(err).Warn("member not found")
		render.ResponseError(w, problems.Forbidden("member not found"))
	case errors.Is(err, errx.ErrorNotOrganizationHead):
		log.WithError(err).Warn("not enough rights to update organization")
		render.ResponseError(w, problems.Forbidden("not enough rights to update organization"))
	case err != nil:
		log.WithError(err).Error("failed to delete organization banner in upload session")
		render.ResponseError(w, problems.InternalError())
	default:
		render.Response(w, http.StatusNoContent, nil)
	}
}
