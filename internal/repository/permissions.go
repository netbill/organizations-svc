package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/umisto/agglomerations-svc/internal/domain/models"
	"github.com/umisto/agglomerations-svc/internal/repository/pgdb"
)

func (s Service) GetPermission(ctx context.Context, ID uuid.UUID) (models.Permission, error) {
	res, err := s.permissionsQ().FilterByID(ID).Get(ctx)
	if err != nil {
		return models.Permission{}, err
	}

	return Permission(res), nil
}

type FilterPermissionsParams struct {
	Description *string
	Code        *string
}

func Permission(p pgdb.Permission) models.Permission {
	return models.Permission{
		ID:          p.ID,
		Code:        p.Code,
		Description: p.Description,
	}
}
