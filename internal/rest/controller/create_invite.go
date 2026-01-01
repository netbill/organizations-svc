package controller

import (
	"errors"
	"net/http"
	"time"

	"github.com/umisto/agglomerations-svc/internal/domain/errx"
	"github.com/umisto/agglomerations-svc/internal/domain/modules/invite"
	"github.com/umisto/agglomerations-svc/internal/rest"
	"github.com/umisto/agglomerations-svc/internal/rest/request"
	"github.com/umisto/agglomerations-svc/internal/rest/responses"
	"github.com/umisto/ape"
	"github.com/umisto/ape/problems"
)

func (s Service) CreateInvite(w http.ResponseWriter, r *http.Request) {
	req, err := request.SentInvite(r)
	if err != nil {
		s.log.WithError(err).Errorf("invalid create invite request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	initiator, err := rest.AccountData(r)
	if err != nil {
		s.log.WithError(err).Errorf("failed to get initiator account data")
		ape.RenderErr(w, problems.Unauthorized("failed to get initiator account data"))
		return
	}

	inv, err := s.core.CreateInvite(r.Context(),
		initiator.ID,
		invite.CreateParams{
			AgglomerationID: req.Data.Attributes.AgglomerationId,
			AccountID:       req.Data.Attributes.AccountId,
			ExpiresAt:       time.Now().UTC().Add(24 * time.Hour),
		})
	if err != nil {
		s.log.WithError(err).Errorf("failed to create invite")
		switch {
		case errors.Is(err, errx.ErrorNotEnoughRights):
			ape.RenderErr(w, problems.Forbidden("not enough rights to create invite"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}
		return
	}

	ape.Render(w, http.StatusCreated, responses.Invite(inv))
}
