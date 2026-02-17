package request

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/netbill/organizations-svc/pkg/resources"
)

func DeleteUploadOrgIcon(r *http.Request) (req resources.DeleteUploadOrgIcon, err error) {
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		return req, err
	}

	errs := validation.Errors{
		"data/id":         validation.Validate(req.Data.Id, validation.Required),
		"data/type":       validation.Validate(req.Data.Type, validation.Required, validation.In("delete_organization_upload_icon")),
		"data/attributes": validation.Validate(req.Data.Attributes, validation.Required),
	}

	if chi.URLParam(r, "organization_id") != req.Data.Id.String() {
		errs["data/id"] = validation.NewError("mismatch", "query organization_id and body data/id do not match")
	}

	return req, errs.Filter()
}
