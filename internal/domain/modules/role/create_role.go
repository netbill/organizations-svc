package role

import (
	"context"

	"github.com/google/uuid"
	"github.com/umisto/cities-svc/internal/domain/entity"
)

type CreateParams struct {
	AgglomerationID uuid.UUID `json:"agglomeration_id"`
	Head            bool      `json:"head"`
	Rank            uint      `json:"rank"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	Color           string    `json:"color"`
}

func (s Service) CreateRole(ctx context.Context, params CreateParams) (entity.Role, error) {
	return s.repo.CreateRole(ctx, params)
}

func (s Service) CreateRoleByUser(
	ctx context.Context,
	accountID uuid.UUID,
	params CreateParams,
) (entity.Role, error) {
	if err := s.CheckPermissionsToManageRole(ctx, accountID, params.AgglomerationID, params.Rank); err != nil {
		return entity.Role{}, err
	}

	if err := s.CheckPermissionsToManageRole(ctx, accountID, params.AgglomerationID, params.Rank); err != nil {
		return entity.Role{}, err
	}

	return s.CreateRole(ctx, params)
}
