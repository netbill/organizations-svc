package profile

import (
	"context"

	"github.com/google/uuid"
)

func (m *Module) DeleteProfile(
	ctx context.Context,
	accountID uuid.UUID,
) error {
	return m.repo.Transaction(ctx, func(ctx context.Context) error {
		err := m.repo.DeleteProfileByAccountID(ctx, accountID)
		if err != nil {
			return err
		}

		err = m.repo.DeleteMembersByAccountID(ctx, accountID)
		if err != nil {
			return err
		}

		return nil
	})
}
