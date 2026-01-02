package controller

import (
	"net/http"

	"github.com/umisto/agglomerations-svc/internal/domain/modules/agglomeration"
	"github.com/umisto/agglomerations-svc/internal/rest/request"
	"github.com/umisto/agglomerations-svc/internal/rest/responses"
	"github.com/umisto/ape"
	"github.com/umisto/ape/problems"
)

func (c Controller) CreateAgglomeration(w http.ResponseWriter, r *http.Request) {
	req, err := request.CreateAgglomeration(r)
	if err != nil {
		c.log.WithError(err).Errorf("invalid create agglomeration request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	res, err := c.core.CreateAgglomeration(
		r.Context(),
		req.Data.Attributes.Head,
		agglomeration.CreateParams{
			Name: req.Data.Attributes.Name,
			Icon: req.Data.Attributes.Icon,
		},
	)
	if err != nil {
		c.log.WithError(err).Errorf("failed to create agglomeration")
		switch {
		default:
			ape.RenderErr(w, problems.InternalError())
		}
		return
	}

	ape.Render(w, http.StatusCreated, responses.Agglomeration(res))
}
