package agglomeration

import (
	"context"

	"github.com/umisto/cities-svc/internal/domain/entity"
)

func (s Service) CreateAgglomeration(ctx context.Context, name string) (entity.Agglomeration, error) {
	return s.repo.CreateAgglomeration(ctx, name)
}
