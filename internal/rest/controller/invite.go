package controller

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/models"
	"github.com/netbill/organizations-svc/internal/rest/requests"
	"github.com/netbill/organizations-svc/internal/rest/responses"
	"github.com/netbill/organizations-svc/internal/rest/scope"
	"github.com/netbill/restkit/pagi"
	"github.com/netbill/restkit/problems"
	"github.com/netbill/restkit/render"
)

type inviteModule interface {
	Create(ctx context.Context, actor models.AccountActor, params core.InviteCreateParams) (models.Invite, error)
	Cancelled(ctx context.Context, actor models.AccountActor, inviteID uuid.UUID) (models.Invite, error)
	Accept(ctx context.Context, actor models.AccountActor, inviteID uuid.UUID) (models.Invite, error)
	Decline(ctx context.Context, actor models.AccountActor, inviteID uuid.UUID) (models.Invite, error)
	GetForAccount(ctx context.Context, accountID uuid.UUID, inviteID uuid.UUID) (models.Invite, error)
	GetListForOrganization(ctx context.Context, actor models.AccountActor, organizationID uuid.UUID, limit, offset uint) (pagi.Page[[]models.Invite], error)
	GetListForAccount(ctx context.Context, actor models.AccountActor, limit, offset uint) (pagi.Page[[]models.Invite], error)
}

type inviteOrganizationModule interface {
	GetByID(ctx context.Context, organizationID uuid.UUID) (models.Organization, error)
	GetByIDs(ctx context.Context, organizationIDs []uuid.UUID) ([]models.Organization, error)
}

type inviteProfileModule interface {
	GetByID(ctx context.Context, accountID uuid.UUID) (models.Profile, error)
	GetByIDs(ctx context.Context, accountIDs []uuid.UUID) ([]models.Profile, error)
}

type InviteRouter struct {
	invites       inviteModule
	organizations inviteOrganizationModule
	profiles      inviteProfileModule
}

func NewInviteRouter(
	invites inviteModule,
	organizations inviteOrganizationModule,
	profiles inviteProfileModule,
) *InviteRouter {
	return &InviteRouter{
		invites:       invites,
		organizations: organizations,
		profiles:      profiles,
	}
}

const operationCreateInvite = "create_invite"

func (c *InviteRouter) Create(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationCreateInvite)

	req, err := requests.SentInvite(r)
	if err != nil {
		log.WithError(err).Warn("invalid create invite requests")
		render.ResponseError(w, problems.BadRequest(err)...)
		return
	}

	log = log.WithField("organization_id", req.Data.Attributes.OrganizationId).
		WithField("account_id", req.Data.Attributes.AccountId)

	inv, err := c.invites.Create(
		r.Context(),
		scope.AccountActor(r),
		core.InviteCreateParams{
			OrganizationID: req.Data.Attributes.OrganizationId,
			AccountID:      req.Data.Attributes.AccountId,
			ExpiresAt:      time.Now().UTC().Add(24 * time.Hour),
		},
	)
	switch {
	case errors.Is(err, errx.ErrorNotOrganizationHead):
		log.WithError(err).Warn("not enough rights to create invite")
		render.ResponseError(w, problems.Forbidden("not enough rights to create invite"))
	case errors.Is(err, errx.ErrorInitiatorNotMemberOfOrganization):
		log.WithError(err).Warn("membership in organization not found")
		render.ResponseError(w, problems.Forbidden("not a member of the organization"))
	case errors.Is(err, errx.ErrorProfileNotExists):
		log.WithError(err).Warn("profile for account not found")
		render.ResponseError(w, problems.NotFound("profile for account not found"))
	case errors.Is(err, errx.ErrorOrganizationNotExists):
		log.WithError(err).Warn("organization not found")
		render.ResponseError(w, problems.NotFound("organization not found"))
	case errors.Is(err, errx.ErrorAccountAlreadyMember):
		log.WithError(err).Warn("account is already a member of the organization")
		render.ResponseError(w, problems.Conflict("account is already a member of the organization"))
	case errors.Is(err, errx.ErrorActiveInviteAlreadyExists):
		log.WithError(err).Warn("active invite for this account and organization already exists")
		render.ResponseError(w, problems.Conflict("active invite for this account and organization already exists"))
	case err != nil:
		log.WithError(err).Error("failed to create invite")
		render.ResponseError(w, problems.InternalError())
	default:
		log.WithField("invite_id", inv.ID).Info("invite created successfully")
		render.Response(w, http.StatusCreated, responses.Invite(inv))
	}
}

const operationGetInvite = "get_invite"

func (c *InviteRouter) Get(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationGetInvite)

	inviteID, err := uuid.Parse(chi.URLParam(r, "invite_id"))
	if err != nil {
		log.WithError(err).Warn("invalid invite id")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"query": fmt.Errorf("invalid invite id: %s", chi.URLParam(r, "invite_id")),
		})...)
		return
	}

	log = log.WithField("invite_id", inviteID)

	invite, err := c.invites.GetForAccount(r.Context(), scope.AccountActor(r), inviteID)
	switch {
	case errors.Is(err, errx.ErrorInviteNotExists):
		log.Info("invite not found")
		render.ResponseError(w, problems.NotFound("invite not found"))
		return
	case err != nil:
		log.WithError(err).Error("failed to get invite")
		render.ResponseError(w, problems.InternalError())
		return
	}

	includes := r.URL.Query()["include"]
	opts := make([]responses.InviteOption, 0)

	if slices.Contains(includes, "profile") {
		profile, err := c.profiles.GetByID(r.Context(), invite.AccountID)
		switch {
		case errors.Is(err, errx.ErrorProfileNotExists):
			log.WithError(err).WithField("account_id", invite.AccountID).Warn("profile not found")
			render.ResponseError(w, problems.NotFound("profile not found"))
			return
		case err != nil:
			log.WithError(err).WithField("account_id", invite.AccountID).Error("failed to get profile")
			render.ResponseError(w, problems.InternalError())
			return
		default:
			opts = append(opts, responses.WithInviteProfile(profile))
		}
	}

	if slices.Contains(includes, "organizations") {
		org, err := c.organizations.GetByID(r.Context(), invite.OrganizationID)
		switch {
		case errors.Is(err, errx.ErrorOrganizationNotExists):
			log.WithError(err).Warn("organization not found")
			render.ResponseError(w, problems.NotFound("organization not found"))
			return
		case err != nil:
			log.WithError(err).Error("failed to get organization")
			render.ResponseError(w, problems.InternalError())
			return
		default:
			opts = append(opts, responses.WithInviteOrganization(org))
		}
	}

	render.Response(w, http.StatusOK, responses.Invite(invite, opts...))
}

const operationGetMyInvites = "get_my_invites"

func (c *InviteRouter) GetMyList(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationGetMyInvites)

	limit, offset := pagi.GetPagination(r)
	if limit > 100 {
		log.Warn("invalid pagination limit")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"query/limit": fmt.Errorf("pagination limit cannot be greater than 100"),
		})...)
		return
	}

	accountID := scope.AccountActor(r)
	log = log.WithField("account_id", accountID).WithField("limit", limit).WithField("offset", offset)

	invites, err := c.invites.GetListForAccount(r.Context(), accountID, limit, offset)
	switch {
	case errors.Is(err, errx.ErrorProfileNotExists):
		log.WithError(err).Warn("account not found")
		render.ResponseError(w, problems.NotFound("account not found"))
		return
	case err != nil:
		log.WithError(err).Error("failed to get invites")
		render.ResponseError(w, problems.InternalError())
		return
	}

	opts := make([]responses.InvitesCollectionOption, 0, 1)
	includesRaw := r.URL.Query()["include"]
	includes := make([]string, 0, 1)

	for _, v := range includesRaw {
		for _, part := range strings.Split(v, ",") {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}
			if !slices.Contains(includes, part) {
				includes = append(includes, part)
			}
		}
	}

	organizationIDs := make([]uuid.UUID, 0, invites.Size)
	for _, invite := range invites.Data {
		if !slices.Contains(organizationIDs, invite.OrganizationID) {
			organizationIDs = append(organizationIDs, invite.OrganizationID)
		}
	}

	if slices.Contains(includes, "organizations") {
		organization, err := c.organizations.GetByIDs(r.Context(), organizationIDs)
		if err != nil {
			log.WithError(err).Error("failed to get organizations for invites")
			render.ResponseError(w, problems.InternalError())
			return
		}
		opts = append(opts, responses.WithCollectionInvitesOrganizations(organization))
	}

	render.Response(w, http.StatusOK, responses.Invites(r, invites, opts...))
}

const operationGetOrganizationInvites = "get_organization_invites"

func (c *InviteRouter) GetListForOrg(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationGetOrganizationInvites)

	organizationID, err := uuid.Parse(chi.URLParam(r, "organization_id"))
	if err != nil {
		log.WithError(err).Warn("invalid organization id")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"query": fmt.Errorf("invalid organization id: %s", chi.URLParam(r, "organization_id")),
		})...)
		return
	}

	limit, offset := pagi.GetPagination(r)
	if limit > 100 {
		log.Warn("invalid pagination limit")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"query": fmt.Errorf("pagination limit cannot be greater than 100"),
		})...)
		return
	}

	log = log.WithField("organization_id", organizationID).
		WithField("limit", limit).
		WithField("offset", offset)

	invites, err := c.invites.GetListForOrganization(
		r.Context(),
		scope.AccountActor(r),
		organizationID,
		limit, offset,
	)
	switch {
	case errors.Is(err, errx.ErrorOrganizationNotExists):
		log.WithError(err).Warn("organization not found")
		render.ResponseError(w, problems.NotFound("organization not found"))
		return
	case errors.Is(err, errx.ErrorInitiatorNotMemberOfOrganization):
		log.WithError(err).Warn("not enough rights to access organization invites")
		render.ResponseError(w, problems.Forbidden("not enough rights to access organization invites"))
		return
	case err != nil:
		log.WithError(err).Error("failed to get organization invites")
		render.ResponseError(w, problems.InternalError())
		return
	}

	opts := make([]responses.InvitesCollectionOption, 0, 2)
	includesRaw := r.URL.Query()["include"]
	includes := make([]string, 0, 2)

	for _, v := range includesRaw {
		for _, part := range strings.Split(v, ",") {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}
			if !slices.Contains(includes, part) {
				includes = append(includes, part)
			}
		}
	}

	if slices.Contains(includes, "profile") {
		profileIDs := make([]uuid.UUID, 0, invites.Size)
		for _, invite := range invites.Data {
			profileIDs = append(profileIDs, invite.AccountID)
		}

		profiles, err := c.profiles.GetByIDs(r.Context(), profileIDs)
		if err != nil {
			log.WithError(err).Error("failed to get profiles for invites")
			render.ResponseError(w, problems.InternalError())
			return
		}
		opts = append(opts, responses.WithCollectionInvitesProfiles(profiles))
	}

	if slices.Contains(includes, "organization") {
		organization, err := c.organizations.GetByID(r.Context(), organizationID)
		if err != nil {
			log.WithError(err).Error("failed to get organizations for invites")
			render.ResponseError(w, problems.InternalError())
			return
		}
		opts = append(opts, responses.WithCollectionInvitesOrganization(organization))
	}

	render.Response(w, http.StatusOK, responses.Invites(r, invites, opts...))
}

const operationAcceptInvite = "accept_invite"

func (c *InviteRouter) Accept(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationAcceptInvite)

	inviteID, err := uuid.Parse(chi.URLParam(r, "invite_id"))
	if err != nil {
		log.WithError(err).Warn("invalid invite id")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"query": fmt.Errorf("invalid invite id: %s", chi.URLParam(r, "invite_id")),
		})...)
		return
	}

	log = log.WithField("invite_id", inviteID)

	res, err := c.invites.Accept(r.Context(), scope.AccountActor(r), inviteID)
	switch {
	case errors.Is(err, errx.ErrorInviteNotExists):
		log.WithError(err).Warn("invite not exists")
		render.ResponseError(w, problems.NotFound("invite not exists"))
	case errors.Is(err, errx.ErrorOrganizationNotExists):
		log.WithError(err).Warn("organization not found for invite")
		render.ResponseError(w, problems.NotFound("organization not found for invite"))
	case errors.Is(err, errx.ErrorOrganizationIsNotActive):
		log.WithError(err).Warn("organization is not active")
		render.ResponseError(w, problems.Forbidden("organization is not active"))
	case errors.Is(err, errx.ErrorInviteNotForInitiator):
		log.WithError(err).Warn("account has no rights to accept this invite")
		render.ResponseError(w, problems.Forbidden("account has no rights to accept this invite"))
	case errors.Is(err, errx.ErrorInviteAlreadyAnswered):
		log.WithError(err).Warn("invite already accepted")
		render.ResponseError(w, problems.Conflict("invite already accepted"))
	case errors.Is(err, errx.ErrorInviteExpired):
		log.WithError(err).Warn("invite has expired")
		render.ResponseError(w, problems.Forbidden("invite has expired"))
	case err != nil:
		log.WithError(err).Error("failed to accept invite")
		render.ResponseError(w, problems.InternalError())
	default:
		render.Response(w, http.StatusOK, responses.Invite(res))
	}
}

const operationDeclineInvite = "decline_invite"

func (c *InviteRouter) Decline(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationDeclineInvite)

	inviteID, err := uuid.Parse(chi.URLParam(r, "invite_id"))
	if err != nil {
		log.WithError(err).Warn("invalid invite id")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"query": fmt.Errorf("invalid invite id: %s", chi.URLParam(r, "invite_id")),
		})...)
		return
	}

	log = log.WithField("invite_id", inviteID)

	res, err := c.invites.Decline(r.Context(), scope.AccountActor(r), inviteID)
	switch {
	case errors.Is(err, errx.ErrorInviteNotExists):
		log.WithError(err).Warn("invite not found")
		render.ResponseError(w, problems.NotFound("invite not found"))
	case errors.Is(err, errx.ErrorInviteNotForInitiator):
		log.WithError(err).Warn("invite not for this account")
		render.ResponseError(w, problems.Forbidden("invite not for this account"))
	case errors.Is(err, errx.ErrorInviteAlreadyAnswered):
		log.WithError(err).Warn("invite already answered")
		render.ResponseError(w, problems.Conflict("invite already answered"))
	case errors.Is(err, errx.ErrorInviteExpired):
		log.WithError(err).Warn("invite has expired")
		render.ResponseError(w, problems.Forbidden("invite has expired"))
	case errors.Is(err, errx.ErrorOrganizationNotExists):
		log.WithError(err).Warn("organization not found for invite")
		render.ResponseError(w, problems.NotFound("organization not found for invite"))
	case errors.Is(err, errx.ErrorOrganizationIsNotActive):
		log.WithError(err).Warn("organization is not active")
		render.ResponseError(w, problems.Forbidden("organization is not active"))
	case err != nil:
		log.WithError(err).Error("failed to decline invite")
		render.ResponseError(w, problems.InternalError())
	default:
		render.Response(w, http.StatusOK, responses.Invite(res))
	}
}

const operationCancelledInvite = "delete_invite"

func (c *InviteRouter) Cancelled(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationCancelledInvite)

	inviteID, err := uuid.Parse(chi.URLParam(r, "invite_id"))
	if err != nil {
		log.WithError(err).Warn("invalid invite id")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"query": fmt.Errorf("invalid invite id: %s", chi.URLParam(r, "invite_id")),
		})...)
		return
	}

	log = log.WithField("invite_id", inviteID)

	invite, err := c.invites.Cancelled(r.Context(), scope.AccountActor(r), inviteID)
	switch {
	case errors.Is(err, errx.ErrorInviteNotExists):
		log.WithError(err).Warn("invite not found")
		render.ResponseError(w, problems.NotFound("invite not found"))
	case errors.Is(err, errx.ErrorInitiatorNotMemberOfOrganization):
		log.WithError(err).Warn("only organization members can cancel invite")
		render.ResponseError(w, problems.NotFound("invite not found"))
	case errors.Is(err, errx.ErrorNotOrganizationHead):
		log.WithError(err).Warn("only organization head can cancel invite")
		render.ResponseError(w, problems.Forbidden("only organization head can cancel invite"))
	case errors.Is(err, errx.ErrorInviteAlreadyAnswered):
		log.WithError(err).Warn("invite already answered")
		render.ResponseError(w, problems.Forbidden("invite already answered"))
	case errors.Is(err, errx.ErrorOrganizationNotExists):
		log.WithError(err).Warn("organization not found for invite")
		render.ResponseError(w, problems.NotFound("organization not found for invite"))
	case errors.Is(err, errx.ErrorOrganizationIsSuspended):
		log.WithError(err).Warn("organization is suspended")
		render.ResponseError(w, problems.Forbidden("organization is suspended"))
	case err != nil:
		log.WithError(err).Error("failed to delete invite")
		render.ResponseError(w, problems.InternalError())
	default:
		log.Info("invite deleted")
		render.Response(w, http.StatusOK, responses.Invite(invite))
	}
}
