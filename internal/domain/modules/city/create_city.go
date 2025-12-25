package city

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/paulmach/orb"
	"github.com/umisto/cities-svc/internal/domain/entity"
	"github.com/umisto/cities-svc/internal/domain/errx"
)

type CreateParams struct {
	AgglomerationID *uuid.UUID
	Name            string
	Slug            *string
	Icon            *string
	Banner          *string
	Point           orb.Point
}

func (s Service) CreateCity(ctx context.Context, params CreateParams) (city entity.City, err error) {
	if params.AgglomerationID != nil {
		if _, err = s.checkAgglomerationIsActiveAndExists(ctx, *params.AgglomerationID); err != nil {
			return entity.City{}, err
		}
	}

	if params.Slug != nil {
		err = s.checkSlugIsAvailable(ctx, *params.Slug)
		if err != nil {
			return entity.City{}, err
		}
	}

	err = s.repo.Transaction(ctx, func(ctx context.Context) error {
		err = s.repo.CreateCity(ctx, params)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to create city: %w", err),
			)
		}

		err = s.messanger.CreateCity(ctx, city)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to send create city message: %w", err),
			)
		}

		return nil
	})

	return s.GetCity(ctx, city.ID)
}
