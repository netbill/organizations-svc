package responses

import (
	"github.com/netbill/organizations-svc/internal/core/models"
	"github.com/netbill/organizations-svc/pkg/resources"
)

func profileData(profile models.Profile) resources.ProfileData {
	return resources.ProfileData{
		Id:   profile.AccountID,
		Type: "profile",
		Attributes: resources.ProfileAttributes{
			Username:  profile.Username,
			Pseudonym: profile.Pseudonym,
			Official:  profile.Official,
			AvatarKey: profile.AvatarKey,
			Version:   profile.Version,
			UpdatedAt: profile.UpdatedAt,
			CreatedAt: profile.CreatedAt,
		},
	}
}
