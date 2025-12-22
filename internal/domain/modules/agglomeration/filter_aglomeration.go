package agglomeration

import (
	"context"

	"github.com/umisto/cities-svc/internal/domain/entity"
	"github.com/umisto/pagi"
)

type FilterParams struct {
	NameLike *string `json:"name_likes,omitempty"`
	Status   *string `json:"status,omitempty"`
}

func (s Service) FilterAgglomerations(
	ctx context.Context,
	params FilterParams,
	pagination pagi.Params,
) (pagi.Page[entity.Agglomeration], error) {
	return s.repo.FilterAgglomerations(ctx, params, pagination)
}
