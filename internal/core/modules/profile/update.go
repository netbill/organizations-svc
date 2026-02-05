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
	Avatar    *string   `json:"avatar"`
	Official  bool      `json:"official"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (m *Module) Update(
	ctx context.Context,
	accountID uuid.UUID,
	params UpdateParams,
) (models.Profile, error) {
	updatedProfile, err := m.repo.UpdateProfile(ctx, accountID, params)
	if err != nil {
		return models.Profile{}, err
	}

	return updatedProfile, nil
}
