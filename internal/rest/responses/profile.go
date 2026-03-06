package responses

import (
	"net/http"

	"github.com/netbill/organizations-svc/internal/models"
	"github.com/netbill/organizations-svc/internal/rest/scope"
	"github.com/netbill/organizations-svc/pkg/resources"
)

func profileData(r *http.Request, profile models.Profile) resources.ProfileData {
	res := resources.ProfileData{
		Id:   profile.AccountID,
		Type: "profile",
		Attributes: resources.ProfileAttributes{
			Username:  profile.Username,
			Pseudonym: profile.Pseudonym,
			Version:   profile.Version,
			UpdatedAt: profile.UpdatedAt,
			CreatedAt: profile.CreatedAt,
		},
	}
	if profile.AvatarKey != nil {
		url := scope.ResolverURL(r, *profile.AvatarKey)
		res.Attributes.AvatarUrl = &url
	}

	return res
}
