package controller

import (
	"fmt"
	"net/http"

	"github.com/netbill/organizations-svc/internal/core/models"
	"github.com/netbill/organizations-svc/internal/rest/contexter"
	"github.com/netbill/organizations-svc/internal/rest/responses"
	"github.com/netbill/restkit/pagi"
	"github.com/netbill/restkit/problems"
)

func (c *Controller) GetMyOrganizations(w http.ResponseWriter, r *http.Request) {
	initiator, err := contexter.AccountData(r.Context())
	if err != nil {
		c.log.WithError(err).Errorf("failed to get initiator account data")
		c.responser.RenderErr(w, problems.Unauthorized("failed to get initiator account data"))
		return
	}

	limit, offset := pagi.GetPagination(r)
	if limit > 100 {
		c.log.WithError(fmt.Errorf("invalid pagination limit %d", limit)).Errorf("invalid pagination limit")
		c.responser.RenderErr(w, problems.BadRequest(fmt.Errorf("pagination limit must be between 1 and 100"))...)
		return
	}

	res, err := c.core.GetOrganizationForUser(
		r.Context(),
		models.InitiatorData{
			AccountID: initiator.GetAccountID(),
		},
		limit, offset,
	)
	if err != nil {
		c.log.WithError(err).Errorf("failed to get organizations")
		c.responser.RenderErr(w, problems.InternalError())
		return
	}

	c.responser.Render(w, http.StatusOK, responses.Organizations(r, res))
}
