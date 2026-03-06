package responses

import (
	"net/http"

	"github.com/netbill/organizations-svc/internal/models"
	"github.com/netbill/organizations-svc/internal/rest/scope"
	resources "github.com/netbill/organizations-svc/pkg/resources"
	"github.com/netbill/restkit/pagi"
)

func Organization(
	r *http.Request,
	organization models.Organization,
) resources.Organization {
	return resources.Organization{
		Data: organizationData(r, organization),
	}
}

func organizationData(r *http.Request, organization models.Organization) resources.OrganizationData {
	res := resources.OrganizationData{
		Id:   organization.ID,
		Type: "organization",
		Attributes: resources.OrganizationDataAttributes{
			Status:    organization.Status,
			Name:      organization.Name,
			Version:   organization.Version,
			CreatedAt: organization.CreatedAt,
			UpdatedAt: organization.UpdatedAt,
		},
	}
	if organization.IconKey != nil {
		url := scope.ResolverURL(r, *organization.IconKey)
		res.Attributes.IconUrl = &url
	}
	if organization.BannerKey != nil {
		url := scope.ResolverURL(r, *organization.BannerKey)
		res.Attributes.BannerUrl = &url
	}

	return res
}

type organizationCollectionResponse struct {
	data     []resources.OrganizationData
	included []resources.MemberData
}

type OrgCollectionOption func(*organizationCollectionResponse)

func WithOrganizationMembers(r *http.Request, members []models.Member) OrgCollectionOption {
	return func(res *organizationCollectionResponse) {
		for _, member := range members {
			res.included = append(res.included, memberData(r, member))
		}
	}
}

func Organizations(
	r *http.Request,
	page pagi.Page[[]models.Organization],
	opts ...OrgCollectionOption,
) resources.OrganizationsCollection {
	data := make([]resources.OrganizationData, len(page.Data))
	for i, org := range page.Data {
		data[i] = organizationData(r, org)
	}

	res := &organizationCollectionResponse{
		data: data,
	}
	for _, opt := range opts {
		opt(res)
	}

	links := pagi.BuildPageLinks(r, page.Page, page.Size, page.Total)

	return resources.OrganizationsCollection{
		Data:     res.data,
		Included: res.included,
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
	r *http.Request,
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
			organizationData(r, organization),
		},
	}
}
