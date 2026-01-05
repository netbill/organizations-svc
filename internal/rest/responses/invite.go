package responses

import (
	"net/http"

	"github.com/netbill/pagi"
	"github.com/netbill/places-svc/internal/core/models"
	"github.com/netbill/places-svc/resources"
)

func Invite(mod models.Invite) resources.Invite {
	return resources.Invite{
		Data: resources.InviteData{
			Id:   mod.ID,
			Type: "invite",
			Attributes: resources.InviteDataAttributes{
				OrganizationId: mod.OrganizationID,
				AccountId:      mod.AccountID,
				Status:         mod.Status,
				CreatedAt:      mod.CreatedAt,
				ExpiresAt:      mod.ExpiresAt,
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
