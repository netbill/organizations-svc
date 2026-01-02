package request

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/umisto/agglomerations-svc/resources"
)

func UpdateRolesRanks(r *http.Request) (req resources.UpdateRolesRanks, err error) {
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		return req, err
	}

	errs := validation.Errors{
		"data/id":         validation.Validate(req.Data.Id, validation.Required),
		"data/type":       validation.Validate(req.Data.Type, validation.Required, validation.In("update_roles_ranks")),
		"data/attributes": validation.Validate(req.Data.Attributes, validation.Required),
	}

	if req.Data.Id.String() != chi.URLParam(r, "agglomeration_id") {
		errs["data/id"] = validation.NewError("mismatch", "query agglomeration_id and body data/id do not match")
	}

	return req, errs.Filter()
}
