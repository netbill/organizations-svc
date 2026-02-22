package role

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/domain"
)

func (m *Module) GetByID(
	ctx context.Context,
	roleID uuid.UUID,
) (domain.Role, error) {
	role, err := m.repo.GetRole(ctx, roleID)
	if err != nil {
		return domain.Role{}, err
	}

	return role, nil
}
