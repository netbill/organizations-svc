package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/modules/member"
	"github.com/netbill/organizations-svc/internal/rest/responses"
	"github.com/netbill/organizations-svc/internal/rest/scope"
	"github.com/netbill/restkit/pagi"
	"github.com/netbill/restkit/problems"
)

const operationGetOrganizationMembers = "get_organization_members"

func (c *Controller) GetOrganizationMembers(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationGetOrganizationMembers)

	organizationID, err := uuid.Parse(chi.URLParam(r, "organization_id"))
	if err != nil {
		log.WithError(err).Info("invalid organization id")
		c.responser.RenderErr(w, problems.BadRequest(fmt.Errorf("invalid organization id"))...)
		return
	}

	limit, offset := pagi.GetPagination(r)
	if limit > 100 {
		log.Info("invalid pagination limit")
		c.responser.RenderErr(w, problems.BadRequest(fmt.Errorf("pagination limit must be between 1 and 100"))...)
		return
	}

	log = log.WithField("organization_id", organizationID).WithField("limit", limit).WithField("offset", offset)

	params := member.FilterParams{
		OrganizationID: &organizationID,
	}

	if v := r.URL.Query().Get("role_id"); v != "" {
		roleID, err := uuid.Parse(v)
		if err != nil {
			log.WithError(err).Info("invalid role id")
			c.responser.RenderErr(w, problems.BadRequest(fmt.Errorf("invalid role id"))...)
			return
		}
		params.RoleID = &roleID
		log = log.WithField("role_id", roleID)
	}

	if v := r.URL.Query().Get("text"); v != "" {
		text := v
		params.BestMatch = &text
	}

	if v := r.URL.Query().Get("role_rank_up"); v != "" {
		roleRankUp64, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			log.WithError(err).Info("invalid role rank up")
			c.responser.RenderErr(w, problems.BadRequest(fmt.Errorf("invalid role rank up"))...)
			return
		}
		roleRankUp := uint(roleRankUp64)
		params.RoleRankUp = &roleRankUp
	}

	if v := r.URL.Query().Get("role_rank_down"); v != "" {
		roleRankDown64, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			log.WithError(err).Info("invalid role rank down")
			c.responser.RenderErr(w, problems.BadRequest(fmt.Errorf("invalid role rank down"))...)
			return
		}
		roleRankDown := uint(roleRankDown64)
		params.RoleRankDown = &roleRankDown
	}

	members, err := c.core.member.GetList(r.Context(), params, limit, offset)
	switch {
	case err != nil:
		log.WithError(err).Error("failed to get organization members")
		c.responser.RenderErr(w, problems.InternalError())
	default:
		c.responser.Render(w, http.StatusOK, responses.Members(r, members))
	}
}
