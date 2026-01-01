package role

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/umisto/cities-svc/internal/domain/errx"
	"github.com/umisto/cities-svc/internal/domain/models"
)

type CreateParams struct {
	AgglomerationID uuid.UUID `json:"agglomeration_id"`
	Rank            uint      `json:"rank"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	Color           string    `json:"color"`
}

func (s Service) CreateRole(ctx context.Context, params CreateParams) (role models.Role, err error) {
	if err = s.repo.Transaction(ctx, func(ctx context.Context) error {
		role, err = s.repo.CreateRole(ctx, params)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to create role: %w", err),
			)
		}

		return s.messenger.WriteRoleCreated(ctx, role)
	}); err != nil {
		return models.Role{}, errx.ErrorInternal.Raise(
			fmt.Errorf("transaction failed when creating role: %w", err),
		)
	}

	return role, nil
}

func (s Service) CreateRoleByUser(
	ctx context.Context,
	accountID uuid.UUID,
	params CreateParams,
) (models.Role, error) {
	if err := s.CheckPermissionsToManageRole(ctx, accountID, params.AgglomerationID, params.Rank); err != nil {
		return models.Role{}, err
	}

	if err := s.CheckPermissionsToManageRole(ctx, accountID, params.AgglomerationID, params.Rank); err != nil {
		return models.Role{}, err
	}

	return s.CreateRole(ctx, params)
}
