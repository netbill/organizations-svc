package role

import (
	"context"

	"github.com/netbill/organizations-svc/internal/core/domain"
)

func (m *Module) GetAllPermissions(ctx context.Context) ([]domain.OrgRolePermission, error) {
	res, err := m.repo.GetAllPermissions(ctx)
	if err != nil {
		return nil, err
	}

	return res, nil
}
