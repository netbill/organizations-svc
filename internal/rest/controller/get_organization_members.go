package controller

import (
	"fmt"
	"net/http"
	"slices"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/modules/member"
	"github.com/netbill/organizations-svc/internal/rest/responses"
	"github.com/netbill/organizations-svc/internal/rest/scope"
	"github.com/netbill/restkit/pagi"
	"github.com/netbill/restkit/problems"
	"github.com/netbill/restkit/render"
)

const operationGetOrganizationMembers = "get_organization_members"

func (c *Controller) GetMembers(w http.ResponseWriter, r *http.Request) {
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

	members, err := c.modules.Member.GetList(r.Context(), params, limit, offset)
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

		profiles, err := c.modules.Profile.GetByIDs(r.Context(), profileIDs)
		if err != nil {
			log.WithError(err).Error("failed to get profiles for members")
			render.ResponseError(w, problems.InternalError())
			return
		}

		opts = append(opts, responses.WithCollectionMembersProfiles(profiles))
	}

	if slices.Contains(includes, "organization") && params.OrganizationID != nil {
		organization, err := c.modules.Organization.GetByID(r.Context(), *params.OrganizationID)
		if err != nil {
			log.WithError(err).Error("failed to get organization for members")
			render.ResponseError(w, problems.InternalError())
			return
		}

		opts = append(opts, responses.WithCollectionMembersOrganization(organization))
	}

	render.Response(w, http.StatusOK, responses.Members(r, members, opts...))
}
