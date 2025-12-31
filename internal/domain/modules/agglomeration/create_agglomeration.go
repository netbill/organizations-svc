package agglomeration

import (
	"context"
	"fmt"

	"github.com/umisto/cities-svc/internal/domain/errx"
	"github.com/umisto/cities-svc/internal/domain/models"
)

type CreateParams struct {
	Name string
	Icon *string
}

func (s Service) CreateAgglomeration(ctx context.Context, params CreateParams) (agglo models.Agglomeration, err error) {
	if err = s.repo.Transaction(ctx, func(ctx context.Context) error {
		agglo, err = s.repo.CreateAgglomeration(ctx, params)
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
		return models.Agglomeration{}, err
	}

	return agglo, nil
}
