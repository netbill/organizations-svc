package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/models"
	"github.com/netbill/organizations-svc/internal/core/modules/member"
	"github.com/netbill/organizations-svc/internal/repository/pgdb"
	"github.com/netbill/pagi"
	"github.com/pkg/errors"
)

func (r Repository) CreateMember(
	ctx context.Context,
	accountID, organizationID uuid.UUID,
) (models.Member, error) {
	row, err := r.orgMembersQ(ctx).Insert(ctx, pgdb.InsertMemberParams{
		AccountID:      accountID,
		OrganizationID: organizationID,
	})
	if err != nil {
		return models.Member{}, fmt.Errorf(
			"failed to create member for organization ID %s and account ID %s cause: %w",
			organizationID, accountID, err,
		)
	}

	return r.GetMember(ctx, row.ID)
}

func (r Repository) UpdateMember(
	ctx context.Context, ID uuid.UUID, params member.UpdateParams) (models.Member, error) {
	q := r.orgMembersQ(ctx).FilterByID(ID)
	if params.Position != nil {
		if *params.Position == "" {
			q.UpdatePosition(pgtype.Text{Valid: false})
		} else {
			q = q.UpdatePosition(pgtype.Text{String: *params.Position, Valid: true})
		}
	}
	if params.Label != nil {
		if *params.Label == "" {
			q.UpdateLabel(pgtype.Text{Valid: false})
		} else {
			q = q.UpdateLabel(pgtype.Text{String: *params.Label, Valid: true})
		}
	}

	row, err := q.UpdateOne(ctx)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return models.Member{}, errx.ErrorMemberNotFound.Raise(
			fmt.Errorf("member with ID %s not found, cause: %w", ID, err),
		)
	case err != nil:
		return models.Member{}, fmt.Errorf("failed to updating member with ID %s, cause: %w", ID, err)
	}

	return r.GetMember(ctx, row.ID)
}

func (r Repository) GetMember(ctx context.Context, memberID uuid.UUID) (models.Member, error) {
	row, err := r.orgMembersQ(ctx).FilterByID(memberID).GetWithUserData(ctx)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return models.Member{}, errx.ErrorMemberNotFound.Raise(
			fmt.Errorf("member with ID %s not found, cause: %w", memberID, err),
		)
	case err != nil:
		return models.Member{}, fmt.Errorf("failed to getting member by id: %w", err)
	}

	return MemberWithUserData(row), nil
}

func (r Repository) GetMemberByAccountAndOrganization(
	ctx context.Context,
	accountID, organizationID uuid.UUID,
) (models.Member, error) {
	row, err := r.orgMembersQ(ctx).
		FilterByAccountID(accountID).
		FilterByOrganizationID(organizationID).
		GetWithUserData(ctx)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return models.Member{}, errx.ErrorMemberNotFound.Raise(
			fmt.Errorf("member with account ID %s and organization ID %s not found, cause: %w",
				accountID, organizationID, err),
		)
	case err != nil:
		return models.Member{}, fmt.Errorf("failed to getting member by account ID and organization ID: %w", err)
	}

	return MemberWithUserData(row), nil
}

func (r Repository) MemberExists(ctx context.Context, accountID, organizationID uuid.UUID) (bool, error) {
	exists, err := r.orgMembersQ(ctx).
		FilterByAccountID(accountID).
		FilterByOrganizationID(organizationID).
		Exists(ctx)
	if err != nil {
		return false, fmt.Errorf("checking member existence by account ID and organization ID: %w", err)
	}

	return exists, nil
}

func (r Repository) GetMembers(
	ctx context.Context,
	filter member.FilterParams,
	limit uint,
	offset uint,
) (pagi.Page[[]models.Member], error) {
	q := r.orgMembersQ(ctx)
	if filter.OrganizationID != nil {
		q = q.FilterByOrganizationID(*filter.OrganizationID)
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

	if limit == 0 {
		limit = 10
	}

	rows, err := q.Page(limit, offset).SelectWithUserData(ctx)
	if err != nil {
		return pagi.Page[[]models.Member]{}, fmt.Errorf(
			"failed to get members with filter, cause: %w", err,
		)
	}

	total, err := q.Count(ctx)
	if err != nil {
		return pagi.Page[[]models.Member]{}, fmt.Errorf(
			"failed to count members with filter, cause: %w", err,
		)
	}

	collection := make([]models.Member, 0, len(rows))
	for _, row := range rows {
		collection = append(collection, MemberWithUserData(row))
	}

	return pagi.Page[[]models.Member]{
		Data:  collection,
		Page:  uint(offset/limit) + 1,
		Size:  uint(len(collection)),
		Total: total,
	}, nil
}

func (r Repository) DeleteMember(ctx context.Context, memberID uuid.UUID) error {
	err := r.orgMembersQ(ctx).FilterByID(memberID).Delete(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete member with ID %s cause: %w", memberID, err)
	}

	return nil
}

func (r Repository) DeleteMembersByAccountID(ctx context.Context, accountID uuid.UUID) error {
	err := r.orgMembersQ(ctx).FilterByAccountID(accountID).Delete(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete members with account ID %s cause: %w", accountID, err)
	}

	return nil
}

//func (s Repository) CanInteract(ctx context.Context, firstMemberID, secondMemberID uuid.UUID) (bool, error) {
//	res, err := s.orgMembersQ(ctx).CanInteract(ctx, firstMemberID, secondMemberID)
//	if err != nil {
//		return false, fmt.Errorf("checking first member can interact: %w", err)
//	}
//
//	return res, nil
//}

func MemberWithUserData(db pgdb.OrganizationMemberWithUserData) models.Member {
	mem := models.Member{
		ID:             db.ID,
		AccountID:      db.AccountID,
		OrganizationID: db.OrganizationID,
		Username:       db.Username,
		Official:       db.Official,
		CreatedAt:      db.CreatedAt,
		UpdatedAt:      db.UpdatedAt,
	}
	if db.Pseudonym.Valid {
		mem.Pseudonym = &db.Pseudonym.String
	}
	if db.Icon.Valid {
		mem.Icon = &db.Icon.String
	}
	if db.Position.Valid {
		mem.Position = &db.Position.String
	}
	if db.Label.Valid {
		mem.Label = &db.Label.String
	}

	return mem
}
