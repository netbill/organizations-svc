package responses

import (
	"net/http"

	"github.com/netbill/organizations-svc/internal/core/models"
	resources "github.com/netbill/organizations-svc/pkg/resources"
	"github.com/netbill/restkit/pagi"
)

func Member(mod models.Member) resources.Member {
	return resources.Member{
		Data: resources.MemberData{
			Id:   mod.ID,
			Type: "member",
			Attributes: resources.MemberDataAttributes{
				Head:      mod.Head,
				Position:  mod.Position,
				Label:     mod.Label,
				CreatedAt: mod.CreatedAt,
				UpdatedAt: mod.UpdatedAt,
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
		Included: []resources.MemberIncludedInner{
			{
				ProfileData: &resources.ProfileData{
					Type: "profile",
					Id:   mod.AccountID,
					Attributes: resources.ProfileAttributes{
						Username:  mod.Username,
						Official:  mod.Official,
						Pseudonym: mod.Pseudonym,
						AvatarKey: mod.AvatarKey,
						Version:   mod.Version,
						CreatedAt: mod.CreatedAt,
						UpdatedAt: mod.UpdatedAt,
					},
				},
			},
		},
	}
}

func Members(r *http.Request, mods pagi.Page[[]models.Member]) resources.MemberCollection {
	data := make([]resources.MemberData, len(mods.Data))
	for i, mod := range mods.Data {
		data[i] = Member(mod).Data
	}

	links := pagi.BuildPageLinks(r, mods.Page, mods.Size, mods.Total)

	return resources.MemberCollection{
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
