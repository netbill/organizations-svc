package profile

import (
	"context"

	"github.com/google/uuid"
)

func (s Service) DeleteProfile(
	ctx context.Context,
	accountID uuid.UUID,
) error {
	return s.repo.Transaction(ctx, func(ctx context.Context) error {
		err := s.repo.DeleteProfileByAccountID(ctx, accountID)
		if err != nil {
			return err
		}

		err = s.repo.DeleteMembersByAccountID(ctx, accountID)
		if err != nil {
			return err
		}

		return nil
	})
}
