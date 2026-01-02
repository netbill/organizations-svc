package controller

import (
	"fmt"
	"net/http"

	"github.com/umisto/agglomerations-svc/internal/rest"
	"github.com/umisto/agglomerations-svc/internal/rest/responses"
	"github.com/umisto/ape"
	"github.com/umisto/ape/problems"
	"github.com/umisto/pagi"
)

func (c Controller) GetMyAgglomerations(w http.ResponseWriter, r *http.Request) {
	initiator, err := rest.AccountData(r)
	if err != nil {
		c.log.WithError(err).Errorf("failed to get initiator account data")
		ape.RenderErr(w, problems.Unauthorized("failed to get initiator account data"))
		return
	}

	limit, offset := pagi.GetPagination(r)
	if limit == 0 || limit > 100 {
		c.log.WithError(fmt.Errorf("invalid pagination limit %d", limit)).Errorf("invalid pagination limit")
		ape.RenderErr(w, problems.BadRequest(fmt.Errorf("pagination limit must be between 1 and 100"))...)
		return
	}

	res, err := c.core.GetAgglomerationForUser(r.Context(), initiator.ID, limit, offset)
	if err != nil {
		c.log.WithError(err).Errorf("failed to get agglomerations")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, http.StatusOK, responses.Agglomerations(r, res))
}
