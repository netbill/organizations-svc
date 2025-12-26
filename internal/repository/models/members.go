package models

import (
	"github.com/umisto/cities-svc/internal/domain/entity"
	"github.com/umisto/cities-svc/internal/repository/pgdb"
)

func MemberWithUserData(db pgdb.MemberWithUserData) entity.Member {
	roles := make([]entity.MemberRole, len(db.Roles))
	for i, r := range db.Roles {
		roles[i] = entity.MemberRole{
			RoleID: r.RoleID,
			Head:   r.Head,
			Rank:   r.Rank,
			Name:   r.Name,
		}
	}

	return entity.Member{
		ID:              db.ID,
		AccountID:       db.AccountID,
		AgglomerationID: db.AgglomerationID,
		Position:        db.Position,
		Label:           db.Label,
		Username:        db.Username,
		Pseudonym:       db.Pseudonym,
		Official:        db.Official,
		Roles:           roles,
		CreatedAt:       db.CreatedAt,
		UpdatedAt:       db.UpdatedAt,
	}
}
