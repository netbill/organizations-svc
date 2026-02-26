package profile

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/models"
)

func (m *Module) GetByID(
	ctx context.Context,
	accountID uuid.UUID,
) (models.Profile, error) {
	profile, err := m.repo.GetProfileByAccountID(ctx, accountID)
	if err != nil {
		return models.Profile{}, err
	}

	return profile, nil
}

func (m *Module) GetByIDs(
	ctx context.Context,
	accountIDs []uuid.UUID,
) ([]models.Profile, error) {
	profiles, err := m.repo.GetProfilesByAccountIDs(ctx, accountIDs)
	if err != nil {
		return nil, err
	}

	return profiles, nil
}

func (m *Module) GetByUsername(
	ctx context.Context,
	username string,
) (models.Profile, error) {
	profile, err := m.repo.GetProfileByUsername(ctx, username)
	if err != nil {
		return models.Profile{}, err
	}

	return profile, nil
}

func (m *Module) ExistsByUsername(
	ctx context.Context,
	username string,
) (bool, error) {
	return m.repo.ExistsProfileByUsername(ctx, username)
}

func (m *Module) ExistsByAccountID(
	ctx context.Context,
	accountID uuid.UUID,
) (bool, error) {
	return m.repo.ExistsProfileByAccountID(ctx, accountID)
}
