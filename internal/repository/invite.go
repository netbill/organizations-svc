package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/umisto/cities-svc/internal/domain/entity"
	"github.com/umisto/cities-svc/internal/domain/modules/invite"
	"github.com/umisto/cities-svc/internal/repository/pgdb"
)

func (s Service) CreateInvite(
	ctx context.Context,
	params invite.CreateParams,
) (entity.Invite, error) {
	row, err := s.invitesQ().Insert(ctx, pgdb.InsertInviteParams{
		AgglomerationID: params.AgglomerationID,
		AccountID:       params.AccountID,
		ExpiresAt:       params.ExpiresAt,
	})
	if err != nil {
		return entity.Invite{}, err
	}

	return Invite(row), nil
}

func (s Service) GetInviteByID(
	ctx context.Context,
	id uuid.UUID,
) (entity.Invite, error) {
	row, err := s.invitesQ().FilterByID(id).Get(ctx)
	if err != nil {
		return entity.Invite{}, err
	}

	return Invite(row), nil
}

func (s Service) UpdateInviteStatus(
	ctx context.Context,
	id uuid.UUID,
	status string,
) (entity.Invite, error) {
	row, err := s.invitesQ().FilterByID(id).UpdateStatus(status).UpdateOne(ctx)
	if err != nil {
		return entity.Invite{}, err
	}

	return Invite(row), nil
}

func (s Service) DeleteInvite(
	ctx context.Context,
	id uuid.UUID,
) error {
	return s.invitesQ().FilterByID(id).Delete(ctx)
}

func (s Service) FilterInvites(
	ctx context.Context,
	filter invite.FilterInviteParams,
) ([]entity.Invite, error) {
	q := s.invitesQ()

	if filter.AgglomerationID != nil {
		q = q.FilterByAgglomerationID(*filter.AgglomerationID)
	}
	if filter.AccountID != nil {
		q = q.FilterByAccountID(*filter.AccountID)
	}
	if filter.Status != nil {
		q = q.FilterByStatus(*filter.Status)
	}

	rows, err := q.Select(ctx)
	if err != nil {
		return []entity.Invite{}, err
	}

	res := make([]entity.Invite, 0, len(rows))
	for _, row := range rows {
		res = append(res, Invite(row))
	}

	return res, nil
}

func Invite(row pgdb.Invite) entity.Invite {
	return entity.Invite{
		ID:              row.ID,
		AgglomerationID: row.AgglomerationID,
		AccountID:       row.AccountID,
		Status:          row.Status,
		ExpiresAt:       row.ExpiresAt,
		CreatedAt:       row.CreatedAt,
	}
}
