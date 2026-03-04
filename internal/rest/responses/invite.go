package responses

import (
	"net/http"

	"github.com/netbill/organizations-svc/internal/core/models"
	"github.com/netbill/organizations-svc/pkg/resources"
	"github.com/netbill/restkit/pagi"
)

type inviteResponse struct {
	invite   resources.InviteData
	included []resources.InviteIncludedInner
}

type InviteOption func(*inviteResponse)

func WithInviteOrganization(organization models.Organization) InviteOption {
	return func(r *inviteResponse) {
		inner := organizationData(organization)
		r.included = append(r.included, resources.InviteIncludedInner{
			OrganizationData: &inner,
		})
	}
}

func WithInviteProfile(model models.Profile) InviteOption {
	return func(r *inviteResponse) {
		inner := profileData(model)
		r.included = append(r.included, resources.InviteIncludedInner{
			ProfileData: &inner,
		})
	}
}

func Invite(mod models.Invite, opts ...InviteOption) resources.Invite {
	r := &inviteResponse{
		invite: inviteData(mod),
	}
	for _, opt := range opts {
		opt(r)
	}
	return resources.Invite{
		Data:     r.invite,
		Included: r.included,
	}
}

func inviteData(invite models.Invite) resources.InviteData {
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

func WithCollectionInvitesProfiles(profile []models.Profile) InvitesCollectionOption {
	return func(r *inviteCollectionResponse) {
		for _, model := range profile {
			inner := profileData(model)
			r.included = append(r.included, resources.InviteIncludedInner{
				ProfileData: &inner,
			})
		}
	}
}

func WithCollectionInvitesOrganization(organization models.Organization) InvitesCollectionOption {
	return func(r *inviteCollectionResponse) {
		inner := organizationData(organization)
		r.included = append(r.included, resources.InviteIncludedInner{
			OrganizationData: &inner,
		})
	}
}

func WithCollectionInvitesOrganizations(organizations []models.Organization) InvitesCollectionOption {
	return func(r *inviteCollectionResponse) {
		for _, model := range organizations {
			inner := organizationData(model)
			r.included = append(r.included, resources.InviteIncludedInner{
				OrganizationData: &inner,
			})
		}
	}
}

func Invites(r *http.Request, mods pagi.Page[[]models.Invite], opts ...InvitesCollectionOption) resources.InvitesCollection {
	data := make([]resources.InviteData, len(mods.Data))
	for i, mod := range mods.Data {
		data[i] = inviteData(mod)
	}

	links := pagi.BuildPageLinks(r, mods.Page, mods.Size, mods.Total)

	resp := &inviteCollectionResponse{
		data: data,
	}
	for _, opt := range opts {
		opt(resp)
	}

	return resources.InvitesCollection{
		Data:     resp.data,
		Included: deduplicateInviteIncluded(resp.included),
		Links: resources.PaginationData{
			First: links.First,
			Last:  links.Last,
			Prev:  links.Prev,
			Next:  links.Next,
			Self:  links.Self,
		},
	}
}

func deduplicateInviteIncluded(items []resources.InviteIncludedInner) []resources.InviteIncludedInner {
	seen := make(map[string]struct{})
	result := make([]resources.InviteIncludedInner, 0, len(items))

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
