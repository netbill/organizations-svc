package responses

import (
	"net/http"

	"github.com/netbill/organizations-svc/internal/core/models"
	"github.com/netbill/organizations-svc/pkg/resources"
	"github.com/netbill/restkit/pagi"
)

func Invite(mod models.Invite) resources.Invite {
	return resources.Invite{
		Data: resources.InviteData{
			Id:   mod.ID,
			Type: "invite",
			Attributes: resources.InviteDataAttributes{
				Status:    mod.Status,
				CreatedAt: mod.CreatedAt,
				ExpiresAt: mod.ExpiresAt,
			},
			Relationships: resources.MemberDataRelationships{
				Organization: resources.MemberDataRelationshipsOrganization{
					Data: resources.MemberDataRelationshipsOrganizationData{
						Id:   mod.OrganizationID,
						Type: "organization",
					},
				},
				Profile: resources.MemberDataRelationshipsProfile{
					Data: resources.MemberDataRelationshipsProfileData{
						Id:   mod.AccountID,
						Type: "profile",
					},
				},
			},
		},
	}
}

func Invites(r *http.Request, mods pagi.Page[[]models.Invite]) resources.InvitesCollection {
	data := make([]resources.InviteData, len(mods.Data))
	for i, mod := range mods.Data {
		data[i] = Invite(mod).Data
	}

	links := pagi.BuildPageLinks(r, mods.Page, mods.Size, mods.Total)

	return resources.InvitesCollection{
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
