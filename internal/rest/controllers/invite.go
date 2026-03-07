package controllers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/invite"
	"github.com/netbill/organizations-svc/internal/errx"
	"github.com/netbill/organizations-svc/internal/models"
	"github.com/netbill/organizations-svc/internal/rest/requests"
	"github.com/netbill/organizations-svc/internal/rest/responses"
	"github.com/netbill/organizations-svc/internal/rest/scope"
	"github.com/netbill/restkit"
	"github.com/netbill/restkit/pagi"
	"github.com/netbill/restkit/problems"
	"github.com/netbill/restkit/render"
)

type inviteCore interface {
	Create(ctx context.Context, actor models.AccountActor, params invite.CreateParams) (models.Invite, error)

	GetForAccount(ctx context.Context, accountID uuid.UUID, inviteID uuid.UUID) (models.Invite, error)
	GetList(
		ctx context.Context,
		actor models.AccountActor,
		params invite.FilterParams,
		limit, offset uint,
	) (pagi.Page[[]models.Invite], error)

	Accept(ctx context.Context, actor models.AccountActor, inviteID uuid.UUID) (models.Invite, error)
	Decline(ctx context.Context, actor models.AccountActor, inviteID uuid.UUID) (models.Invite, error)
	Cancelled(ctx context.Context, actor models.AccountActor, inviteID uuid.UUID) (models.Invite, error)
}

type inviteGetter interface {
	GetByID(ctx context.Context, organizationID uuid.UUID) (models.Organization, error)
	GetByIDs(ctx context.Context, organizationIDs []uuid.UUID) ([]models.Organization, error)
}

type InviteController struct {
	invites       inviteCore
	organizations inviteGetter
	profiles      profileGetter
}

type InviteControllerDeps struct {
	Invites       inviteCore
	Organizations inviteGetter
	Profiles      profileGetter
}

func NewInviteController(deps InviteControllerDeps) *InviteController {
	return &InviteController{
		invites:       deps.Invites,
		organizations: deps.Organizations,
		profiles:      deps.Profiles,
	}
}

const operationCreateInvite = "create_invite"

func (c *InviteController) Create(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationCreateInvite)

	req, err := requests.CreateInvite(r)
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
		invite.CreateParams{
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
	case errors.Is(err, errx.ErrorOrganizationNotExists),
		errors.Is(err, errx.ErrorOrganizationDeleted):
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
		render.Response(w, http.StatusCreated, responses.Invite(r, inv))
	}
}

const operationGetInvite = "get_invite"

func (c *InviteController) Get(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationGetInvite)

	inviteID, err := uuid.Parse(chi.URLParam(r, "invite_id"))
	if err != nil {
		log.WithError(err).Warn("invalid invite id")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"path/invite_id": fmt.Errorf("invalid invite id: %s", chi.URLParam(r, "invite_id")),
		})...)
		return
	}

	log = log.WithField("invite_id", inviteID)

	res, err := c.invites.GetForAccount(r.Context(), scope.AccountActor(r), inviteID)
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

	opts := make([]responses.InviteOption, 0)
	includes := restkit.ParseIncludes(r)

	if slices.Contains(includes, "profile") {
		profile, err := c.profiles.GetByID(r.Context(), res.AccountID)
		if err != nil {
			log.WithError(err).Error("failed to get profile for invite")
		}

		opts = append(opts, responses.WithInviteProfile(r, profile))
	}

	if slices.Contains(includes, "organizations") {
		org, err := c.organizations.GetByID(r.Context(), res.OrganizationID)
		if err != nil {
			log.WithError(err).Error("failed to get organization for invite")
		}

		opts = append(opts, responses.WithInviteOrganization(r, org))
	}

	render.Response(w, http.StatusOK, responses.Invite(r, res, opts...))
}

const operationGetInvites = "get_invites"

func (c *InviteController) GetList(w http.ResponseWriter, r *http.Request) {
	log := scope.Log(r).WithOperation(operationGetInvites)

	limit, offset := pagi.GetPagination(r)
	if limit > 100 {
		log.Warn("invalid pagination limit")
		render.ResponseError(w, problems.BadRequest(validation.Errors{
			"query/size": fmt.Errorf("pagination limit cannot be greater than 100"),
		})...)
		return
	}

	log = log.WithField("limit", limit).WithField("offset", offset)

	params := invite.FilterParams{}
	if accountIdStr := r.URL.Query().Get("account_id"); accountIdStr != "" {
		accountID, err := uuid.Parse(accountIdStr)
		if err != nil {
			log.WithError(err).Warn("invalid account id in query")
			render.ResponseError(w, problems.BadRequest(validation.Errors{
				"query/account_id": fmt.Errorf("invalid account id: %s", accountIdStr),
			})...)
			return
		}

		params.AccountID = &accountID
	}

	if organizationIdStr := r.URL.Query().Get("organization_id"); organizationIdStr != "" {
		organizationID, err := uuid.Parse(organizationIdStr)
		if err != nil {
			log.WithError(err).Warn("invalid organization id in query")
			render.ResponseError(w, problems.BadRequest(validation.Errors{
				"query/organization_id": fmt.Errorf("invalid organization id: %s", organizationIdStr),
			})...)
			return
		}

		params.OrganizationID = &organizationID
	}

	invites, err := c.invites.GetList(r.Context(), scope.AccountActor(r), params, limit, offset)
	switch {
	case errors.Is(err, errx.ErrorInitiatorNotMemberOfOrganization):
		log.WithError(err).Warn("initiator is not a member of organization, cannot get invites")
		render.ResponseError(w, problems.Forbidden("only members of the organization can get invites"))
		return
	case errors.Is(err, errx.ErrorCannotGetInvitesForOtherAccount):
		log.WithError(err).Warn("cannot get invites for other account")
		render.ResponseError(w, problems.Forbidden("cannot get invites for other account"))
		return
	case err != nil:
		log.WithError(err).Error("failed to get invites")
		render.ResponseError(w, problems.InternalError())
		return
	}

	opts := make([]responses.InvitesCollectionOption, 0, 2)
	includes := restkit.ParseIncludes(r)

	if slices.Contains(includes, "organizations") {
		organizationIDs := make([]uuid.UUID, 0, len(invites.Data))
		for _, i := range invites.Data {
			if !slices.Contains(organizationIDs, i.OrganizationID) {
				organizationIDs = append(organizationIDs, i.OrganizationID)
			}
		}

		organization, err := c.organizations.GetByIDs(r.Context(), organizationIDs)
		if err != nil {
			log.WithError(err).Error("failed to get organizations for invites")
		}

		opts = append(opts, responses.WithCollectionInvitesOrganizations(r, organization))
	}

	if slices.Contains(includes, "profiles") {
		accountIDs := make([]uuid.UUID, 0, len(invites.Data))
		for _, i := range invites.Data {
			if !slices.Contains(accountIDs, i.AccountID) {
				accountIDs = append(accountIDs, i.AccountID)
			}
		}

		profiles, err := c.profiles.GetByIDs(r.Context(), accountIDs)
		if err != nil {
			log.WithError(err).Error("failed to get profiles for invites")
		}

		opts = append(opts, responses.WithCollectionInvitesProfiles(r, profiles))
	}

	render.Response(w, http.StatusOK, responses.Invites(r, invites, opts...))
}
