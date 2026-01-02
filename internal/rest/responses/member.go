package responses

import (
	"net/http"

	"github.com/netbill/organizations-svc/internal/domain/models"
	"github.com/netbill/organizations-svc/resources"
	"github.com/netbill/pagi"
)

func Member(mod models.Member) resources.Member {
	return resources.Member{
		Data: resources.MemberData{
			Id:   mod.ID,
			Type: "member",
			Attributes: resources.MemberDataAttributes{
				OrganizationId: mod.OrganizationID,
				AccountId:      mod.AccountID,
				Position:       mod.Position,
				Label:          mod.Label,
				Username:       mod.Username,
				Official:       mod.Official,
				CreatedAt:      mod.CreatedAt,
				UpdatedAt:      mod.UpdatedAt,
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
