package role

import (
	"context"

	"github.com/netbill/organizations-svc/internal/core/models"
)

func (m *Module) GetAllPermissions(ctx context.Context) ([]models.OrgRolePermission, error) {
	res, err := m.repo.GetAllPermissions(ctx)
	if err != nil {
		return nil, err
	}

	return res, nil
}
