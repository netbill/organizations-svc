package profile

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/places-svc/internal/core/errx"
	"github.com/netbill/places-svc/internal/core/models"
)

func (s Service) UpdateUsername(ctx context.Context, accountID uuid.UUID, username string) (models.Profile, error) {
	res, err := s.repo.UpdateUsername(ctx, accountID, username)
	if err != nil {
		return models.Profile{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to update username for accountID %s: %w", accountID, err),
		)
	}

	return res, nil
}
