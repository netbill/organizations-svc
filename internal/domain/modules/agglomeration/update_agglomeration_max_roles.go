package agglomeration

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/umisto/agglomerations-svc/internal/domain/errx"
	"github.com/umisto/agglomerations-svc/internal/domain/models"
)

func (s Service) UpdateAgglomerationMaxRoles(
	ctx context.Context,
	agglomerationID uuid.UUID,
	maxRoles uint,
) (models.Agglomeration, error) {
	agglo, err := s.repo.UpdateAgglomerationMaxRoles(ctx, agglomerationID, maxRoles)
	if err != nil {
		return models.Agglomeration{}, errx.ErrorInternal.Raise(
			fmt.Errorf("failed to update agglomeration max roles: %w", err),
		)
	}

	return agglo, nil
}
