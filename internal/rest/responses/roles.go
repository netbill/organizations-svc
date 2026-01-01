package responses

import (
	"net/http"

	"github.com/umisto/agglomerations-svc/internal/domain/models"
	"github.com/umisto/agglomerations-svc/resources"
	"github.com/umisto/pagi"
)

func Role(mod models.Role) resources.Role {
	return resources.Role{
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
}

func Roles(r *http.Request, mods pagi.Page[[]models.Role]) resources.RolesCollection {
	data := make([]resources.RoleData, len(mods.Data))
	for i, mod := range mods.Data {
		data[i] = Role(mod).Data
	}

	links := BuildPageLinks(r, mods.Page, mods.Size, mods.Total)

	return resources.RolesCollection{
		Data:  data,
		Links: links,
	}
}
