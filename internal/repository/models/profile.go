package models

import (
	"github.com/umisto/cities-svc/internal/domain/entity"
	"github.com/umisto/cities-svc/internal/repository/pgdb"
)

func Profile(row pgdb.Profile) entity.Profile {
	return entity.Profile{
		AccountID: row.AccountID,
		Username:  row.Username,
		Official:  row.Official,
		Pseudonym: row.Pseudonym,
	}
}
