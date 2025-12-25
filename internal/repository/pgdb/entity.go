package pgdb

import (
	"github.com/umisto/cities-svc/internal/domain/entity"
)

func (m Member) ToEntity() entity.Member {
	ent := entity.Member{
		ID:              m.ID,
		AccountID:       m.AccountID,
		AgglomerationID: m.AgglomerationID,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
	}

	if m.Position.Valid {
		ent.Position = &m.Position.String
	}
	if m.Label.Valid {
		ent.Label = &m.Label.String
	}

	return ent
}

func (r GetMemberRow) ToEntity() entity.Member {
	ent := entity.Member{
		ID:              r.MemberID,
		AccountID:       r.AccountID,
		AgglomerationID: r.AgglomerationID,
		CreatedAt:       r.CreatedAt,
		UpdatedAt:       r.UpdatedAt,
		Username:        r.Username,
		Official:        r.Official,
	}

	if r.Position.Valid {
		ent.Position = &r.Position.String
	}
	if r.Label.Valid {
		ent.Label = &r.Label.String
	}

	if r.Pseudonym.Valid {
		ent.Pseudonym = &r.Pseudonym.String
	}

	var roles RoleDTOs
	if err := roles.Scan(r.Roles); err == nil {
		for _, role := range roles {
			ent.Roles = append(ent.Roles, entity.MemberRole{
				RoleID: role.RoleID,
				Head:   role.Head,
				Rank:   role.Rank,
				Name:   role.Name,
			})
		}
	}

	return ent
}

func (r FilterMembersRow) ToEntity() entity.Member {
	ent := entity.Member{
		ID:              r.MemberID,
		AccountID:       r.AccountID,
		AgglomerationID: r.AgglomerationID,
		CreatedAt:       r.CreatedAt,
		UpdatedAt:       r.UpdatedAt,
		Username:        r.Username,
		Official:        r.Official,
	}

	if r.Position.Valid {
		ent.Position = &r.Position.String
	}
	if r.Label.Valid {
		ent.Label = &r.Label.String
	}

	if r.Pseudonym.Valid {
		ent.Pseudonym = &r.Pseudonym.String
	}

	var roles RoleDTOs
	if err := roles.Scan(r.Roles); err == nil {
		for _, role := range roles {
			ent.Roles = append(ent.Roles, entity.MemberRole{
				RoleID: role.RoleID,
				Head:   role.Head,
				Rank:   role.Rank,
				Name:   role.Name,
			})
		}
	}

	return ent
}

func (r Role) ToEntity() entity.Role {
	return entity.Role{
		ID:              r.ID,
		AgglomerationID: r.AgglomerationID,
		Head:            r.Head,
		Editable:        r.Editable,
		Rank:            uint(r.Rank),
		Name:            r.Name,
		CreatedAt:       r.CreatedAt,
		UpdatedAt:       r.UpdatedAt,
	}
}

func (r GetRoleRow) ToEntity() entity.Role {
	ent := entity.Role{
		ID:              r.ID,
		AgglomerationID: r.AgglomerationID,
		Head:            r.Head,
		Editable:        r.Editable,
		Rank:            uint(r.Rank),
		Name:            r.Name,
		CreatedAt:       r.CreatedAt,
		UpdatedAt:       r.UpdatedAt,
	}

	var perms PermissionDTOs
	if err := perms.Scan(r.Permissions); err == nil {
		for _, perm := range perms {
			ent.Permissions = append(ent.Permissions, entity.Permission{
				ID:          perm.PermissionID,
				Code:        entity.CodeRolePermission(perm.Code),
				Description: perm.Description,
			})
		}
	}

	return ent
}

func (r UpdateRoleRankRow) ToEntity() entity.Role {
	return entity.Role{
		ID:              r.ID,
		AgglomerationID: r.AgglomerationID,
		Head:            r.Head,
		Editable:        r.Editable,
		Rank:            uint(r.Rank),
		Name:            r.Name,
		CreatedAt:       r.CreatedAt,
		UpdatedAt:       r.UpdatedAt,
	}
}
func (r FilterRolesRow) ToEntity() entity.Role {
	ent := entity.Role{
		ID:              r.ID,
		AgglomerationID: r.AgglomerationID,
		Head:            r.Head,
		Editable:        r.Editable,
		Rank:            uint(r.Rank),
		Name:            r.Name,
		CreatedAt:       r.CreatedAt,
		UpdatedAt:       r.UpdatedAt,
	}

	var perms PermissionDTOs
	if err := perms.Scan(r.Permissions); err == nil {
		for _, perm := range perms {
			ent.Permissions = append(ent.Permissions, entity.Permission{
				ID:          perm.PermissionID,
				Code:        entity.CodeRolePermission(perm.Code),
				Description: perm.Description,
			})
		}
	}

	return ent
}

func (p Permission) ToEntity() entity.Permission {
	return entity.Permission{
		ID:          p.ID,
		Code:        entity.CodeRolePermission(p.Code),
		Description: p.Description,
	}
}
