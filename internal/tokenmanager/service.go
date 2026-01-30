package tokenmanager

import (
	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/bucket"
	"github.com/netbill/restkit/tokens"
)

type Manager struct {
	issuer   string
	uploadSK string
}

const (
	UploadOrgResource = "organization"
)

func New(issuer, uploadSK string) Manager {
	return Manager{
		issuer:   issuer,
		uploadSK: uploadSK,
	}
}

func (m Manager) NewUploadOrganizationMediaToken(
	OwnerAccountID uuid.UUID,
	ResourceID uuid.UUID,
	UploadSessionID uuid.UUID,
) (string, error) {
	return tokens.NewUploadFileToken(
		tokens.GenerateUploadFilesJwtRequest{
			OwnerAccountID:  OwnerAccountID,
			UploadSessionID: UploadSessionID,
			ResourceID:      ResourceID.String(),
			Resource:        UploadOrgResource,
			Issuer:          m.issuer,
			Audience:        []string{m.issuer},
			Ttl:             bucket.OrganizationUploadTTL,
		}, m.uploadSK)
}
