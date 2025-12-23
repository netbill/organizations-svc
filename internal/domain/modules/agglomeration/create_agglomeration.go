package agglomeration

import (
	"context"
	"fmt"

	"github.com/umisto/cities-svc/internal/domain/entity"
	"github.com/umisto/cities-svc/internal/domain/errx"
)

func (s Service) CreateAgglomeration(ctx context.Context, name string) (entity.Agglomeration, error) {
	res, err := s.repo.CreateAgglomeration(ctx, name)
	if err != nil {
		return entity.Agglomeration{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to create agglomeration: %w", err),
		)
	}

	err = s.messager.WriteAgglomerationCreated(ctx, res)
	if err != nil {
		return entity.Agglomeration{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to publish agglomeration create event: %w", err),
		)
	}

	return res, nil
}
