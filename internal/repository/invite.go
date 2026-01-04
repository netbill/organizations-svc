package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/models"
	"github.com/netbill/organizations-svc/internal/core/modules/invite"
	"github.com/netbill/organizations-svc/internal/repository/pgdb"
	"github.com/netbill/pagi"
)

func (s Service) CreateInvite(
	ctx context.Context,
	params invite.CreateParams,
) (models.Invite, error) {
	row, err := s.invitesQ(ctx).Insert(ctx, pgdb.InsertInviteParams{
		OrganizationID: params.OrganizationID,
		AccountID:      params.AccountID,
		ExpiresAt:      params.ExpiresAt,
	})
	if err != nil {
		return models.Invite{}, err
	}

	return Invite(row), nil
}

func (s Service) GetInvite(
	ctx context.Context,
	id uuid.UUID,
) (models.Invite, error) {
	row, err := s.invitesQ(ctx).FilterByID(id).Get(ctx)
	if err != nil {
		return models.Invite{}, err
	}

	return Invite(row), nil
}

func (s Service) UpdateInviteStatus(
	ctx context.Context,
	id uuid.UUID,
	status string,
) (models.Invite, error) {
	row, err := s.invitesQ(ctx).FilterByID(id).UpdateStatus(status).UpdateOne(ctx)
	if err != nil {
		return models.Invite{}, err
	}

	return Invite(row), nil
}

func (s Service) DeleteInvite(
	ctx context.Context,
	id uuid.UUID,
) error {
	return s.invitesQ(ctx).FilterByID(id).Delete(ctx)
}

func (s Service) GetOrganizationInvites(
	ctx context.Context,
	organizationID uuid.UUID,
	limit, offset uint,
) (pagi.Page[[]models.Invite], error) {
	if limit == 0 {
		limit = 10
	}

	rows, err := s.invitesQ(ctx).
		FilterByOrganizationID(organizationID).
		Page(limit, offset).
		Select(ctx)
	if err != nil {
		return pagi.Page[[]models.Invite]{}, err
	}

	total, err := s.invitesQ(ctx).
		FilterByOrganizationID(organizationID).
		Count(ctx)
	if err != nil {
		return pagi.Page[[]models.Invite]{}, err
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

func (s Service) GetAccountInvites(
	ctx context.Context,
	accountID uuid.UUID,
	limit, offset uint,
) (pagi.Page[[]models.Invite], error) {
	if limit == 0 {
		limit = 10
	}

	rows, err := s.invitesQ(ctx).
		FilterByAccountID(accountID).
		Page(limit, offset).
		Select(ctx)
	if err != nil {
		return pagi.Page[[]models.Invite]{}, err
	}

	total, err := s.invitesQ(ctx).
		FilterByAccountID(accountID).
		Count(ctx)
	if err != nil {
		return pagi.Page[[]models.Invite]{}, err
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

func Invite(row pgdb.Invite) models.Invite {
	return models.Invite{
		ID:             row.ID,
		OrganizationID: row.OrganizationID,
		AccountID:      row.AccountID,
		Status:         row.Status,
		ExpiresAt:      row.ExpiresAt,
		CreatedAt:      row.CreatedAt,
	}
}
