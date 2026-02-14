package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/rest/responses"
	"github.com/netbill/organizations-svc/internal/rest/scope"
	"github.com/netbill/restkit/problems"
)

const operationGetMember = "get_member"

func (c *Controller) GetMember(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationGetMember)

	memberID, err := uuid.Parse(chi.URLParam(r, "member_id"))
	if err != nil {
		log.WithError(err).Info("invalid member id")
		c.responser.RenderErr(w, problems.BadRequest(fmt.Errorf("invalid member id"))...)
		return
	}

	log = log.WithField("member_id", memberID)

	res, err := c.modules.Member.GetByID(r.Context(), memberID)
	switch {
	case errors.Is(err, errx.ErrorMemberNotFound):
		log.Info("member not found")
		c.responser.RenderErr(w, problems.NotFound("member not found"))
	case err != nil:
		log.WithError(err).Error("failed to get member by id")
		c.responser.RenderErr(w, problems.InternalError())
	default:
		c.responser.Render(w, http.StatusOK, responses.Member(res))
	}
}
