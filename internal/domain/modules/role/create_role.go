package role

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/umisto/agglomerations-svc/internal/domain/errx"
	"github.com/umisto/agglomerations-svc/internal/domain/models"
)

type CreateParams struct {
	AgglomerationID uuid.UUID `json:"agglomeration_id"`
	Rank            uint      `json:"rank"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	Color           string    `json:"color"`
}

func (s Service) CreateRole(
	ctx context.Context,
	accountID uuid.UUID,
	params CreateParams,
) (role models.Role, err error) {
	initiator, err := s.getInitiator(ctx, accountID, params.AgglomerationID)
	if err != nil {
		return role, err
	}

	if err = s.checkPermissionsToManageRole(ctx, initiator.ID, params.Rank); err != nil {
		return models.Role{}, err
	}

	if err = s.repo.Transaction(ctx, func(ctx context.Context) error {
		role, err = s.repo.CreateRole(ctx, params)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to create role: %w", err),
			)
		}

		err = s.messenger.WriteRoleCreated(ctx, role)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to write role to member %w", err),
			)
		}

		return nil
	}); err != nil {
		return models.Role{}, err
	}

	return role, nil
}
