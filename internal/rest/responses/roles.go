package responses

import (
	"net/http"

	"github.com/umisto/agglomerations-svc/internal/domain/models"
	"github.com/umisto/agglomerations-svc/resources"
	"github.com/umisto/pagi"
)

func Role(mod models.Role, perms []models.Permission) resources.Role {
	res := resources.Role{
		Data: resources.RoleData{
			Id:   mod.ID,
			Type: "role",
			Attributes: resources.RoleDataAttributes{
				AgglomerationId: mod.AgglomerationID,
				Head:            mod.Head,
				Rank:            mod.Rank,
				Name:            mod.Name,
				Description:     mod.Description,
				Color:           mod.Color,
				CreatedAt:       mod.CreatedAt,
				UpdatedAt:       mod.UpdatedAt,
			},
		},
	}

	if perms != nil {
		res.Data.Relationships = &resources.RoleDataRelationships{
			Permissions: RolePermissions(perms),
		}
	}

	return res
}

func Roles(r *http.Request, mods pagi.Page[[]models.Role]) resources.RolesCollection {
	data := make([]resources.RoleData, len(mods.Data))
	for i, mod := range mods.Data {
		data[i] = Role(mod, nil).Data
	}

	links := BuildPageLinks(r, mods.Page, mods.Size, mods.Total)

	return resources.RolesCollection{
		Data:  data,
		Links: links,
	}
}

func RolePermissions(mods []models.Permission) []resources.RolePermission {
	result := make([]resources.RolePermission, len(mods))
	for i, mod := range mods {
		result[i] = resources.RolePermission{
			Id:          mod.ID,
			Code:        string(mod.Code),
			Description: mod.Description,
		}
	}

	return result
}
