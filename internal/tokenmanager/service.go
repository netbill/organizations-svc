package tokenmanager

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/netbill/restkit/tokens"
)

type Manager struct {
	uploadSK string

	organizationMediaTokenTTL time.Duration
}

const (
	OrganizationsActor = "organizations-svc"
	UploadOrgResource  = "organization"
)

func New(uploadSK string, organizationMediaTokenTTL time.Duration) Manager {
	return Manager{
		uploadSK:                  uploadSK,
		organizationMediaTokenTTL: organizationMediaTokenTTL,
	}
}

func (m Manager) NewUploadOrganizationMediaToken(
	OwnerAccountID uuid.UUID,
	ResourceID uuid.UUID,
	UploadSessionID uuid.UUID,
) (string, error) {
	tkn, err := tokens.NewUploadFileToken(
		tokens.GenerateUploadFilesJwtRequest{
			OwnerAccountID:  OwnerAccountID,
			UploadSessionID: UploadSessionID,
			ResourceID:      ResourceID.String(),
			Resource:        UploadOrgResource,
			Issuer:          OrganizationsActor,
			Audience:        []string{OrganizationsActor},
			Ttl:             m.organizationMediaTokenTTL,
		}, m.uploadSK,
	)
	if err != nil {
		return "", fmt.Errorf("failed to generate upload organization media token, cause: %w", err)
	}

	return tkn, nil
}
