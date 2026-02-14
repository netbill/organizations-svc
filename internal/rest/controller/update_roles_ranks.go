package controller

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/rest/request"
	"github.com/netbill/organizations-svc/internal/rest/scope"
	"github.com/netbill/restkit/problems"
)

const operationUpdateRolesRanks = "update_roles_ranks"

func (c *Controller) UpdateRolesRanks(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationUpdateRolesRanks)

	req, err := request.UpdateRolesRanks(r)
	if err != nil {
		log.WithError(err).Info("invalid update roles ranks request")
		c.responser.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	log = log.WithField("organization_id", req.Data.Id)

	dict := make(map[uuid.UUID]uint)
	for _, item := range req.Data.Attributes.Roles {
		dict[item.Id] = item.Rank
	}

	err = c.modules.Role.UpdateRanks(r.Context(), scope.AccountActor(r), req.Data.Id, dict)
	switch {
	case errors.Is(err, errx.ErrorNotEnoughRights):
		log.Info("not enough rights to update roles ranks")
		c.responser.RenderErr(w, problems.Forbidden("not enough rights to update roles ranks"))
	case err != nil:
		log.WithError(err).Error("failed to update roles ranks")
		c.responser.RenderErr(w, problems.InternalError())
	default:
		w.WriteHeader(http.StatusOK)
	}
}
