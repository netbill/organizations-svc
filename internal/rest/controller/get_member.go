package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/netbill/ape"
	"github.com/netbill/ape/problems"
	"github.com/netbill/organizations-svc/internal/core/errx"
)

func (c Controller) GetMember(w http.ResponseWriter, r *http.Request) {
	memberId, err := uuid.Parse(r.URL.Query().Get("member_id"))
	if err != nil {
		c.log.Errorf("failed to parse member id, cause %s", err)
		ape.RenderErr(w, problems.BadRequest(fmt.Errorf("invalid member id"))...)
		return
	}

	res, err := c.core.GetMemberByID(r.Context(), memberId)
	if err != nil {
		c.log.WithError(err).Errorf("failed to get member by id")
		switch {
		case errors.Is(err, errx.ErrorMemberNotFound):
			ape.RenderErr(w, problems.NotFound("member not found"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}
		return
	}

	ape.Render(w, http.StatusOK, res)
}
