package request

import (
	"encoding/json"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/netbill/organizations-svc/resources"
)

func CreateRole(r *http.Request) (req resources.Role, err error) {
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		return req, err
	}

	errs := validation.Errors{
		"data/type":       validation.Validate(req.Data.Type, validation.Required, validation.In("create_role")),
		"data/attributes": validation.Validate(req.Data.Attributes, validation.Required),
	}

	return req, errs.Filter()
}
