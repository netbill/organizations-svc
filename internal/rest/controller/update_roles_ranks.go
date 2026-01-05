package controller

import (
	"errors"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/netbill/ape"
	"github.com/netbill/ape/problems"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/rest"
	"github.com/netbill/organizations-svc/internal/rest/request"
)

func (c Controller) UpdateRolesRanks(w http.ResponseWriter, r *http.Request) {
	req, err := request.UpdateRolesRanks(r)
	if err != nil {
		c.log.WithError(err).Errorf("invalid update roles ranks request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	initiator, err := rest.AccountData(r)
	if err != nil {
		c.log.WithError(err).Errorf("failed to get initiator account data")
		ape.RenderErr(w, problems.Unauthorized("failed to get initiator account data"))
		return
	}

	dict := make(map[uuid.UUID]uint)
	for _, item := range req.Data.Attributes.Roles {
		dict[item.Id] = item.Rank
	}

	err = c.core.UpdateRolesRanks(r.Context(), initiator.ID, req.Data.Id, dict)
	if err != nil {
		c.log.WithError(err).Errorf("failed to update roles ranks")
		switch {
		case errors.Is(err, errx.ErrorNotEnoughRights):
			ape.RenderErr(w, problems.Forbidden("not enough rights to update roles ranks"))
		case errors.Is(err, errx.ErrorCannotUpdateHeadRoleRank):
			ape.RenderErr(w, problems.Forbidden("cannot update head role rank"))
		case errors.Is(err, errx.ErrorInvalidInput):
			ape.RenderErr(w, problems.BadRequest(validation.Errors{
				"roles": fmt.Errorf(err.Error()),
			})...)
		default:
			ape.RenderErr(w, problems.InternalError())
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}
