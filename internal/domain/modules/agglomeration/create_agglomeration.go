package agglomeration

import (
	"context"
	"fmt"

	"github.com/umisto/cities-svc/internal/domain/entity"
	"github.com/umisto/cities-svc/internal/domain/errx"
)

func (s Service) CreateAgglomeration(ctx context.Context, name string) (agglo entity.Agglomeration, err error) {
	if err = s.repo.Transaction(ctx, func(ctx context.Context) error {
		agglo, err = s.repo.CreateAgglomeration(ctx, name)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to create agglomeration: %w", err),
			)
		}

		err = s.messenger.WriteAgglomerationCreated(ctx, agglo)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to publish agglomeration create event: %w", err),
			)
		}

		return nil
	}); err != nil {
		return entity.Agglomeration{}, err
	}

	return agglo, nil
}
