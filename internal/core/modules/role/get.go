package role

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/models"
)

func (m *Module) GetByID(
	ctx context.Context,
	roleID uuid.UUID,
) (models.Role, error) {
	role, err := m.repo.GetRole(ctx, roleID)
	if err != nil {
		return models.Role{}, err
	}

	return role, nil
}
