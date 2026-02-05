package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/rest/responses"
	"github.com/netbill/restkit/problems"
)

func (c *Controller) GetMember(w http.ResponseWriter, r *http.Request) {
	memberId, err := uuid.Parse(chi.URLParam(r, "member_id"))
	if err != nil {
		c.log.Errorf("failed to parse member id, cause %s", err)
		c.responser.RenderErr(w, problems.BadRequest(fmt.Errorf("invalid member id"))...)
		return
	}

	res, err := c.core.member.GetByID(r.Context(), memberId)
	if err != nil {
		c.log.WithError(err).Errorf("failed to get member by id")
		switch {
		case errors.Is(err, errx.ErrorMemberNotFound):
			c.responser.RenderErr(w, problems.NotFound("member not found"))
		default:
			c.responser.RenderErr(w, problems.InternalError())
		}
		return
	}

	c.responser.Render(w, http.StatusOK, responses.Member(res))
}
