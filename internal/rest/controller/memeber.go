package controller

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/umisto/ape"
	"github.com/umisto/ape/problems"
	"github.com/umisto/cities-svc/internal/domain/errx"
	"github.com/umisto/cities-svc/internal/domain/models"
	"github.com/umisto/cities-svc/internal/domain/modules/member"
	"github.com/umisto/cities-svc/internal/rest/request"
	"github.com/umisto/cities-svc/internal/rest/responses"
	"github.com/umisto/logium"
	"github.com/umisto/pagi"
)

type Member interface {
	GetMemberByID(ctx context.Context, ID uuid.UUID) (models.Member, error)
	GetMemberByAccountAndAgglomeration(ctx context.Context, accountID, agglomerationID uuid.UUID) (models.Member, error)
	GetInitiatorMember(ctx context.Context, accountID, agglomerationID uuid.UUID) (models.Member, error)

	FilterMembers(
		ctx context.Context,
		filter member.FilterParams,
		offset uint,
		limit uint,
	) (pagi.Page[[]models.Member], error)

	UpdateMemberByUser(
		ctx context.Context,
		accountID, memberID uuid.UUID,
		params member.UpdateParams,
	) (models.Member, error)

	DeleteMemberByUser(
		ctx context.Context,
		accountID, memberID uuid.UUID,
	) error
}

type MemberController struct {
	domain Member
	log    logium.Logger
}

func (c MemberController) GetMemberByID(w http.ResponseWriter, r *http.Request) {
	memberId, err := uuid.Parse(r.URL.Query().Get("member_id"))
	if err != nil {
		c.log.Errorf("failed to parse member id, cause %s", err)
		ape.RenderErr(w, problems.BadRequest(fmt.Errorf("invalid member id"))...)
		return
	}

	res, err := c.domain.GetMemberByID(r.Context(), memberId)
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

func (c MemberController) UpdateMember(w http.ResponseWriter, r *http.Request) {
	req, err := request.UpdateMember(r)
	if err != nil {
		c.log.Errorf("invalid update member request, cause %s", err)
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	initiatorID, err := uuid.Parse(r.URL.Query().Get("initiator_id"))
	if err != nil {
		c.log.Errorf("failed to parse initiator id, cause %s", err)
		ape.RenderErr(w, problems.BadRequest(fmt.Errorf("invalid initiator id"))...)
		return
	}

	memberId, err := uuid.Parse(r.URL.Query().Get("member_id"))
	if err != nil {
		c.log.Errorf("failed to parse member id, cause %s", err)
		ape.RenderErr(w, problems.BadRequest(fmt.Errorf("invalid member id"))...)
		return
	}

	res, err := c.domain.UpdateMemberByUser(r.Context(), initiatorID, memberId, member.UpdateParams{
		Position: req.Data.Attributes.Position,
		Label:    req.Data.Attributes.Label,
	})
	if err != nil {
		c.log.WithError(err).Errorf("failed to update member")
		switch {
		case errors.Is(err, errx.ErrorMemberNotFound):
			ape.RenderErr(w, problems.NotFound("member not found"))
		case errors.Is(err, errx.ErrorNotEnoughRights):
			ape.RenderErr(w, problems.Forbidden("not enough rights to update member"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}
		return
	}

	ape.Render(w, http.StatusOK, responses.Member(res))
}

func (c MemberController) DeleteMember(w http.ResponseWriter, r *http.Request) {
	initiatorID, err := uuid.Parse(r.URL.Query().Get("initiator_id"))
	if err != nil {
		c.log.Errorf("failed to parse initiator id, cause %s", err)
		ape.RenderErr(w, problems.BadRequest(fmt.Errorf("invalid initiator id"))...)
		return
	}

	memberId, err := uuid.Parse(r.URL.Query().Get("member_id"))
	if err != nil {
		c.log.Errorf("failed to parse member id, cause %s", err)
		ape.RenderErr(w, problems.BadRequest(fmt.Errorf("invalid member id"))...)
		return
	}

	err = c.domain.DeleteMemberByUser(r.Context(), initiatorID, memberId)
	if err != nil {
		c.log.WithError(err).Errorf("failed to delete member")
		switch {
		case errors.Is(err, errx.ErrorMemberNotFound):
			ape.RenderErr(w, problems.NotFound("member not found"))
		case errors.Is(err, errx.ErrorNotEnoughRights):
			ape.RenderErr(w, problems.Forbidden("not enough rights to delete member"))
		default:
			ape.RenderErr(w, problems.InternalError())
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
