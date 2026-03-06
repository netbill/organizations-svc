package requests

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/netbill/organizations-svc/pkg/resources"
	"github.com/netbill/restkit"
)

func CreateOrganization(r *http.Request) (req resources.CreateOrganization, err error) {
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		err = restkit.NewDecodeError("body", err)
		return
	}

	errs := validation.Errors{
		"data/type":       validation.Validate(req.Data.Type, validation.Required, validation.In("organization")),
		"data/attributes": validation.Validate(req.Data.Attributes, validation.Required),
	}

	return req, errs.Filter()
}

func UpdateOrganization(r *http.Request) (req resources.UpdateOrganization, err error) {
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		err = restkit.NewDecodeError("body", err)
		return
	}

	errs := validation.Errors{
		"data/id":         validation.Validate(req.Data.Id, validation.Required),
		"data/type":       validation.Validate(req.Data.Type, validation.Required, validation.In("organization")),
		"data/attributes": validation.Validate(req.Data.Attributes, validation.Required),
	}

	if chi.URLParam(r, "organization_id") != req.Data.Id.String() {
		errs["data/id"] = validation.NewError("mismatch", "query organization_id and body data/id do not match")
	}

	return req, errs.Filter()
}

func DeleteUploadOrgMedia(r *http.Request) (req resources.DeleteUploadOrgMedia, err error) {
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		return req, err
	}

	errs := validation.Errors{
		"data/id":         validation.Validate(req.Data.Id, validation.Required),
		"data/type":       validation.Validate(req.Data.Type, validation.Required, validation.In("organization")),
		"data/attributes": validation.Validate(req.Data.Attributes, validation.Required),
	}

	if chi.URLParam(r, "organization_id") != req.Data.Id.String() {
		errs["data/id"] = validation.NewError("mismatch", "query organization_id and body data/id do not match")
	}

	return req, errs.Filter()
}
