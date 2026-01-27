package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/models"
	"github.com/netbill/organizations-svc/internal/core/modules/invite"
	"github.com/netbill/organizations-svc/internal/repository/pgdb"
	"github.com/netbill/pagi"
)

func (r Repository) CreateInvite(
	ctx context.Context,
	params invite.CreateParams,
) (models.Invite, error) {
	row, err := r.orgInvitesQ(ctx).Insert(ctx, pgdb.InsertInviteParams{
		OrganizationID: params.OrganizationID,
		AccountID:      params.AccountID,
		ExpiresAt:      params.ExpiresAt,
	})
	if err != nil {
		return models.Invite{}, fmt.Errorf(
			"failed to create invite for organization ID %s and account ID %s cause: %w",
			params.OrganizationID, params.AccountID, err,
		)
	}

	return Invite(row), nil
}

func (r Repository) GetInvite(
	ctx context.Context,
	id uuid.UUID,
) (models.Invite, error) {
	row, err := r.orgInvitesQ(ctx).FilterByID(id).Get(ctx)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return models.Invite{}, errx.ErrorInviteNotFound.Raise(
			fmt.Errorf("invite with ID %s not found, cause: %w", id, err),
		)
	case err != nil:
		return models.Invite{}, fmt.Errorf("failed to get invite with ID %s cause: %w", id, err)
	}

	return Invite(row), nil
}

func (r Repository) UpdateInviteStatus(
	ctx context.Context,
	id uuid.UUID,
	status string,
) (models.Invite, error) {
	row, err := r.orgInvitesQ(ctx).FilterByID(id).UpdateStatus(status).UpdateOne(ctx)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return models.Invite{}, errx.ErrorInviteNotFound.Raise(
			fmt.Errorf("invite with ID %s not found: %w", id, err),
		)
	case err != nil:
		return models.Invite{}, fmt.Errorf("failed to update invite with ID %s cause: %w", id, err)
	}

	return Invite(row), nil
}

func (r Repository) DeleteInvite(
	ctx context.Context,
	ID uuid.UUID,
) error {
	err := r.orgInvitesQ(ctx).FilterByID(ID).Delete(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete invite with ID %s cause: %w", ID, err)
	}

	return nil
}

func (r Repository) GetOrganizationInvites(
	ctx context.Context,
	organizationID uuid.UUID,
	limit, offset uint,
) (pagi.Page[[]models.Invite], error) {
	if limit == 0 {
		limit = 10
	}

	rows, err := r.orgInvitesQ(ctx).
		FilterByOrganizationID(organizationID).
		Page(limit, offset).
		Select(ctx)
	if err != nil {
		return pagi.Page[[]models.Invite]{}, fmt.Errorf(
			"failed to get invites for organization ID %s, cause: %w",
			organizationID, err,
		)
	}

	total, err := r.orgInvitesQ(ctx).
		FilterByOrganizationID(organizationID).
		Count(ctx)
	if err != nil {
		return pagi.Page[[]models.Invite]{}, fmt.Errorf(
			"failed to count invites for organization ID %s, cause: %w",
			organizationID, err,
		)
	}

	res := make([]models.Invite, 0, len(rows))
	for _, row := range rows {
		res = append(res, Invite(row))
	}

	return pagi.Page[[]models.Invite]{
		Data:  res,
		Page:  uint(offset/limit) + 1,
		Size:  uint(len(res)),
		Total: total,
	}, nil
}

func (r Repository) GetAccountInvites(
	ctx context.Context,
	accountID uuid.UUID,
	limit, offset uint,
) (pagi.Page[[]models.Invite], error) {
	if limit == 0 {
		limit = 10
	}

	rows, err := r.orgInvitesQ(ctx).
		FilterByAccountID(accountID).
		Page(limit, offset).
		Select(ctx)
	if err != nil {
		return pagi.Page[[]models.Invite]{}, fmt.Errorf(
			"failed to get invites for account ID %s, cause: %w",
			accountID, err,
		)
	}

	total, err := r.orgInvitesQ(ctx).
		FilterByAccountID(accountID).
		Count(ctx)
	if err != nil {
		return pagi.Page[[]models.Invite]{}, fmt.Errorf(
			"failed to count invites for account ID %s, cause: %w",
			accountID, err,
		)
	}

	res := make([]models.Invite, 0, len(rows))
	for _, row := range rows {
		res = append(res, Invite(row))
	}

	return pagi.Page[[]models.Invite]{
		Data:  res,
		Page:  uint(offset/limit) + 1,
		Size:  uint(len(res)),
		Total: total,
	}, nil
}

func Invite(row pgdb.OrganizationInvite) models.Invite {
	return models.Invite{
		ID:             row.ID,
		OrganizationID: row.OrganizationID,
		AccountID:      row.AccountID,
		Status:         row.Status,
		ExpiresAt:      row.ExpiresAt,
		CreatedAt:      row.CreatedAt,
	}
}
