package requests

import (
	"encoding/json"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/netbill/organizations-svc/pkg/resources"
	"github.com/netbill/restkit"
)

func CreateInvite(r *http.Request) (params resources.CreateInvite, err error) {
	if err = json.NewDecoder(r.Body).Decode(&params); err != nil {
		err = restkit.NewDecodeError("body", err)
		return
	}

	errs := validation.Errors{
		"data/type":       validation.Validate(params.Data.Type, validation.Required, validation.In("invite")),
		"data/attributes": validation.Validate(params.Data.Attributes, validation.Required),
	}

	return params, errs.Filter()
}
