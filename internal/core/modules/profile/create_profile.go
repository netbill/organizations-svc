package profile

import (
	"context"
	"fmt"

	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/models"
)

func (s Service) CreateProfile(ctx context.Context, profile models.Profile) (models.Profile, error) {
	createdProfile, err := s.repo.CreateProfile(ctx, profile)
	if err != nil {
		return models.Profile{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to create profile: %w", err),
		)
	}

	return createdProfile, nil
}
