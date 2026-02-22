package profile

import (
	"context"

	"github.com/netbill/organizations-svc/internal/core/domain"
)

func (m *Module) Create(
	ctx context.Context,
	profile domain.Profile,
) (domain.Profile, error) {
	createdProfile, err := m.repo.CreateProfile(ctx, profile)
	if err != nil {
		return domain.Profile{}, err
	}

	return createdProfile, nil
}
