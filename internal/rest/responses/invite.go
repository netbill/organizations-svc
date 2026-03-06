package responses

import (
	"fmt"
	"net/http"

	"github.com/netbill/organizations-svc/internal/models"
	"github.com/netbill/organizations-svc/pkg/resources"
	"github.com/netbill/restkit/pagi"
)

type inviteResponse struct {
	data     resources.InviteData
	included []resources.InviteIncludedInner
}

type InviteOption func(*inviteResponse)

func WithInviteOrganization(r *http.Request, organization models.Organization) InviteOption {
	return func(res *inviteResponse) {
		org := organizationData(r, organization)
		res.included = append(res.included, resources.InviteIncludedInner{
			OrganizationData: &org,
		})
	}
}

func WithInviteProfile(r *http.Request, profile models.Profile) InviteOption {
	return func(res *inviteResponse) {
		prof := profileData(r, profile)
		res.included = append(res.included, resources.InviteIncludedInner{
			ProfileData: &prof,
		})
	}
}

func Invite(
	r *http.Request,
	mod models.Invite,
	opts ...InviteOption,
) resources.Invite {
	res := &inviteResponse{
		data: inviteData(r, mod),
	}
	for _, opt := range opts {
		opt(res)
	}

	return resources.Invite{
		Data:     res.data,
		Included: res.included,
	}
}

func inviteData(r *http.Request, invite models.Invite) resources.InviteData {
	return resources.InviteData{
		Id:   invite.ID,
		Type: "invite",
		Attributes: resources.InviteDataAttributes{
			Status:    invite.Status,
			UpdatedAt: invite.UpdatedAt,
			CreatedAt: invite.CreatedAt,
			ExpiresAt: invite.ExpiresAt,
		},
		Relationships: resources.MemberDataRelationships{
			Organization: resources.MemberDataRelationshipsOrganization{
				Data: resources.MemberDataRelationshipsOrganizationData{
					Id:   invite.OrganizationID,
					Type: "organization",
				},
			},
			Profile: resources.MemberDataRelationshipsProfile{
				Data: resources.MemberDataRelationshipsProfileData{
					Id:   invite.AccountID,
					Type: "profile",
				},
			},
		},
	}
}

type inviteCollectionResponse struct {
	data     []resources.InviteData
	included []resources.InviteIncludedInner
}

type InvitesCollectionOption func(*inviteCollectionResponse)

func WithCollectionInvitesProfiles(r *http.Request, profiles []models.Profile) InvitesCollectionOption {
	return func(res *inviteCollectionResponse) {
		for _, model := range profiles {
			prof := profileData(r, model)
			res.included = append(res.included, resources.InviteIncludedInner{
				ProfileData: &prof,
			})
		}
	}
}

func WithCollectionInvitesOrganization(r *http.Request, organization models.Organization) InvitesCollectionOption {
	return func(res *inviteCollectionResponse) {
		org := organizationData(r, organization)
		res.included = append(res.included, resources.InviteIncludedInner{
			OrganizationData: &org,
		})
	}
}

func WithCollectionInvitesOrganizations(r *http.Request, organizations []models.Organization) InvitesCollectionOption {
	return func(res *inviteCollectionResponse) {
		for _, org := range organizations {
			orgData := organizationData(r, org)
			res.included = append(res.included, resources.InviteIncludedInner{
				OrganizationData: &orgData,
			})
		}
	}
}

func Invites(
	r *http.Request,
	page pagi.Page[[]models.Invite],
	opts ...InvitesCollectionOption,
) resources.InvitesCollection {
	data := make([]resources.InviteData, len(page.Data))
	for i, mod := range page.Data {
		data[i] = inviteData(r, mod)
	}

	resp := &inviteCollectionResponse{
		data: data,
	}
	for _, opt := range opts {
		opt(resp)
	}

	links := pagi.BuildPageLinks(r, page.Page, page.Size, page.Total)

	return resources.InvitesCollection{
		Data:     data,
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

func deduplicateIncluded(included []resources.InviteIncludedInner) []resources.InviteIncludedInner {
	seen := make(map[string]struct{})
	result := make([]resources.InviteIncludedInner, 0, len(included))

	for _, item := range included {
		key := includeKey(item)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		result = append(result, item)
	}
	return result
}

func includeKey(item resources.InviteIncludedInner) string {
	if item.OrganizationData != nil {
		return fmt.Sprintf("organization:%s", item.OrganizationData.Id)
	}
	if item.ProfileData != nil {
		return fmt.Sprintf("profile:%s", item.ProfileData.Id)
	}
	return ""
}
