package agglomeration

import (
	"context"

	"github.com/umisto/cities-svc/internal/domain/entity"
	"github.com/umisto/pagi"
)

type FilterParams struct {
	NameContains *string `json:"name_contains,omitempty"`
	Active       *bool   `json:"active,omitempty"`
}

func (s Service) FilterAgglomerations(
	ctx context.Context,
	params FilterParams,
	pagination pagi.Params,
) (pagi.Page[entity.Agglomeration], error) {
	return s.repo.FilterAgglomerations(ctx, params, pagination)
}
