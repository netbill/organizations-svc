package agglomeration

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/umisto/cities-svc/internal/domain/errx"
)

func (s Service) DeleteAgglomeration(ctx context.Context, ID uuid.UUID) error {
	if err := s.repo.Transaction(ctx, func(ctx context.Context) error {
		err := s.repo.DeleteAgglomeration(ctx, ID)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to delete agglomeration: %w", err),
			)
		}

		err = s.messenger.WriteAgglomerationDeleted(ctx, ID)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to publish agglomeration delete event: %w", err),
			)
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
