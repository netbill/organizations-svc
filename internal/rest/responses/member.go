package responses

import (
	"net/http"

	"github.com/umisto/agglomerations-svc/internal/domain/models"
	"github.com/umisto/agglomerations-svc/resources"
	"github.com/umisto/pagi"
)

func Member(mod models.Member) resources.Member {
	return resources.Member{
		Data: resources.MemberData{
			Id:   mod.ID,
			Type: "member",
			Attributes: resources.MemberDataAttributes{
				AgglomerationId: mod.AgglomerationID,
				AccountId:       mod.AccountID,
				Position:        mod.Position,
				Label:           mod.Label,
				Username:        mod.Username,
				Official:        mod.Official,
				CreatedAt:       mod.CreatedAt,
				UpdatedAt:       mod.UpdatedAt,
			},
		},
	}
}

func Members(r *http.Request, mods pagi.Page[[]models.Member]) resources.MemberCollection {
	data := make([]resources.MemberData, len(mods.Data))
	for i, mod := range mods.Data {
		data[i] = Member(mod).Data
	}

	links := BuildPageLinks(r, mods.Page, mods.Size, mods.Total)

	return resources.MemberCollection{
		Data:  data,
		Links: links,
	}
}
