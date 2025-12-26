package models

import (
	"github.com/umisto/cities-svc/internal/domain/entity"
	"github.com/umisto/cities-svc/internal/repository/pgdb"
)

func Permission(p pgdb.Permission) entity.Permission {
	return entity.Permission{
		ID:          p.ID,
		Code:        entity.CodeRolePermission(p.Code),
		Description: p.Description,
	}

}
