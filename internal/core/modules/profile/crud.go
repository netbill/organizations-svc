package profile

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/models"
)

func (m *Module) Create(
	ctx context.Context,
	profile models.Profile,
) (models.Profile, error) {
	exists, err := m.repo.ExistsProfileByAccountID(ctx, profile.AccountID)
	if err != nil {
		return models.Profile{}, err
	}
	if exists {
		return models.Profile{}, errx.ErrorProfileAlreadyExists.Raise(
			errors.New("profile with given account id already exists"),
		)
	}

	buried, err := m.repo.ProfileIsBuried(ctx, profile.AccountID)
	if err != nil {
		return models.Profile{}, err
	}
	if buried {
		return models.Profile{}, errx.ErrorProfileDeleted.Raise(
			errors.New("profile with given account id is deleted"),
		)
	}

	return m.repo.CreateProfile(ctx, profile)
}

func (m *Module) GetByID(
	ctx context.Context,
	accountID uuid.UUID,
) (models.Profile, error) {
	return m.repo.GetProfileByAccountID(ctx, accountID)
}

func (m *Module) GetByIDs(
	ctx context.Context,
	accountIDs []uuid.UUID,
) ([]models.Profile, error) {
	return m.repo.GetProfilesByAccountIDs(ctx, accountIDs)
}

type UpdateParams struct {
	Username  string    `json:"username"`
	Pseudonym *string   `json:"pseudonym"`
	AvatarKey *string   `json:"avatar"`
	Version   int32     `json:"version"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (m *Module) Update(
	ctx context.Context,
	accountID uuid.UUID,
	params UpdateParams,
) (models.Profile, error) {
	profile, err := m.repo.GetProfileByAccountID(ctx, accountID)
	if errors.Is(err, errx.ErrorProfileNotExists) {
		buried, err := m.repo.ProfileIsBuried(ctx, accountID)
		if err != nil {
			return models.Profile{}, err
		}
		if buried {
			return models.Profile{}, errx.ErrorProfileDeleted.Raise(
				fmt.Errorf("profile with account id %s is already deleted", accountID),
			)
		}
	}
	if err != nil {
		return models.Profile{}, err
	}
	if profile.Version >= params.Version {
		return profile, nil
	}

	updatedProfile, err := m.repo.UpdateProfile(ctx, accountID, params)
	if err != nil {
		return models.Profile{}, err
	}

	return updatedProfile, nil
}

func (m *Module) Delete(
	ctx context.Context,
	accountID uuid.UUID,
) error {
	buried, err := m.repo.ProfileIsBuried(ctx, accountID)
	if err != nil {
		return err
	}
	if buried {
		return errx.ErrorProfileDeleted.Raise(
			fmt.Errorf("profile with account id %s is already deleted", accountID),
		)
	}

	return m.repo.Transaction(ctx, func(ctx context.Context) error {
		if err := m.repo.BuryProfile(ctx, accountID); err != nil {
			return err
		}

		return m.repo.DeleteProfileByAccountID(ctx, accountID)
	})
}
