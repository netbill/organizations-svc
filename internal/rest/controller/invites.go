package controller

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/umisto/ape"
	"github.com/umisto/ape/problems"
	"github.com/umisto/cities-svc/internal/domain/models"
	"github.com/umisto/cities-svc/internal/domain/modules/invite"
	"github.com/umisto/cities-svc/internal/rest"
	"github.com/umisto/cities-svc/internal/rest/request"
	"github.com/umisto/cities-svc/internal/rest/responses"
	"github.com/umisto/logium"
)

type Invite interface {
	CreateInvite(ctx context.Context, params invite.CreateParams) (models.Invite, error)
	SentInviteByUser(
		ctx context.Context,
		accountID uuid.UUID,
		params invite.CreateParams,
	) (models.Invite, error)

	GetInvite(ctx context.Context, id uuid.UUID) (models.Invite, error)
	FilterInvites(
		ctx context.Context,
		filter invite.FilterInviteParams,
	) ([]models.Invite, error)

	DeclineInvite(
		ctx context.Context,
		accountID, inviteID uuid.UUID,
	) (models.Invite, error)
	AcceptInvite(
		ctx context.Context,
		accountID, inviteID uuid.UUID,
	) (models.Invite, error)

	DeleteInvite(
		ctx context.Context,
		accountID, inviteID uuid.UUID,
	) error
}

type InviteController struct {
	domain Invite
	log    logium.Logger
}

func (c InviteController) SentInvite(w http.ResponseWriter, r *http.Request) {
	req, err := request.SentInvite(r)
	if err != nil {
		c.log.WithError(err).Errorf("failed to parse sent invite request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	initiator, err := rest.AccountData(r)
	if err != nil {
		c.log.WithError(err).Errorf("failed to get account data from context")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	res, err := c.domain.SentInviteByUser(
		r.Context(),
		initiator.ID,
		invite.CreateParams{
			AgglomerationID: req.Data.Attributes.AgglomerationId,
			AccountID:       req.Data.Attributes.AccountId,
			ExpiresAt:       time.Now().Add(24 * time.Hour),
		},
	)
	if err != nil {
		c.log.WithError(err).Errorf("failed to create invite")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, http.StatusCreated, responses.Invite(res))
}
