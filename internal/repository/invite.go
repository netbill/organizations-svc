package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/umisto/agglomerations-svc/internal/domain/models"
	"github.com/umisto/agglomerations-svc/internal/domain/modules/invite"
	"github.com/umisto/agglomerations-svc/internal/repository/pgdb"
	"github.com/umisto/pagi"
)

func (s Service) CreateInvite(
	ctx context.Context,
	params invite.CreateParams,
) (models.Invite, error) {
	row, err := s.invitesQ().Insert(ctx, pgdb.InsertInviteParams{
		AgglomerationID: params.AgglomerationID,
		AccountID:       params.AccountID,
		ExpiresAt:       params.ExpiresAt,
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
	row, err := s.invitesQ().FilterByID(id).Get(ctx)
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
	row, err := s.invitesQ().FilterByID(id).UpdateStatus(status).UpdateOne(ctx)
	if err != nil {
		return models.Invite{}, err
	}

	return Invite(row), nil
}

func (s Service) DeleteInvite(
	ctx context.Context,
	id uuid.UUID,
) error {
	return s.invitesQ().FilterByID(id).Delete(ctx)
}

func (s Service) GetAgglomerationInvites(
	ctx context.Context,
	agglomerationID uuid.UUID,
	limit, offset uint,
) (pagi.Page[[]models.Invite], error) {
	rows, err := s.invitesQ().
		FilterByAgglomerationID(agglomerationID).
		Select(ctx)
	if err != nil {
		return pagi.Page[[]models.Invite]{}, err
	}

	total, err := s.invitesQ().
		FilterByAgglomerationID(agglomerationID).
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
		Total: uint(total),
	}, nil
}

func (s Service) GetAccountInvites(
	ctx context.Context,
	accountID uuid.UUID,
	limit, offset uint,
) (pagi.Page[[]models.Invite], error) {
	rows, err := s.invitesQ().
		FilterByAccountID(accountID).
		Select(ctx)
	if err != nil {
		return pagi.Page[[]models.Invite]{}, err
	}

	total, err := s.invitesQ().
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
		Total: uint(total),
	}, nil
}

func Invite(row pgdb.Invite) models.Invite {
	return models.Invite{
		ID:              row.ID,
		AgglomerationID: row.AgglomerationID,
		AccountID:       row.AccountID,
		Status:          row.Status,
		ExpiresAt:       row.ExpiresAt,
		CreatedAt:       row.CreatedAt,
	}
}
