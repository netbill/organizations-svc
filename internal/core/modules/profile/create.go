package profile

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/domain"
)

type UpdateParams struct {
	Username  string    `json:"username"`
	Pseudonym *string   `json:"pseudonym"`
	AvatarKey *string   `json:"avatar"`
	Official  bool      `json:"official"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (m *Module) Update(
	ctx context.Context,
	accountID uuid.UUID,
	params UpdateParams,
) (domain.Profile, error) {
	updatedProfile, err := m.repo.UpdateProfile(ctx, accountID, params)
	if err != nil {
		return domain.Profile{}, err
	}

	return updatedProfile, nil
}
