package invite

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/models"
	"github.com/netbill/restkit/pagi"
)

func (m *Module) GetForOrganizations(
	ctx context.Context,
	initiator models.Initiator,
	organizationID uuid.UUID,
	limit, offset uint,
) (pagi.Page[[]models.Invite], error) {
	_, err := m.getInitiator(ctx, initiator, organizationID)
	if err != nil {
		return pagi.Page[[]models.Invite]{}, err
	}

	res, err := m.repo.GetOrganizationInvites(ctx, organizationID, limit, offset)
	if err != nil {
		return pagi.Page[[]models.Invite]{}, err
	}

	return res, nil
}
