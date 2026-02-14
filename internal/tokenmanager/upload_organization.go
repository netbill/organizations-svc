package tokenmanager

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/netbill/restkit/tokens"
)

const (
	OrgResource = "organization"
)

func (m *Manager) GenerateUploadOrganizationMediaToken(
	OwnerAccountID uuid.UUID,
	OrganizationID uuid.UUID,
	UploadSessionID uuid.UUID,
) (string, error) {
	tkn, err := tokens.UploadContentClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   OwnerAccountID.String(),
			Issuer:    m.issuer,
			Audience:  []string{m.issuer},
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(m.mediaTTL)),
		},
		UploadSessionID: UploadSessionID,
		ResourceID:      OrganizationID.String(),
		ResourceType:    OrgResource,
	}.GenerateJWT(m.uploadSK)
	if err != nil {
		return "", fmt.Errorf("failed to generate upload organization media token, cause: %w", err)
	}

	return tkn, nil
}

func (m *Manager) ParseUploadOrganizationContentToken(token string) (tokens.UploadContentClaims, error) {
	res, err := tokens.ParseUploadFilesClaims(token, m.uploadSK)
	if err != nil {
		return tokens.UploadContentClaims{}, fmt.Errorf(
			"failed to validate upload organization media token, cause: %w", err,
		)
	}

	if res.ResourceType != OrgResource {
		return tokens.UploadContentClaims{}, fmt.Errorf("invalid upload token resource type")
	}

	audSuccess := false
	for _, aud := range res.Audience {
		if aud == m.issuer {
			audSuccess = true
			break
		}
	}
	if !audSuccess {
		return tokens.UploadContentClaims{}, fmt.Errorf("invalid upload token audience")
	}

	_, err = uuid.Parse(res.Subject)
	if err != nil {
		return tokens.UploadContentClaims{}, fmt.Errorf("invalid upload token subject: %w", err)
	}

	_, err = uuid.Parse(res.ResourceID)
	if err != nil {
		return tokens.UploadContentClaims{}, fmt.Errorf("invalid upload token resource ID: %w", err)
	}

	return res, nil
}
