package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/models"
	"github.com/netbill/organizations-svc/internal/rest/contexter"
	"github.com/netbill/restkit/problems"
)

func (c *Controller) DeleteMember(w http.ResponseWriter, r *http.Request) {
	initiator, err := contexter.AccountData(r.Context())
	if err != nil {
		c.log.WithError(err).Errorf("failed to get initiator account data")
		c.responser.RenderErr(w, problems.Unauthorized("failed to get initiator account data"))
		return
	}

	memberId, err := uuid.Parse(chi.URLParam(r, "member_id"))
	if err != nil {
		c.log.Errorf("failed to parse member id, cause %s", err)
		c.responser.RenderErr(w, problems.BadRequest(fmt.Errorf("invalid member id"))...)
		return
	}

	err = c.core.DeleteMember(
		r.Context(),
		models.InitiatorData{
			AccountID: initiator.GetAccountID(),
		},
		memberId,
	)
	if err != nil {
		c.log.WithError(err).Errorf("failed to delete member")
		switch {
		case errors.Is(err, errx.ErrorMemberNotFound):
			c.responser.RenderErr(w, problems.NotFound("member not found"))
		case errors.Is(err, errx.ErrorNotEnoughRights):
			c.responser.RenderErr(w, problems.Forbidden("not enough rights to delete member"))
		case errors.Is(err, errx.ErrorCannotDeleteOrganizationHeadMember):
			c.responser.RenderErr(w, problems.Forbidden("cannot delete organization head member"))
		default:
			c.responser.RenderErr(w, problems.InternalError())
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
