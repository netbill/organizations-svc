package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/netbill/ape"
	"github.com/netbill/ape/problems"
	"github.com/netbill/organizations-svc/internal/domain/errx"
)

func (c Controller) DeleteMember(w http.ResponseWriter, r *http.Request) {
	initiatorID, err := uuid.Parse(r.URL.Query().Get("initiator_id"))
	if err != nil {
		c.log.Errorf("failed to parse initiator id, cause %s", err)
		ape.RenderErr(w, problems.BadRequest(fmt.Errorf("invalid initiator id"))...)
		return
	}

	memberId, err := uuid.Parse(r.URL.Query().Get("member_id"))
	if err != nil {
		c.log.Errorf("failed to parse member id, cause %s", err)
		ape.RenderErr(w, problems.BadRequest(fmt.Errorf("invalid member id"))...)
		return
	}

	err = c.core.DeleteMember(r.Context(), initiatorID, memberId)
	if err != nil {
		c.log.WithError(err).Errorf("failed to delete member")
		switch {
		case errors.Is(err, errx.ErrorMemberNotFound):
			ape.RenderErr(w, problems.NotFound("member not found"))
		case errors.Is(err, errx.ErrorNotEnoughRights):
			ape.RenderErr(w, problems.Forbidden("not enough rights to delete member"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
