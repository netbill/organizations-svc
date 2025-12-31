package controller

import (
	"net/http"

	"github.com/umisto/ape"
	"github.com/umisto/ape/problems"
	"github.com/umisto/cities-svc/internal/domain/modules/agglomeration"
	"github.com/umisto/cities-svc/internal/rest/request"
	"github.com/umisto/cities-svc/internal/rest/responses"
)

func (s Service) CreateAgglomeration(w http.ResponseWriter, r *http.Request) {
	req, err := request.CreateAgglomeration(r)
	if err != nil {
		s.log.WithError(err).Errorf("invalid create agglomeration request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	res, err := s.domain.Agglomeration.CreateAgglomeration(
		r.Context(),
		agglomeration.CreateParams{
			Name: req.Data.Attributes.Name,
			Icon: req.Data.Attributes.Icon,
		},
	)
	if err != nil {
		s.log.WithError(err).Errorf("failed to create agglomeration")
		switch {
		default:
			ape.RenderErr(w, problems.InternalError())
		}
		return
	}

	ape.Render(w, http.StatusCreated, responses.Agglomeration(res))
}
