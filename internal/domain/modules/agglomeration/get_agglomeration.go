package agglomeration

import (
	"context"

	"github.com/google/uuid"
	"github.com/umisto/cities-svc/internal/domain/entity"
)

func (s Service) GetAgglomeration(ctx context.Context, agglomerationID uuid.UUID) (entity.Agglomeration, error) {
	return s.repo.GetAgglomerationByID(ctx, agglomerationID)
}
