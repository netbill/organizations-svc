package controllers

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
	"github.com/netbill/organizations-svc/internal/core/member"
	"github.com/netbill/organizations-svc/internal/errx"
	"github.com/netbill/organizations-svc/internal/models"
	"github.com/netbill/organizations-svc/internal/rest/requests"
	"github.com/netbill/organizations-svc/internal/rest/responses"
	"github.com/netbill/organizations-svc/internal/rest/scope"
	"github.com/netbill/restkit"
	"github.com/netbill/restkit/pagi"
	"github.com/netbill/restkit/problems"
	"github.com/netbill/restkit/render"
)

type memberCore interface {
	GetByID(ctx context.Context, memberID uuid.UUID) (models.Member, error)
	GetList(ctx context.Context, filter member.FilterParams, limit, offset uint) (pagi.Page[[]models.Member], error)

	Update(ctx context.Context, actor models.AccountActor, memberID uuid.UUID, params member.UpdateParams) (models.Member, error)

	Delete(ctx context.Context, actor models.AccountActor, memberID uuid.UUID) error
}

type organizationGetter interface {
	GetByID(ctx context.Context, organizationID uuid.UUID) (models.Organization, error)
	GetByIDs(ctx context.Context, organizationIDs []uuid.UUID) ([]models.Organization, error)
}

type profileGetter interface {
	GetByID(ctx context.Context, accountID uuid.UUID) (models.Profile, error)
	GetByIDs(ctx context.Context, accountIDs []uuid.UUID) ([]models.Profile, error)
}

type MemberController struct {
	members       memberCore
	profiles      profileGetter
	organizations organizationGetter
}

type MemberControllerDeps struct {
	Members       memberCore
	Profiles      profileGetter
	Organizations organizationGetter
}

func NewMemberController(deps MemberControllerDeps) *MemberController {
	return &MemberController{
		members:       deps.Members,
		profiles:      deps.Profiles,
		organizations: deps.Organizations,
	}
}

const operationGetMember = "get_member"

func (c *MemberController) Get(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationGetMember)

	memberID, err := uuid.Parse(chi.URLParam(r, "member_id"))
	if err != nil {
		log.WithError(err).Warn("invalid member id")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"path/member_id": fmt.Errorf("invalid member id: %s", chi.URLParam(r, "member_id")),
		})...)
		return
	}

	log = log.WithField("member_id", memberID)

	res, err := c.members.GetByID(r.Context(), memberID)
	switch {
	case errors.Is(err, errx.ErrorMemberNotExists),
		errors.Is(err, errx.ErrorMemberDeleted):
		log.WithError(err).Warn("member not found")
		render.ResponseError(w, problems.NotFound("member not found"))
		return
	case err != nil:
		log.WithError(err).Error("failed to get member by id")
		render.ResponseError(w, problems.InternalError())
		return
	}

	includes := restkit.ParseIncludes(r)
	opts := make([]responses.MemberOption, 0)

	if slices.Contains(includes, "profile") {
		profile, err := c.profiles.GetByID(r.Context(), res.AccountID)
		if err != nil {
			log.WithError(err).WithField("account_id", res.AccountID).
				Error("failed to get profile for member")
		}

		opts = append(opts, responses.WithMemberProfile(r, profile))
	}

	if slices.Contains(includes, "organizations") {
		org, err := c.organizations.GetByID(r.Context(), res.OrganizationID)
		if err != nil {
			log.WithError(err).WithField("organization_id", res.OrganizationID).
				Error("failed to get organization for member")
		}

		opts = append(opts, responses.WithMemberOrganization(r, org))
	}

	render.Response(w, http.StatusOK, responses.Member(r, res, opts...))
}

const operationGetOrganizationMembers = "get_organization_members"

func (c *MemberController) GetList(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationGetOrganizationMembers)

	limit, offset := pagi.GetPagination(r)
	if limit > 100 {
		log.Warn("invalid pagination limit")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"query/size": fmt.Errorf("pagination limit cannot be greater than 100"),
		})...)
		return
	}

	log = log.WithField("limit", limit).WithField("offset", offset)

	params := member.FilterParams{}

	if v := r.URL.Query().Get("organization_id"); v != "" {
		id, err := uuid.Parse(strings.TrimSpace(v))
		if err != nil {
			render.ResponseError(w, problems.BadRequest(validation.Errors{
				"query/organization_id": fmt.Errorf("invalid uuid: %s", v),
			})...)
			return
		}
		params.OrganizationID = &id
	}

	if v := r.URL.Query().Get("account_id"); v != "" {
		id, err := uuid.Parse(strings.TrimSpace(v))
		if err != nil {
			render.ResponseError(w, problems.BadRequest(validation.Errors{
				"query/account_id": fmt.Errorf("invalid uuid: %s", v),
			})...)
			return
		}
		params.AccountID = &id
	}

	if v := r.URL.Query().Get("text"); v != "" {
		text := v
		params.BestMatch = &text
	}

	members, err := c.members.GetList(r.Context(), params, limit, offset)
	switch {
	case err != nil:
		log.WithError(err).Error("failed to get organization members")
		render.ResponseError(w, problems.InternalError())
		return
	}

	opts := make([]responses.MembersCollectionOption, 0, 2)
	includes := restkit.ParseIncludes(r)

	if slices.Contains(includes, "profile") {
		profileIDs := make([]uuid.UUID, 0, members.Size)
		for _, m := range members.Data {
			profileIDs = append(profileIDs, m.AccountID)
		}

		profiles, err := c.profiles.GetByIDs(r.Context(), profileIDs)
		if err != nil {
			log.WithError(err).Error("failed to get profiles for members")
		}

		opts = append(opts, responses.WithCollectionMembersProfiles(r, profiles))
	}

	if slices.Contains(includes, "organization") && params.OrganizationID != nil {
		organization, err := c.organizations.GetByID(r.Context(), *params.OrganizationID)
		if err != nil {
			log.WithError(err).Error("failed to get organization for members")
		}

		opts = append(opts, responses.WithCollectionMembersOrganization(r, organization))
	}

	render.Response(w, http.StatusOK, responses.Members(r, members, opts...))
}

const operationUpdateMember = "update_member"

func (c *MemberController) Update(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationUpdateMember)

	req, err := requests.UpdateMember(r)
	if err != nil {
		log.WithError(err).Warn("invalid update member requests")
		render.ResponseError(w, problems.BadRequest(err)...)
		return
	}

	memberID, err := uuid.Parse(chi.URLParam(r, "member_id"))
	if err != nil {
		log.WithError(err).Warn("invalid member id")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"path/member_id": fmt.Errorf("invalid member id: %s", chi.URLParam(r, "member_id")),
		})...)
		return
	}

	log = log.WithField("member_id", memberID)

	res, err := c.members.Update(
		r.Context(),
		scope.AccountActor(r),
		memberID,
		member.UpdateParams{
			Position: req.Data.Attributes.Position,
			Label:    req.Data.Attributes.Label,
		},
	)
	switch {
	case errors.Is(err, errx.ErrorMemberNotExists),
		errors.Is(err, errx.ErrorMemberDeleted):
		log.WithError(err).Warn("member not found")
		render.ResponseError(w, problems.NotFound("member not found"))
	case errors.Is(err, errx.ErrorNotOrganizationHead):
		log.WithError(err).Warn("not enough rights to update member")
		render.ResponseError(w, problems.Forbidden("not enough rights to update member"))
	case errors.Is(err, errx.ErrorOrganizationNotExists),
		errors.Is(err, errx.ErrorOrganizationDeleted):
		log.WithError(err).Warn("organization not found")
		render.ResponseError(w, problems.NotFound("organization not found"))
	case errors.Is(err, errx.ErrorOrganizationIsSuspended):
		log.WithError(err).Warn("organization is suspended")
		render.ResponseError(w, problems.Forbidden("organization is suspended"))
	case err != nil:
		log.WithError(err).Error("failed to update member")
		render.ResponseError(w, problems.InternalError())
	default:
		log.Info("member updated")
		render.Response(w, http.StatusOK, responses.Member(r, res))
	}
}

const operationDeleteMember = "delete_member"

func (c *MemberController) Delete(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationDeleteMember)

	memberID, err := uuid.Parse(chi.URLParam(r, "member_id"))
	if err != nil {
		log.WithError(err).Warn("invalid member id")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"path/member_id": fmt.Errorf("invalid member id: %s", chi.URLParam(r, "member_id")),
		})...)
		return
	}

	log = log.WithField("member_id", memberID)

	err = c.members.Delete(r.Context(), scope.AccountActor(r), memberID)
	switch {
	case errors.Is(err, errx.ErrorMemberNotExists),
		errors.Is(err, errx.ErrorMemberDeleted):
		log.WithError(err).Warn("member not found")
		render.Response(w, http.StatusNoContent, nil)
	case errors.Is(err, errx.ErrorInitiatorNotMemberOfOrganization):
		log.WithError(err).Warn("membership in organization not found")
		render.ResponseError(w, problems.Forbidden("not a member of the organization"))
	case errors.Is(err, errx.ErrorNotOrganizationHead):
		log.WithError(err).Warn("not enough rights to delete member")
		render.ResponseError(w, problems.Forbidden("not enough rights to delete member"))
	case errors.Is(err, errx.ErrorOrganizationNotExists),
		errors.Is(err, errx.ErrorOrganizationDeleted):
		log.WithError(err).Warn("organization not found")
		render.ResponseError(w, problems.NotFound("organization not found"))
	case errors.Is(err, errx.ErrorOrganizationIsSuspended):
		log.WithError(err).Warn("organization is suspended")
		render.ResponseError(w, problems.Forbidden("organization is suspended"))
	case errors.Is(err, errx.ErrorCannotDeleteOrganizationHeadMember):
		log.WithError(err).Warn("cannot delete organization head member")
		render.ResponseError(w, problems.Forbidden("cannot delete organization head member"))
	case err != nil:
		log.WithError(err).Error("failed to delete member")
		render.ResponseError(w, problems.InternalError())
	default:
		log.Info("member deleted")
		render.Response(w, http.StatusNoContent, nil)
	}
}
