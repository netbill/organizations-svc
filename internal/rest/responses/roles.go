package responses

import (
	"net/http"

	"github.com/netbill/organizations-svc/internal/core/models"
	"github.com/netbill/organizations-svc/resources"
	"github.com/netbill/pagi"
)

func Role(mod models.Role, perms map[models.Permission]bool) resources.Role {
	res := resources.Role{
		Data: resources.RoleData{
			Id:   mod.ID,
			Type: "role",
			Attributes: resources.RoleDataAttributes{
				OrganizationId: mod.OrganizationID,
				Head:           mod.Head,
				Rank:           mod.Rank,
				Name:           mod.Name,
				Description:    mod.Description,
				Color:          mod.Color,
				CreatedAt:      mod.CreatedAt,
				UpdatedAt:      mod.UpdatedAt,
			},
		},
	}

	if perms != nil {
		ps := make([]resources.RoleDataRelationshipsPermissionsInner, 0, len(perms))

		for perm, has := range perms {
			ps = append(ps, resources.RoleDataRelationshipsPermissionsInner{
				Id:          perm.ID,
				Code:        perm.Code,
				Description: perm.Description,
				Enabled:     has,
			})
		}

		res.Data.Relationships = &resources.RoleDataRelationships{
			Permissions: ps,
		}
	}

	return res
}

func Roles(r *http.Request, mods pagi.Page[[]models.Role]) resources.RolesCollection {
	data := make([]resources.RoleData, len(mods.Data))
	for i, mod := range mods.Data {
		data[i] = Role(mod, nil).Data
	}

	links := pagi.BuildPageLinks(r, mods.Page, mods.Size, mods.Total)

	return resources.RolesCollection{
		Data: data,
		Links: resources.PaginationData{
			First: links.First,
			Last:  links.Last,
			Prev:  links.Prev,
			Next:  links.Next,
			Self:  links.Self,
		},
	}
}

func RolePermissions(mods []models.Permission) resources.RolePermissions {
	result := make([]resources.RolePermissionsDataInner, len(mods))
	for i, mod := range mods {
		result[i] = resources.RolePermissionsDataInner{
			Id:          mod.ID,
			Code:        mod.Code,
			Description: mod.Description,
		}
	}

	return resources.RolePermissions{
		Data: result,
	}
}
