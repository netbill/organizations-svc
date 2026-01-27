package profile

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
)

func (s Service) DeleteProfile(
	ctx context.Context,
	accountID uuid.UUID,
) error {
	return s.repo.Transaction(ctx, func(ctx context.Context) error {
		err := s.repo.DeleteProfileByAccountID(ctx, accountID)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to delete profile: %w", err),
			)
		}

		err = s.repo.DeleteMembersByAccountID(ctx, accountID)
		if err != nil {
			return errx.ErrorInternal.Raise(
				fmt.Errorf("failed to delete memberships for accountID %s: %w", accountID, err),
			)
		}

		return nil
	})
}
