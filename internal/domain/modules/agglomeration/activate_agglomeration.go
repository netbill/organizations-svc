package agglomeration

import (
	"context"

	"github.com/google/uuid"
)

func (s Service) ActivateAgglomeration(ctx context.Context, ID uuid.UUID) error {
	return s.repo.ActivateAgglomeration(ctx, ID)
}
