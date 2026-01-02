package controller

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/netbill/ape"
	"github.com/netbill/ape/problems"
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
		ape.RenderErr(w, problems.InternalError())
		return
	}

	w.WriteHeader(http.StatusOK)
}
