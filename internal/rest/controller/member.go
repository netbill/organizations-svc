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

type memberModule interface {
	GetByID(ctx context.Context, memberID uuid.UUID) (models.Member, error)
	GetList(ctx context.Context, filter core.MemberFilterParams, limit, offset uint) (pagi.Page[[]models.Member], error)
	Update(ctx context.Context, actor models.AccountActor, memberID uuid.UUID, params core.MemberUpdateParams) (models.Member, error)
	Delete(ctx context.Context, actor models.AccountActor, memberID uuid.UUID) error
	DeleteSelf(ctx context.Context, actor models.AccountActor, orgID uuid.UUID) error
}

type memberProfileModule interface {
	GetByID(ctx context.Context, accountID uuid.UUID) (models.Profile, error)
	GetByIDs(ctx context.Context, accountIDs []uuid.UUID) ([]models.Profile, error)
}

type memberOrganizationModule interface {
	GetByID(ctx context.Context, organizationID uuid.UUID) (models.Organization, error)
}

type MemberRouter struct {
	members       memberModule
	profiles      memberProfileModule
	organizations memberOrganizationModule
}

func NewMemberRouter(
	members memberModule,
	profiles memberProfileModule,
	organizations memberOrganizationModule,
) *MemberRouter {
	return &MemberRouter{
		members:       members,
		profiles:      profiles,
		organizations: organizations,
	}
}

const operationGetMember = "get_member"

func (c *MemberRouter) Get(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationGetMember)

	memberID, err := uuid.Parse(chi.URLParam(r, "member_id"))
	if err != nil {
		log.WithError(err).Warn("invalid member id")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"query": fmt.Errorf("invalid member id: %s", chi.URLParam(r, "member_id")),
		})...)
		return
	}

	log = log.WithField("member_id", memberID)

	member, err := c.members.GetByID(r.Context(), memberID)
	switch {
	case errors.Is(err, errx.ErrorMemberNotExists):
		log.WithError(err).Warn("member not found")
		render.ResponseError(w, problems.NotFound("member not found"))
		return
	case err != nil:
		log.WithError(err).Error("failed to get member by id")
		render.ResponseError(w, problems.InternalError())
		return
	}

	includes := r.URL.Query()["include"]
	opts := make([]responses.MemberOption, 0)

	if slices.Contains(includes, "profile") {
		profile, err := c.profiles.GetByID(r.Context(), member.AccountID)
		switch {
		case errors.Is(err, errx.ErrorProfileNotExists):
			log.WithField("account_id", member.AccountID).Warn("profile not found")
			render.ResponseError(w, problems.NotFound("profile not found"))
			return
		case err != nil:
			log.WithError(err).WithField("account_id", member.AccountID).Error("failed to get profile")
			render.ResponseError(w, problems.InternalError())
			return
		default:
			opts = append(opts, responses.WithMemberProfile(profile))
		}
	}

	if slices.Contains(includes, "organizations") {
		org, err := c.organizations.GetByID(r.Context(), member.OrganizationID)
		switch {
		case errors.Is(err, errx.ErrorOrganizationNotExists):
			log.WithError(err).Warn("organization not found")
			render.ResponseError(w, problems.NotFound("organization not found"))
			return
		case err != nil:
			log.WithError(err).Error("failed to get organization")
			render.ResponseError(w, problems.InternalError())
			return
		default:
			opts = append(opts, responses.WithMemberOrganization(org))
		}
	}

	render.Response(w, http.StatusOK, responses.Member(member, opts...))
}

const operationGetOrganizationMembers = "get_organization_members"

func (c *MemberRouter) GetList(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationGetOrganizationMembers)

	limit, offset := pagi.GetPagination(r)
	if limit > 100 {
		log.Warn("invalid pagination limit")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"query": fmt.Errorf("pagination limit cannot be greater than 100"),
		})...)
		return
	}

	log = log.WithField("limit", limit).WithField("offset", offset)

	params := core.MemberFilterParams{}

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
	includesRaw := r.URL.Query()["include"]
	includes := make([]string, 0, 2)

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

	if slices.Contains(includes, "profile") {
		profileIDs := make([]uuid.UUID, 0, members.Size)
		for _, m := range members.Data {
			profileIDs = append(profileIDs, m.AccountID)
		}

		profiles, err := c.profiles.GetByIDs(r.Context(), profileIDs)
		if err != nil {
			log.WithError(err).Error("failed to get profiles for members")
			render.ResponseError(w, problems.InternalError())
			return
		}
		opts = append(opts, responses.WithCollectionMembersProfiles(profiles))
	}

	if slices.Contains(includes, "organization") && params.OrganizationID != nil {
		organization, err := c.organizations.GetByID(r.Context(), *params.OrganizationID)
		if err != nil {
			log.WithError(err).Error("failed to get organization for members")
			render.ResponseError(w, problems.InternalError())
			return
		}
		opts = append(opts, responses.WithCollectionMembersOrganization(organization))
	}

	render.Response(w, http.StatusOK, responses.Members(r, members, opts...))
}

const operationUpdateMember = "update_member"

func (c *MemberRouter) Update(w http.ResponseWriter, r *http.Request) {
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
			"query": fmt.Errorf("invalid member id: %s", chi.URLParam(r, "member_id")),
		})...)
		return
	}

	log = log.WithField("member_id", memberID)

	res, err := c.members.Update(
		r.Context(),
		scope.AccountActor(r),
		memberID,
		core.MemberUpdateParams{
			Position: req.Data.Attributes.Position,
			Label:    req.Data.Attributes.Label,
		},
	)
	switch {
	case errors.Is(err, errx.ErrorMemberNotExists):
		log.WithError(err).Warn("member not found")
		render.ResponseError(w, problems.NotFound("member not found"))
	case errors.Is(err, errx.ErrorNotOrganizationHead):
		log.WithError(err).Warn("not enough rights to update member")
		render.ResponseError(w, problems.Forbidden("not enough rights to update member"))
	case errors.Is(err, errx.ErrorOrganizationNotExists):
		log.WithError(err).Warn("organization not found")
		render.ResponseError(w, problems.NotFound("organization not found"))
	case errors.Is(err, errx.ErrorOrganizationIsSuspended):
		log.WithError(err).Warn("organization is suspended")
		render.ResponseError(w, problems.Forbidden("organization is suspended"))
	case err != nil:
		log.WithError(err).Error("failed to update member")
		render.ResponseError(w, problems.InternalError())
	default:
		render.Response(w, http.StatusOK, responses.Member(res))
	}
}

const operationDeleteMember = "delete_member"

func (c *MemberRouter) Delete(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationDeleteMember)

	memberID, err := uuid.Parse(chi.URLParam(r, "member_id"))
	if err != nil {
		log.WithError(err).Warn("invalid member id")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"query": fmt.Errorf("invalid member id: %s", chi.URLParam(r, "member_id")),
		})...)
		return
	}

	log = log.WithField("member_id", memberID)

	err = c.members.Delete(r.Context(), scope.AccountActor(r), memberID)
	switch {
	case errors.Is(err, errx.ErrorCannotDeleteSelf):
		log.WithError(err).Warn("cannot delete self")
		render.ResponseError(w, problems.Forbidden("cannot delete self"))
	case errors.Is(err, errx.ErrorMemberNotExists):
		log.WithError(err).Warn("member not found")
		render.ResponseError(w, problems.NotFound("member not found"))
	case errors.Is(err, errx.ErrorMemberDeleted):
		log.WithError(err).Warn("member already deleted")
		render.Response(w, http.StatusNoContent, nil)
	case errors.Is(err, errx.ErrorInitiatorNotMemberOfOrganization):
		log.WithError(err).Warn("membership in organization not found")
		render.ResponseError(w, problems.Forbidden("not a member of the organization"))
	case errors.Is(err, errx.ErrorNotOrganizationHead):
		log.WithError(err).Warn("not enough rights to delete member")
		render.ResponseError(w, problems.Forbidden("not enough rights to delete member"))
	case errors.Is(err, errx.ErrorOrganizationNotExists):
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

const operationLeaveOrganization = "leave_organization"

func (c *MemberRouter) LeaveFromOrg(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationLeaveOrganization)

	organizationID, err := uuid.Parse(chi.URLParam(r, "organization_id"))
	if err != nil {
		log.WithError(err).Warn("invalid organization id")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"query": fmt.Errorf("invalid organization id: %s", chi.URLParam(r, "organization_id")),
		})...)
		return
	}

	log = log.WithField("organization_id", organizationID)

	err = c.members.DeleteSelf(r.Context(), scope.AccountActor(r), organizationID)
	switch {
	case errors.Is(err, errx.ErrorOrganizationNotExists):
		log.WithError(err).Warn("organization not found")
		render.ResponseError(w, problems.NotFound("organization not found"))
	case errors.Is(err, errx.ErrorInitiatorNotMemberOfOrganization):
		log.WithError(err).Warn("account is not a member of organization")
		render.ResponseError(w, problems.Forbidden("account is not a member of organization"))
	case errors.Is(err, errx.ErrorCannotDeleteOrganizationHeadMember):
		log.WithError(err).Warn("cannot leave organization as head member")
		render.ResponseError(w, problems.Forbidden("cannot leave organization as head member"))
	case err != nil:
		log.WithError(err).Error("failed to leave organization")
		render.ResponseError(w, problems.InternalError())
	default:
		log.Info("left organization")
		render.Response(w, http.StatusNoContent, nil)
	}
}
