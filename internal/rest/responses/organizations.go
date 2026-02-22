package responses

import (
	"net/http"

	"github.com/netbill/organizations-svc/internal/core/domain"
	resources2 "github.com/netbill/organizations-svc/pkg/resources"
	"github.com/netbill/restkit/pagi"
)

func Organization(organization domain.Organization) resources2.Organization {
	return resources2.Organization{
		Data: resources2.OrganizationData{
			Id:   organization.ID,
			Type: "organization",
			Attributes: resources2.OrganizationDataAttributes{
				Status:    organization.Status,
				Name:      organization.Name,
				CreatedAt: organization.CreatedAt,
				UpdatedAt: organization.UpdatedAt,
			},
		},
	}
}

func Organizations(r *http.Request, page pagi.Page[[]domain.Organization]) resources2.OrganizationsCollection {
	data := make([]resources2.OrganizationData, len(page.Data))
	for i, ag := range page.Data {
		data[i] = Organization(ag).Data
	}

	links := pagi.BuildPageLinks(r, page.Page, page.Size, page.Total)

	return resources2.OrganizationsCollection{
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

func UploadOrganizationMediaLinks(organization domain.Organization, uploadLinks domain.UploadOrgMediaLinks) resources2.UploadOrgMediaLinks {
	return resources2.UploadOrgMediaLinks{
		Data: resources2.UploadOrgMediaLinksData{
			Id:   organization.ID,
			Type: "update_organization_session",
			Attributes: resources2.UploadOrgMediaLinksDataAttributes{
				Icon: resources2.UploadResourcesLink{
					Key:        uploadLinks.Icon.Key,
					UploadUrl:  uploadLinks.Icon.UploadURL,
					PreloadUrl: uploadLinks.Icon.PreloadUrl,
				},
				Banner: resources2.UploadResourcesLink{
					Key:        uploadLinks.Banner.Key,
					UploadUrl:  uploadLinks.Banner.UploadURL,
					PreloadUrl: uploadLinks.Banner.PreloadUrl,
				},
			},
			Relationships: resources2.UploadOrgMediaLinksDataRelationships{
				Organization: &resources2.UploadOrgMediaLinksDataRelationshipsOrganization{
					Data: resources2.UploadOrgMediaLinksDataRelationshipsOrganizationData{
						Id:   organization.ID,
						Type: "organization",
					},
				},
			},
		},
		Included: []resources2.OrganizationData{
			Organization(organization).Data,
		},
	}
}
