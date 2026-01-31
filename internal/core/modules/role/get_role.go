package role

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/models"
	"github.com/netbill/restkit/pagi"
)

func (s Service) GetRole(ctx context.Context, roleID uuid.UUID) (models.Role, error) {
	role, err := s.repo.GetRole(ctx, roleID)
	if err != nil {
		return models.Role{}, err
	}

	return role, nil
}

func (s Service) GetRoleWithPermissions(ctx context.Context, accountID, roleID uuid.UUID) (models.Role, map[models.Permission]bool, error) {
	role, err := s.GetRole(ctx, roleID)
	if err != nil {
		return models.Role{}, nil, err
	}

	_, err = s.getInitiator(ctx, accountID, role.OrganizationID)
	if err != nil {
		return models.Role{}, nil, err
	}

	permissions, err := s.repo.GetRolePermissions(ctx, roleID)
	if err != nil {
		return models.Role{}, nil, err
	}

	return role, permissions, nil
}

type FilterParams struct {
	OrganizationID *uuid.UUID
	RolesID        *[]uuid.UUID
	Head           *bool
	Rank           *int
	Name           *string
}

func (s Service) GetRoles(
	ctx context.Context,
	params FilterParams,
	limit, offset uint,
) (pagi.Page[[]models.Role], error) {
	res, err := s.repo.GetRoles(ctx, params, limit, offset)
	if err != nil {
		return pagi.Page[[]models.Role]{}, err
	}

	return res, nil
}
