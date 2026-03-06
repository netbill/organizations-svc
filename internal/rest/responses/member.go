package responses

import (
	"fmt"
	"net/http"

	"github.com/netbill/organizations-svc/internal/models"
	resources "github.com/netbill/organizations-svc/pkg/resources"
	"github.com/netbill/restkit/pagi"
)

type memberResponse struct {
	data     resources.MemberData
	included []resources.MemberIncludedInner
}

type MemberOption func(*memberResponse)

func WithMemberOrganization(r *http.Request, organization models.Organization) MemberOption {
	return func(res *memberResponse) {
		org := organizationData(r, organization)
		res.included = append(res.included, resources.MemberIncludedInner{
			OrganizationData: &org,
		})
	}
}

func WithMemberProfile(r *http.Request, profile models.Profile) MemberOption {
	return func(res *memberResponse) {
		prof := profileData(r, profile)
		res.included = append(res.included, resources.MemberIncludedInner{
			ProfileData: &prof,
		})
	}
}

func Member(
	r *http.Request,
	mod models.Member,
	opts ...MemberOption,
) resources.Member {
	res := &memberResponse{
		data: memberData(r, mod),
	}
	for _, opt := range opts {
		opt(res)
	}

	return resources.Member{
		Data:     res.data,
		Included: deduplicateMemberIncluded(res.included),
	}
}

func memberData(r *http.Request, mod models.Member) resources.MemberData {
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

func WithCollectionMembersProfiles(r *http.Request, profiles []models.Profile) MembersCollectionOption {
	return func(res *memberCollectionResponse) {
		for _, profile := range profiles {
			prof := profileData(r, profile)
			res.included = append(res.included, resources.MemberIncludedInner{
				ProfileData: &prof,
			})
		}
	}
}

func WithCollectionMembersOrganization(r *http.Request, organization models.Organization) MembersCollectionOption {
	return func(res *memberCollectionResponse) {
		org := organizationData(r, organization)
		res.included = append(res.included, resources.MemberIncludedInner{
			OrganizationData: &org,
		})
	}
}

func Members(
	r *http.Request,
	page pagi.Page[[]models.Member],
	opts ...MembersCollectionOption,
) resources.MembersCollection {
	data := make([]resources.MemberData, len(page.Data))
	for i, mod := range page.Data {
		data[i] = memberData(r, mod)
	}

	res := &memberCollectionResponse{
		data: data,
	}
	for _, opt := range opts {
		opt(res)
	}

	links := pagi.BuildPageLinks(r, page.Page, page.Size, page.Total)

	return resources.MembersCollection{
		Data:     res.data,
		Included: deduplicateMemberIncluded(res.included),
		Links: resources.PaginationData{
			First: links.First,
			Last:  links.Last,
			Prev:  links.Prev,
			Next:  links.Next,
			Self:  links.Self,
		},
	}
}

func deduplicateMemberIncluded(included []resources.MemberIncludedInner) []resources.MemberIncludedInner {
	seen := make(map[string]struct{})
	result := make([]resources.MemberIncludedInner, 0, len(included))

	for _, item := range included {
		key := memberIncludeKey(item)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		result = append(result, item)
	}
	return result
}

func memberIncludeKey(item resources.MemberIncludedInner) string {
	if item.OrganizationData != nil {
		return fmt.Sprintf("organization:%s", item.OrganizationData.Id)
	}
	if item.ProfileData != nil {
		return fmt.Sprintf("profile:%s", item.ProfileData.Id)
	}
	return ""
}
