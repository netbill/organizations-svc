package responses

import (
	"net/http"

	"github.com/netbill/organizations-svc/internal/core/domain"
	resources2 "github.com/netbill/organizations-svc/pkg/resources"
	"github.com/netbill/restkit/pagi"
)

func Member(mod domain.Member) resources2.Member {
	return resources2.Member{
		Data: resources2.MemberData{
			Id:   mod.ID,
			Type: "member",
			Attributes: resources2.MemberDataAttributes{
				OrganizationId: mod.OrganizationID,
				AccountId:      mod.AccountID,
				Head:           mod.Head,
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

func Members(r *http.Request, mods pagi.Page[[]domain.Member]) resources2.MemberCollection {
	data := make([]resources2.MemberData, len(mods.Data))
	for i, mod := range mods.Data {
		data[i] = Member(mod).Data
	}

	links := pagi.BuildPageLinks(r, mods.Page, mods.Size, mods.Total)

	return resources2.MemberCollection{
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
