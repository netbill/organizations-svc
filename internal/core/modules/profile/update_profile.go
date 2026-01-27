package profile

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/models"
)

type UpdateParams struct {
	Username  string    `json:"username"`
	Pseudonym *string   `json:"pseudonym"`
	Icon      *string   `json:"icon"`
	Official  bool      `json:"official"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (s Service) UpdateProfile(ctx context.Context, accountID uuid.UUID, params UpdateParams) (models.Profile, error) {
	updatedProfile, err := s.repo.UpdateProfile(ctx, accountID, params)
	if err != nil {
		return models.Profile{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to update profile: %w", err),
		)
	}

	return updatedProfile, nil
}
