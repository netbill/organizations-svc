package request

import (
	"encoding/json"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/netbill/organizations-svc/resources"
)

func SentInvite(r *http.Request) (params resources.CreateInvite, err error) {
	if err = json.NewDecoder(r.Body).Decode(&params); err != nil {
		err = newDecodeError("body", err)
		return
	}

	errs := validation.Errors{
		"data/type":       validation.Validate(params.Data.Type, validation.Required, validation.In("create_invite")),
		"data/attributes": validation.Validate(params.Data.Attributes, validation.Required),
	}

	return params, errs.Filter()
}
