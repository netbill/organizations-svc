package agglomeration

import (
	"context"

	"github.com/google/uuid"
	"github.com/umisto/cities-svc/internal/domain/entity"
)

type UpdateParams struct {
	Name *string `json:"name,omitempty"`
	Icon *string `json:"icon,omitempty"`
}

func (s Service) UpdateAgglomeration(ctx context.Context, ID uuid.UUID, params UpdateParams) (entity.Agglomeration, error) {
	return s.repo.UpdateAgglomeration(ctx, ID, params)
}
