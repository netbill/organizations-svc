package responses

import (
	"net/http"

	"github.com/netbill/organizations-svc/internal/core/models"
	resources "github.com/netbill/organizations-svc/pkg/resources"
	"github.com/netbill/restkit/pagi"
)

func Organization(organization models.Organization) resources.Organization {
	return resources.Organization{
		Data: organizationData(organization),
	}
}

func organizationData(organization models.Organization) resources.OrganizationData {
	return resources.OrganizationData{
		Id:   organization.ID,
		Type: "organization",
		Attributes: resources.OrganizationDataAttributes{
			Status:    organization.Status,
			Name:      organization.Name,
			IconKey:   organization.IconKey,
			BannerKey: organization.BannerKey,
			Version:   organization.Version,
			CreatedAt: organization.CreatedAt,
			UpdatedAt: organization.UpdatedAt,
		},
	}
}

type organizationResponse struct {
	data     []resources.OrganizationData
	included []resources.MemberData
}

type OrgCollectionOption func(*organizationResponse)

func WithOrganizationMembers(members []models.Member) OrgCollectionOption {
	return func(r *organizationResponse) {
		for _, member := range members {
			r.included = append(r.included, memberData(member))
		}
	}
}

func Organizations(r *http.Request, page pagi.Page[[]models.Organization], opts ...OrgCollectionOption) resources.OrganizationsCollection {
	data := make([]resources.OrganizationData, len(page.Data))
	for i, ag := range page.Data {
		data[i] = Organization(ag).Data
	}

	links := pagi.BuildPageLinks(r, page.Page, page.Size, page.Total)

	resp := &organizationResponse{
		data: data,
	}
	for _, opt := range opts {
		opt(resp)
	}

	return resources.OrganizationsCollection{
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

func UploadOrganizationMediaLinks(
	organization models.Organization,
	uploadLinks models.UploadOrgMediaLinks,
) resources.UploadOrgMediaLinks {
	return resources.UploadOrgMediaLinks{
		Data: resources.UploadOrgMediaLinksData{
			Id:   organization.ID,
			Type: "organization_upload_links",
			Attributes: resources.UploadOrgMediaLinksDataAttributes{
				Icon: resources.UploadResourcesLink{
					Key:        uploadLinks.Icon.Key,
					UploadUrl:  uploadLinks.Icon.UploadURL,
					PreloadUrl: uploadLinks.Icon.PreloadUrl,
				},
				Banner: resources.UploadResourcesLink{
					Key:        uploadLinks.Banner.Key,
					UploadUrl:  uploadLinks.Banner.UploadURL,
					PreloadUrl: uploadLinks.Banner.PreloadUrl,
				},
			},
			Relationships: resources.UploadOrgMediaLinksDataRelationships{
				Organization: &resources.UploadOrgMediaLinksDataRelationshipsOrganization{
					Data: resources.UploadOrgMediaLinksDataRelationshipsOrganizationData{
						Id:   organization.ID,
						Type: "organization",
					},
				},
			},
		},
		Included: []resources.OrganizationData{
			organizationData(organization),
		},
	}
}
