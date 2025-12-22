package agglomeration

import (
	"context"

	"github.com/google/uuid"
)

func (s Service) DeleteAgglomeration(ctx context.Context, ID uuid.UUID) error {
	return s.repo.DeleteAgglomeration(ctx, ID)
}
