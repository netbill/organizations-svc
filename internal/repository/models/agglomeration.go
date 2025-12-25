package models

import (
	"github.com/umisto/cities-svc/internal/domain/entity"
	"github.com/umisto/cities-svc/internal/repository/pgdb"
)

func AgglomerationRow(db pgdb.Agglomeration) entity.Agglomeration {
	return entity.Agglomeration{
		ID:        db.ID,
		Status:    string(db.Status),
		Name:      db.Name,
		Icon:      db.Icon.String,
		CreatedAt: db.CreatedAt,
		UpdatedAt: db.UpdatedAt,
	}
}
