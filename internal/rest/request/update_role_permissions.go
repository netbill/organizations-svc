package request

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/netbill/organizations-svc/resources"
)

func UpdateRolePermissions(r *http.Request) (req resources.UpdateRolePermissions, err error) {
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		return req, err
	}

	errs := validation.Errors{
		"data/id":         validation.Validate(req.Data.Id, validation.Required),
		"data/type":       validation.Validate(req.Data.Type, validation.Required, validation.In("update_role_permissions")),
		"data/attributes": validation.Validate(req.Data.Attributes, validation.Required),
	}

	if req.Data.Id.String() != chi.URLParam(r, "role_id") {
		errs["data/id"] = validation.NewError("mismatch", "query role_id and body data/id do not match")
	}

	return req, errs.Filter()
}
