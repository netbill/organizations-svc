package tokenmanager

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/netbill/restkit/tokens"
)

type Manager struct {
	uploadSK string

	organizationMediaTokenTTL time.Duration
}

const (
	OrganizationsActor = "organizations-svc"
	OrgResource        = "organization"
)

func New(uploadSK string, organizationMediaTokenTTL time.Duration) *Manager {
	return &Manager{
		uploadSK:                  uploadSK,
		organizationMediaTokenTTL: organizationMediaTokenTTL,
	}
}

func (m *Manager) NewUploadOrganizationMediaToken(
	OwnerAccountID uuid.UUID,
	ResourceID uuid.UUID,
	UploadSessionID uuid.UUID,
) (string, error) {
	tkn, err := tokens.UploadContentClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   OwnerAccountID.String(),
			Issuer:    OrganizationsActor,
			Audience:  []string{OrganizationsActor},
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(m.organizationMediaTokenTTL)),
		},
		UploadSessionID: UploadSessionID,
		ResourceID:      ResourceID.String(),
		Resource:        OrgResource,
	}.GenerateJWT(m.uploadSK)
	if err != nil {
		return "", fmt.Errorf("failed to generate upload organization media token, cause: %w", err)
	}

	return tkn, nil
}
