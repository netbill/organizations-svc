package profile

import (
	"context"
	"fmt"

	"github.com/netbill/organizations-svc/internal/domain/errx"
	"github.com/netbill/organizations-svc/internal/domain/models"
)

func (s Service) UpsertProfile(ctx context.Context, profile models.Profile) (models.Profile, error) {
	res, err := s.repo.UpsertProfile(ctx, profile)
	if err != nil {
		return models.Profile{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to upsert profile: %w", err),
		)
	}

	return res, nil
}
