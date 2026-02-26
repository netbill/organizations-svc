package responses

import (
	"net/http"

	"github.com/netbill/organizations-svc/internal/core/models"
	resources "github.com/netbill/organizations-svc/pkg/resources"
	"github.com/netbill/restkit/pagi"
)

type memberResponse struct {
	member   resources.MemberData
	included []resources.MemberIncludedInner
}

type MemberOption func(*memberResponse)

func WithMemberOrganization(organization models.Organization) MemberOption {
	return func(r *memberResponse) {
		inner := organizationData(organization)
		r.included = append(r.included, resources.MemberIncludedInner{
			OrganizationData: &inner,
		})
	}
}

func WithMemberProfile(profile models.Profile) MemberOption {
	return func(r *memberResponse) {
		inner := profileData(profile)
		r.included = append(r.included, resources.MemberIncludedInner{
			ProfileData: &inner,
		})
	}
}

func Member(mod models.Member, opts ...MemberOption) resources.Member {
	r := &memberResponse{
		member: memberData(mod),
	}
	for _, opt := range opts {
		opt(r)
	}
	return resources.Member{
		Data:     r.member,
		Included: r.included,
	}
}

func memberData(mod models.Member) resources.MemberData {
	return resources.MemberData{
		Id:   mod.ID,
		Type: "member",
		Attributes: resources.MemberDataAttributes{
			Head:      mod.Head,
			Position:  mod.Position,
			Label:     mod.Label,
			Version:   mod.Version,
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
	}
}

type memberCollectionResponse struct {
	data     []resources.MemberData
	included []resources.MemberIncludedInner
}

type MembersCollectionOption func(*memberCollectionResponse)

func WithCollectionMembersProfiles(profiles []models.Profile) MembersCollectionOption {
	return func(r *memberCollectionResponse) {
		for _, model := range profiles {
			inner := profileData(model)
			r.included = append(r.included, resources.MemberIncludedInner{
				ProfileData: &inner,
			})
		}
	}
}

func WithCollectionMembersOrganization(organization models.Organization) MembersCollectionOption {
	return func(r *memberCollectionResponse) {
		inner := organizationData(organization)
		r.included = append(r.included, resources.MemberIncludedInner{
			OrganizationData: &inner,
		})
	}
}

func Members(r *http.Request, mods pagi.Page[[]models.Member], opts ...MembersCollectionOption) resources.MemberCollection {
	data := make([]resources.MemberData, len(mods.Data))
	for i, mod := range mods.Data {
		data[i] = memberData(mod)
	}

	links := pagi.BuildPageLinks(r, mods.Page, mods.Size, mods.Total)

	resp := &memberCollectionResponse{
		data: data,
	}
	for _, opt := range opts {
		opt(resp)
	}

	return resources.MemberCollection{
		Data:     resp.data,
		Included: deduplicateIncluded(resp.included),
		Links: resources.PaginationData{
			First: links.First,
			Last:  links.Last,
			Prev:  links.Prev,
			Next:  links.Next,
			Self:  links.Self,
		},
	}
}

func deduplicateIncluded(items []resources.MemberIncludedInner) []resources.MemberIncludedInner {
	seen := make(map[string]struct{})
	result := make([]resources.MemberIncludedInner, 0, len(items))

	for _, item := range items {
		var key string
		switch {
		case item.ProfileData != nil:
			key = "profile:" + item.ProfileData.Id.String()
		case item.OrganizationData != nil:
			key = "organization:" + item.OrganizationData.Id.String()
		default:
			continue
		}

		if _, exists := seen[key]; !exists {
			seen[key] = struct{}{}
			result = append(result, item)
		}
	}

	return result
}
