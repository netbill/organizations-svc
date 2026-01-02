package request

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/netbill/organizations-svc/resources"
)

func UpdateMember(r *http.Request) (params resources.UpdateMember, err error) {
	if err = json.NewDecoder(r.Body).Decode(&params); err != nil {
		err = newDecodeError("body", err)
		return
	}

	errs := validation.Errors{
		"data/id":         validation.Validate(params.Data.Id, validation.Required),
		"data/type":       validation.Validate(params.Data.Type, validation.Required, validation.In("update_member")),
		"data/attributes": validation.Validate(params.Data.Attributes, validation.Required),
	}

	if chi.URLParam(r, "member_id") != params.Data.Id.String() {
		errs["data/id"] = validation.NewError("mismatch", "query member_id and body data/id do not match")
	}

	return params, errs.Filter()
}
