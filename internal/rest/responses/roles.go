package responses

import (
	"net/http"

	"github.com/netbill/organizations-svc/internal/core/domain"
	resources2 "github.com/netbill/organizations-svc/pkg/resources"
	"github.com/netbill/restkit/pagi"
)

func Role(role domain.Role, perms *domain.OrgRolePermissionsWithDetailsForRole) resources2.Role {
	res := resources2.Role{
		Data: resources2.RoleData{
			Id:   role.ID,
			Type: "role",
			Attributes: resources2.RoleDataAttributes{
				OrganizationId: role.OrganizationID,
				Rank:           role.Rank,
				Name:           role.Name,
				Description:    role.Description,
				Color:          role.Color,
				CreatedAt:      role.CreatedAt,
				UpdatedAt:      role.UpdatedAt,
			},
		},
	}

	if perms != nil {
		ps := make([]resources2.RoleDataIncludedPermissionsInner, 0, len(*perms))

		for _, details := range *perms {
			ps = append(ps, resources2.RoleDataIncludedPermissionsInner{
				Code:        details.Code,
				Description: details.Description,
				Enabled:     details.Enabled,
			})
		}

		res.Data.Included = &resources2.RoleDataIncluded{
			Permissions: ps,
		}
	}

	return res
}

func Roles(r *http.Request, mods pagi.Page[[]domain.Role]) resources2.RolesCollection {
	data := make([]resources2.RoleData, len(mods.Data))
	for i, mod := range mods.Data {
		data[i] = Role(mod, nil).Data
	}

	links := pagi.BuildPageLinks(r, mods.Page, mods.Size, mods.Total)

	return resources2.RolesCollection{
		Data: data,
		Links: resources2.PaginationData{
			First: links.First,
			Last:  links.Last,
			Prev:  links.Prev,
			Next:  links.Next,
			Self:  links.Self,
		},
	}
}

func RolePermissions(mods []domain.OrgRolePermission) resources2.RolePermissions {
	result := make([]resources2.RolePermissionsDataInner, len(mods))
	for i, mod := range mods {
		result[i] = resources2.RolePermissionsDataInner{
			Code:        mod.Code,
			Description: mod.Description,
		}
	}

	return resources2.RolePermissions{
		Data: result,
	}
}
