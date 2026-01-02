package controller

import (
	"net/http"

	"github.com/umisto/agglomerations-svc/internal/rest/responses"
	"github.com/umisto/ape"
	"github.com/umisto/ape/problems"
)

func (s Service) GetAllPermissions(w http.ResponseWriter, r *http.Request) {
	perms, err := s.core.GetAllPermissions(r.Context())
	if err != nil {
		s.log.WithError(err).Errorf("failed to get all permissions")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, http.StatusOK, responses.RolePermissions(perms))
}
