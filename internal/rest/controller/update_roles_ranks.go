package controller

import (
	"errors"
	"net/http"

	"github.com/google/uuid"

	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/rest/contexter"
	"github.com/netbill/organizations-svc/internal/rest/request"
	"github.com/netbill/restkit/problems"
)

func (c *Controller) UpdateRolesRanks(w http.ResponseWriter, r *http.Request) {
	req, err := request.UpdateRolesRanks(r)
	if err != nil {
		c.log.WithError(err).Errorf("invalid update roles ranks request")
		c.responser.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	initiator, err := contexter.AccountData(r.Context())
	if err != nil {
		c.log.WithError(err).Errorf("failed to get initiator account data")
		c.responser.RenderErr(w, problems.Unauthorized("failed to get initiator account data"))
		return
	}

	dict := make(map[uuid.UUID]uint)
	for _, item := range req.Data.Attributes.Roles {
		dict[item.Id] = item.Rank
	}

	err = c.core.UpdateRolesRanks(r.Context(), initiator.GetAccountID(), req.Data.Id, dict)
	if err != nil {
		c.log.WithError(err).Errorf("failed to update roles ranks")
		switch {
		case errors.Is(err, errx.ErrorNotEnoughRights):
			c.responser.RenderErr(w, problems.Forbidden("not enough rights to update roles ranks"))
		case errors.Is(err, errx.ErrorCannotUpdateHeadRoleRank):
			c.responser.RenderErr(w, problems.Forbidden("cannot update head role rank"))
		//case errors.Is(err, errx.ErrorInvalidInput):
		//	c.responser.RenderErr(w, problems.BadRequest(validation.Errors{
		//		"roles": fmt.Errorf(err.Error()),
		//	})...)
		default:
			c.responser.RenderErr(w, problems.InternalError())
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}
