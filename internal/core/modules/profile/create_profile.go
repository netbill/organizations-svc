package profile

import (
	"context"

	"github.com/netbill/organizations-svc/internal/core/models"
)

func (s Service) CreateProfile(ctx context.Context, profile models.Profile) (models.Profile, error) {
	createdProfile, err := s.repo.CreateProfile(ctx, profile)
	if err != nil {
		return models.Profile{}, err
	}

	return createdProfile, nil
}
