package request

import (
	"encoding/json"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/umisto/agglomerations-svc/resources"
)

func newDecodeError(what string, err error) error {
	return validation.Errors{
		what: fmt.Errorf("decode request %s: %w", what, err),
	}
}

func CreateAgglomeration(r *http.Request) (req resources.CreateAgglomeration, err error) {
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		err = newDecodeError("body", err)
		return
	}

	errs := validation.Errors{
		"data/type":       validation.Validate(req.Data.Type, validation.Required, validation.In("create_agglomeration")),
		"data/attributes": validation.Validate(req.Data.Attributes, validation.Required),
	}

	return req, errs.Filter()
}
