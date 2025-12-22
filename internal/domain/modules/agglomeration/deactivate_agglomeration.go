package agglomeration

import (
	"context"

	"github.com/google/uuid"
)

func (s Service) DeactivateAgglomeration(ctx context.Context, ID uuid.UUID) error {
	return s.repo.DeactivateAgglomeration(ctx, ID)
}
