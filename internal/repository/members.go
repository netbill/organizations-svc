package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/umisto/agglomerations-svc/internal/domain/models"
	"github.com/umisto/agglomerations-svc/internal/domain/modules/member"
	"github.com/umisto/agglomerations-svc/internal/repository/pgdb"
	"github.com/umisto/pagi"
)

func (s Service) CreateMember(ctx context.Context, accountID, agglomerationID uuid.UUID) (models.Member, error) {
	row, err := s.membersQ().Insert(ctx, pgdb.InsertMemberParams{
		AccountID:       accountID,
		AgglomerationID: agglomerationID,
	})
	if err != nil {
		return models.Member{}, err
	}

	return s.GetMember(ctx, row.ID)
}

func (s Service) UpdateMember(ctx context.Context, ID uuid.UUID, params member.UpdateParams) (models.Member, error) {
	q := s.membersQ().FilterByID(ID)
	if params.Position != nil {
		if *params.Position == "" {
			q.UpdatePosition(sql.NullString{Valid: false})
		} else {
			q = q.UpdatePosition(sql.NullString{String: *params.Position, Valid: true})
		}
	}
	if params.Label != nil {
		if *params.Label == "" {
			q.UpdateLabel(sql.NullString{Valid: false})
		} else {
			q = q.UpdateLabel(sql.NullString{String: *params.Label, Valid: true})
		}
	}

	row, err := q.UpdateOne(ctx)
	if err != nil {
		return models.Member{}, err
	}

	return s.GetMember(ctx, row.ID)
}

func (s Service) GetMember(ctx context.Context, memberID uuid.UUID) (models.Member, error) {
	row, err := s.membersQ().FilterByID(memberID).GetWithUserData(ctx)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return models.Member{}, nil
	case err != nil:
		return models.Member{}, fmt.Errorf("getting member by id: %w", err)
	}

	return MemberWithUserData(row), nil
}

func (s Service) GetMemberByAccountAndAgglomeration(
	ctx context.Context,
	accountID, agglomerationID uuid.UUID,
) (models.Member, error) {
	row, err := s.membersQ().
		FilterByAccountID(accountID).
		FilterByAgglomerationID(agglomerationID).
		GetWithUserData(ctx)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return models.Member{}, nil
	case err != nil:
		return models.Member{}, fmt.Errorf("getting member by account and agglomeration: %w", err)
	}

	return MemberWithUserData(row), nil
}

func (s Service) GetMembers(
	ctx context.Context,
	filter member.FilterParams,
	offset uint,
	limit uint,
) (pagi.Page[[]models.Member], error) {
	q := s.membersQ()
	if filter.AgglomerationID != nil {
		q = q.FilterByAgglomerationID(*filter.AgglomerationID)
	}
	if filter.AccountID != nil {
		q = q.FilterByAccountID(*filter.AccountID)
	}
	if filter.Username != nil {
		q = q.FilterByUsername(*filter.Username)
	}
	if filter.BestMatch != nil {
		q = q.FilterLikeUsername(*filter.BestMatch)
	}
	if filter.RoleID != nil {
		q = q.FilterRoleID(*filter.RoleID)
	}
	if filter.PermissionCode != nil {
		q = q.FilterByPermissionCode(*filter.PermissionCode)
	}
	if filter.RoleRankUp != nil {
		q = q.FilterByRoleRankUp(*filter.RoleRankUp)
	}
	if filter.RoleRankDown != nil {
		q = q.FilterByRoleRankDown(*filter.RoleRankDown)
	}
	if filter.Label != nil {
		q = q.FilterLikeLabel(*filter.Label)
	}
	if filter.Position != nil {
		q = q.FilterLikePosition(*filter.Position)
	}

	rows, err := q.Page(limit, offset).SelectWithUserData(ctx)
	if err != nil {
		return pagi.Page[[]models.Member]{}, fmt.Errorf("filtering members: %w", err)
	}

	total, err := q.Count(ctx)
	if err != nil {
		return pagi.Page[[]models.Member]{}, fmt.Errorf("counting members: %w", err)
	}

	collection := make([]models.Member, 0, len(rows))
	for _, row := range rows {
		collection = append(collection, MemberWithUserData(row))
	}

	return pagi.Page[[]models.Member]{
		Data:  collection,
		Page:  uint(offset/limit) + 1,
		Size:  uint(len(collection)),
		Total: uint(total),
	}, nil
}

func (s Service) DeleteMember(ctx context.Context, memberID uuid.UUID) error {
	return s.membersQ().FilterByID(memberID).Delete(ctx)
}

func (s Service) DeleteMembershipsByAccountID(ctx context.Context, accountID uuid.UUID) error {
	return s.membersQ().FilterByAccountID(accountID).Delete(ctx)
}

func (s Service) CheckMemberHavePermission(
	ctx context.Context,
	memberID uuid.UUID,
	permissionCode string,
) (bool, error) {
	have, err := s.membersQ().
		FilterByID(memberID).
		FilterByPermissionCode(permissionCode).Exists(ctx)
	if err != nil {
		return false, fmt.Errorf("checking member have permission: %w", err)
	}

	return have, nil
}

func (s Service) CanInteract(ctx context.Context, firstMemberID, secondMemberID uuid.UUID) (bool, error) {
	res, err := s.membersQ().CanInteract(ctx, firstMemberID, secondMemberID)
	if err != nil {
		return false, fmt.Errorf("checking first member can interact: %w", err)
	}

	return res, nil
}

func MemberWithUserData(db pgdb.MemberWithUserData) models.Member {
	return models.Member{
		ID:              db.ID,
		AccountID:       db.AccountID,
		AgglomerationID: db.AgglomerationID,
		Position:        db.Position,
		Label:           db.Label,
		Username:        db.Username,
		Pseudonym:       db.Pseudonym,
		Official:        db.Official,
		CreatedAt:       db.CreatedAt,
		UpdatedAt:       db.UpdatedAt,
	}
}
