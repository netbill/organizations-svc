package tokenmanager

import (
	"github.com/google/uuid"
	"github.com/netbill/profiles-svc/internal/bucket"
	"github.com/netbill/restkit/tokens"
)

type Manager struct {
	issuer   string
	uploadSK string
}

const UploadProfileAvatarScope = "upload:profile_avatar"

func New(issuer, uploadSK string) Manager {
	return Manager{
		issuer:   issuer,
		uploadSK: uploadSK,
	}
}

func (m Manager) NewUploadProfileAvatarToken(
	sessionID uuid.UUID,
) (string, error) {
	return tokens.NewUploadFileToken(
		tokens.GenerateUploadFilesJwtRequest{
			SessionID: sessionID,
			Issuer:    m.issuer,
			Audience:  []string{m.issuer},
			Scope:     UploadProfileAvatarScope,
			Ttl:       bucket.ProfileAvatarUploadTTL,
		}, m.uploadSK)
}
