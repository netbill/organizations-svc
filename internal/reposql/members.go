package repol

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/umisto/cities-svc/internal/domain/entity"
	"github.com/umisto/cities-svc/internal/domain/modules/memeber"
	"github.com/umisto/cities-svc/internal/reposql/pgdbsql"

	"github.com/umisto/nilx"
	"github.com/umisto/pagi"
)

func (s Service) CreateMember(ctx context.Context, accountID, agglomerationID uuid.UUID) (entity.Member, error) {
	res, err := s.sql(ctx).CreateMember(ctx, pgdbsql.CreateMemberParams{
		AccountID:       accountID,
		AgglomerationID: agglomerationID,
	})

	if err != nil {
		return entity.Member{}, err
	}

	return res.ToEntity(), nil
}

func (s Service) UpdateMember(ctx context.Context, member) (entity.Member, error) {
	res, err := s.sql(ctx).UpdateMember(ctx, pgdbsql.UpdateMemberParams{
		ID:       member.ID,
		Position: nilx.String(member.Position),
		Label:    nilx.String(member.Label),
	})
	if err != nil {
		return entity.Member{}, err
	}

	return res.ToEntity(), nil
}

func (s Service) GetMember(ctx context.Context, memberID uuid.UUID) (entity.Member, error) {
	res, err := s.sql(ctx).GetMember(ctx, memberID)
	if err != nil {
		return entity.Member{}, err
	}

	return res.ToEntity(), nil
}

func (s Service) FilterMembers(
	ctx context.Context,
	filter memeber.FilterParams,
	pagination pagi.Params,
) (pagi.Page[entity.Member], error) {
	params := pgdbsql.FilterMembersParams{
		AgglomerationID: nilx.UUID(filter.AgglomerationID),
		Username:        nilx.String(filter.Username),
		AccountID:       nilx.UUID(filter.AccountID),
		RoleID:          nilx.UUID(filter.RoleID),
		PermissionCode:  nilx.String(filter.PermissionCode),
	}

	if pagination.Cursor != nil {
		usernameCursor, ok := pagination.Cursor["username"]
		if !ok || usernameCursor == "" {
			return pagi.Page[entity.Member]{}, fmt.Errorf("missing username in pagination cursor")
		}
		params.CursorUsername = sql.NullString{String: usernameCursor, Valid: true}

		idCursor, ok := pagination.Cursor["id"]
		if !ok || idCursor == "" {
			return pagi.Page[entity.Member]{}, fmt.Errorf("missing id in pagination cursor")
		}

		afterID, err := uuid.Parse(idCursor)
		if err != nil {
			return pagi.Page[entity.Member]{}, fmt.Errorf("invalid id in pagination cursor: %w", err)
		}
		params.CursorMemberID = uuid.NullUUID{UUID: afterID, Valid: true}
	}

	limit := pagi.CalculateLimit(pagination.Limit, 50, 100)
	params.Limit = int32(limit)

	members, err := s.sql(ctx).FilterMembers(ctx, params)
	if err != nil {
		return pagi.Page[entity.Member]{}, err
	}

	count, err := s.sql(ctx).CountMembers(ctx, pgdbsql.CountMembersParams{
		AgglomerationID: nilx.UUID(filter.AgglomerationID),
		Username:        nilx.String(filter.Username),
		AccountID:       nilx.UUID(filter.AccountID),
		RoleID:          nilx.UUID(filter.RoleID),
		PermissionCode:  nilx.String(filter.PermissionCode),
	})
	if err != nil {
		return pagi.Page[entity.Member]{}, err
	}

	entities := make([]entity.Member, len(members))
	for i, m := range members {
		entities[i] = m.ToEntity()
	}

	var nextCursor map[string]string
	if len(members) == limit {
		lastMember := members[len(members)-1]
		nextCursor = map[string]string{
			"username": lastMember.Username,
			"id":       lastMember.MemberID.String(),
		}
	}

	return pagi.Page[entity.Member]{
		Data:       entities,
		Total:      int(count),
		NextCursor: nextCursor,
	}, nil
}

func (s Service) DeleteMember(ctx context.Context, memberID uuid.UUID) error {
	return s.sql(ctx).DeleteMember(ctx, memberID)
}
