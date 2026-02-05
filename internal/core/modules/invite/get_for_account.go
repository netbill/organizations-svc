package invite

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/models"
)

func (m *Module) GetForAccount(
	ctx context.Context,
	accountID uuid.UUID,
	inviteID uuid.UUID,
) (models.Invite, error) {
	res, err := m.repo.GetInvite(ctx, inviteID)
	if err != nil {
		return models.Invite{}, err
	}

	if res.AccountID != accountID {
		_, err = m.repo.GetMemberByAccountAndOrganization(
			ctx,
			accountID,
			res.OrganizationID,
		)
		if err != nil {
			if errors.Is(err, errx.ErrorMemberNotFound) {
				return models.Invite{}, errx.ErrorInviteNotFound.Raise(
					fmt.Errorf("account has no rights to view this invite"),
				)
			}
			return models.Invite{}, err
		}
	}

	return res, nil
}
