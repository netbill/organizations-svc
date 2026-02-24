package profile

import (
	"context"

	"github.com/google/uuid"
)

// TODO - delete profile by account id, and cascade delete all members of the profile
func (m *Module) Delete(
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
