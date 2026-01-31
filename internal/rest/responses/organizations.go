package responses

import (
	"net/http"

	"github.com/netbill/organizations-svc/internal/core/models"
	"github.com/netbill/organizations-svc/resources"
	"github.com/netbill/restkit/pagi"
)

func Organization(organization models.Organization) resources.Organization {
	return resources.Organization{
		Data: resources.OrganizationData{
			Id:   organization.ID,
			Type: "organization",
			Attributes: resources.OrganizationDataAttributes{
				Status:    organization.Status,
				Name:      organization.Name,
				CreatedAt: organization.CreatedAt,
				UpdatedAt: organization.UpdatedAt,
			},
		},
	}
}

func Organizations(r *http.Request, page pagi.Page[[]models.Organization]) resources.OrganizationsCollection {
	data := make([]resources.OrganizationData, len(page.Data))
	for i, ag := range page.Data {
		data[i] = Organization(ag).Data
	}

	links := pagi.BuildPageLinks(r, page.Page, page.Size, page.Total)

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

func UpdateOrganizationSession(uploadLinks models.UpdateOrganizationMedia, organization models.Organization) resources.UpdateOrganizationLinks {
	return resources.UpdateOrganizationLinks{
		Data: resources.UpdateOrganizationLinksData{
			Id:   uploadLinks.UploadSessionID,
			Type: "update_organization_session",
			Attributes: resources.UpdateOrganizationLinksDataAttributes{
				UploadToken:     uploadLinks.UploadToken,
				IconGetUrl:      uploadLinks.Links.IconGetURL,
				IconUploadUrl:   uploadLinks.Links.IconUploadURL,
				BannerGetUrl:    uploadLinks.Links.BannerGetURL,
				BannerUploadUrl: uploadLinks.Links.BannerUploadURL,
			},
			Relationships: resources.UpdateOrganizationLinksDataRelationships{
				Organization: &resources.UpdateOrganizationLinksDataRelationshipsOrganization{
					Data: resources.UpdateOrganizationLinksDataRelationshipsOrganizationData{
						Id:   organization.ID,
						Type: "organization",
					},
				},
			},
		},
		Included: []resources.OrganizationData{
			Organization(organization).Data,
		},
	}
}
