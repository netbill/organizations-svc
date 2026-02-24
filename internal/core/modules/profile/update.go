package profile

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/models"
)

type UpdateParams struct {
	Username  string    `json:"username"`
	Pseudonym *string   `json:"pseudonym"`
	AvatarKey *string   `json:"avatar"`
	Official  bool      `json:"official"`
	Version   int       `json:"version"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (m *Module) Update(
	ctx context.Context,
	accountID uuid.UUID,
	params UpdateParams,
) (models.Profile, error) {
	profile, err := m.repo.GetProfileByAccountID(ctx, accountID)
	if err != nil {
		return models.Profile{}, err
	}
	if profile.Version >= params.Version {
		return profile, nil
	}

	updatedProfile, err := m.repo.UpdateProfile(ctx, accountID, params)
	if err != nil {
		return models.Profile{}, err
	}

	return updatedProfile, nil
}
