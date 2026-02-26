package controller

import (
	"errors"
	"fmt"
	"net/http"
	"slices"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/rest/responses"
	"github.com/netbill/organizations-svc/internal/rest/scope"
	"github.com/netbill/restkit/problems"
	"github.com/netbill/restkit/render"
)

const operationGetMember = "get_member"

func (c *Controller) GetMember(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationGetMember)

	memberID, err := uuid.Parse(chi.URLParam(r, "member_id"))
	if err != nil {
		log.WithError(err).Info("invalid member id")
		render.ResponseError(w, problems.BadRequest(fmt.Errorf("invalid member id"))...)
		return
	}

	log = log.WithField("member_id", memberID)

	member, err := c.modules.Member.GetByID(r.Context(), memberID)
	switch {
	case errors.Is(err, errx.ErrorMemberNotFound):
		log.Info("member not found")
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
		profile, err := c.modules.Profile.GetByID(r.Context(), member.AccountID)
		switch {
		case errors.Is(err, errx.ErrorProfileNotFound):
			log.WithField("account_id", member.AccountID).Info("profile not found")
			render.ResponseError(w, problems.NotFound("profile not found"))
			return
		case err != nil:
			log.WithField("account_id", member.AccountID).WithError(err).Error("failed to get profile")
			render.ResponseError(w, problems.InternalError())
			return
		default:
			opts = append(opts, responses.WithMemberProfile(profile))
		}
	}

	if slices.Contains(includes, "organizations") {
		org, err := c.modules.Organization.GetByID(r.Context(), member.OrganizationID)
		switch {
		case errors.Is(err, errx.ErrorOrganizationNotFound):
			log.Info("organization not found")
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
