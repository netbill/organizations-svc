package responses

import (
	"net/http"
	"strconv"

	"github.com/netbill/organizations-svc/internal/domain/models"
	"github.com/netbill/organizations-svc/resources"
	"github.com/netbill/pagi"
)

func Organization(organization models.Organization) resources.Organization {
	return resources.Organization{
		Data: resources.OrganizationData{
			Id:   organization.ID,
			Type: "organization",
			Attributes: resources.OrganizationDataAttributes{
				Status:    organization.Status,
				Name:      organization.Name,
				Icon:      organization.Icon,
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

	links := BuildPageLinks(r, page.Page, page.Size, page.Total)

	return resources.OrganizationsCollection{
		Data:  data,
		Links: links,
	}
}

func BuildPageLinks(r *http.Request, page, size, total uint) resources.PaginationData {
	if page == 0 {
		page = 1
	}
	if size == 0 {
		size = 20
	}

	lastPage := uint(1)
	if total > 0 {
		lastPage = (total + size - 1) / size
		if lastPage == 0 {
			lastPage = 1
		}
	}

	self := buildURLWithPage(r, page, size)

	var first *string
	if page != 1 {
		v := buildURLWithPage(r, 1, size)
		first = &v
	}

	var last *string
	if page != lastPage {
		v := buildURLWithPage(r, lastPage, size)
		last = &v
	}

	var prev *string
	if page > 1 {
		v := buildURLWithPage(r, page-1, size)
		prev = &v
	}

	var next *string
	if page < lastPage {
		v := buildURLWithPage(r, page+1, size)
		next = &v
	}

	return resources.PaginationData{
		Self:  self,
		First: first,
		Last:  last,
		Prev:  prev,
		Next:  next,
	}
}

func buildURLWithPage(r *http.Request, page, size uint) string {
	u := *r.URL
	q := u.Query()

	q.Set("page[number]", strconv.FormatUint(uint64(page), 10))
	q.Set("page[size]", strconv.FormatUint(uint64(size), 10))

	u.RawQuery = q.Encode()
	return u.String()
}
