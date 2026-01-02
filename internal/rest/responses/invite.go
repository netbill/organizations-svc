package responses

import (
	"net/http"

	"github.com/umisto/agglomerations-svc/internal/domain/models"
	"github.com/umisto/agglomerations-svc/resources"
	"github.com/umisto/pagi"
)

func Invite(mod models.Invite) resources.Invite {
	return resources.Invite{
		Data: resources.InviteData{
			Id:   mod.ID,
			Type: "invite",
			Attributes: resources.InviteDataAttributes{
				AgglomerationId: mod.AgglomerationID,
				AccountId:       mod.AccountID,
				Status:          mod.Status,
				CreatedAt:       mod.CreatedAt,
				ExpiresAt:       mod.ExpiresAt,
			},
		},
	}
}

func Invites(r *http.Request, mods pagi.Page[[]models.Invite]) resources.InvitesCollection {
	data := make([]resources.InviteData, len(mods.Data))
	for i, mod := range mods.Data {
		data[i] = Invite(mod).Data
	}

	links := BuildPageLinks(r, mods.Page, mods.Size, mods.Total)

	return resources.InvitesCollection{
		Data:  data,
		Links: links,
	}
}
