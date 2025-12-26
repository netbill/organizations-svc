package models

import (
	"github.com/umisto/cities-svc/internal/domain/entity"
	"github.com/umisto/cities-svc/internal/repository/pgdb"
)

func Role(r pgdb.Role) entity.Role {
	return entity.Role{
		ID:              r.ID,
		AgglomerationID: r.AgglomerationID,
		Head:            r.Head,
		Editable:        r.Editable,
		Rank:            r.Rank,
		Name:            r.Name,
		CreatedAt:       r.CreatedAt,
		UpdatedAt:       r.UpdatedAt,
	}
}
