package invite

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/umisto/cities-svc/internal/domain/entity"
	"github.com/umisto/cities-svc/internal/domain/errx"
)

func (s Service) DeclineInvite(
	ctx context.Context,
	accountID, inviteID uuid.UUID,
) (entity.Invite, error) {
	invite, err := s.repo.GetInviteByID(ctx, inviteID)
	if err != nil {
		return entity.Invite{}, err
	}

	if invite.AccountID != accountID {
		return entity.Invite{}, errx.ErrorNotEnoughRights.Raise(
			fmt.Errorf("account has no rights to decline this invite"),
		)
	}
	if invite.Status != entity.InviteStatusSent {
		return entity.Invite{}, err
	}

	invite, err = s.repo.UpdateInviteStatus(ctx, inviteID, entity.InviteStatusDeclined)
	if err != nil {
		return entity.Invite{}, err
	}

	return invite, nil
}
