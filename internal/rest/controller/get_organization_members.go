package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/netbill/ape"
	"github.com/netbill/ape/problems"
	"github.com/netbill/organizations-svc/internal/core/modules/member"
	"github.com/netbill/organizations-svc/internal/rest/responses"
	"github.com/netbill/pagi"
)

func (c Controller) GetOrganizationMembers(w http.ResponseWriter, r *http.Request) {
	organizationID, err := uuid.Parse(chi.URLParam(r, "organization_id"))
	if err != nil {
		c.log.Errorf("failed to parse organization id, cause %s", err)
		ape.RenderErr(w, problems.BadRequest(fmt.Errorf("invalid organization id"))...)
		return
	}

	limit, offset := pagi.GetPagination(r)
	if limit > 100 {
		c.log.WithError(fmt.Errorf("invalid pagination limit %d", limit)).Errorf("invalid pagination limit")
		ape.RenderErr(w, problems.BadRequest(fmt.Errorf("pagination limit must be between 1 and 100"))...)
		return
	}

	params := member.FilterParams{
		OrganizationID: &organizationID,
	}
	if _, ok := r.URL.Query()["role_id"]; ok {
		roleID, err := uuid.Parse(r.URL.Query().Get("role_id"))
		if err != nil {
			c.log.WithError(err).Errorf("invalid role id")
			ape.RenderErr(w, problems.BadRequest(fmt.Errorf("invalid role id"))...)
			return
		}
		params.RoleID = &roleID
	}

	if _, ok := r.URL.Query()["text"]; ok {
		text := r.URL.Query().Get("text")
		params.BestMatch = &text
	}

	if _, ok := r.URL.Query()["role_rank_up"]; ok {
		roleRankUp64, err := strconv.ParseUint(r.URL.Query().Get("role_rank_up"), 10, 64)
		if err != nil {
			c.log.WithError(err).Errorf("invalid role rank up")
			ape.RenderErr(w, problems.BadRequest(fmt.Errorf("invalid role rank up"))...)
			return
		}
		roleRankUp := uint(roleRankUp64)
		params.RoleRankUp = &roleRankUp
	}
	if _, ok := r.URL.Query()["role_rank_down"]; ok {
		roleRankDown64, err := strconv.ParseUint(r.URL.Query().Get("role_rank_down"), 10, 64)
		if err != nil {
			c.log.WithError(err).Errorf("invalid role rank down")
			ape.RenderErr(w, problems.BadRequest(fmt.Errorf("invalid role rank down"))...)
			return
		}
		roleRankDown := uint(roleRankDown64)
		params.RoleRankDown = &roleRankDown
	}

	members, err := c.core.GetMembers(r.Context(), params, limit, offset)
	if err != nil {
		c.log.WithError(err).Errorf("failed to get organization members")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, http.StatusOK, responses.Members(r, members))
}
