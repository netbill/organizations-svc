package responses

import (
	"net/http"

	"github.com/netbill/organizations-svc/internal/core/models"
	resources2 "github.com/netbill/organizations-svc/pkg/resources"
	"github.com/netbill/restkit/pagi"
)

func Invite(mod models.Invite) resources2.Invite {
	return resources2.Invite{
		Data: resources2.InviteData{
			Id:   mod.ID,
			Type: "invite",
			Attributes: resources2.InviteDataAttributes{
				OrganizationId: mod.OrganizationID,
				AccountId:      mod.AccountID,
				Status:         mod.Status,
				CreatedAt:      mod.CreatedAt,
				ExpiresAt:      mod.ExpiresAt,
			},
		},
	}
}

func Invites(r *http.Request, mods pagi.Page[[]models.Invite]) resources2.InvitesCollection {
	data := make([]resources2.InviteData, len(mods.Data))
	for i, mod := range mods.Data {
		data[i] = Invite(mod).Data
	}

	links := pagi.BuildPageLinks(r, mods.Page, mods.Size, mods.Total)

	return resources2.InvitesCollection{
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
