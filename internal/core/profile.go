package core

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/models"
)

type profileRepository interface {
	CreateProfile(ctx context.Context, profile models.Profile) (models.Profile, error)

	GetProfileByAccountID(ctx context.Context, accountID uuid.UUID) (models.Profile, error)
	GetProfilesByAccountIDs(ctx context.Context, accountIDs []uuid.UUID) ([]models.Profile, error)
	GetProfileByUsername(ctx context.Context, username string) (models.Profile, error)
	ExistsProfileByAccountID(ctx context.Context, accountID uuid.UUID) (bool, error)
	ExistsProfileByUsername(ctx context.Context, username string) (bool, error)

	UpdateProfile(ctx context.Context, accountID uuid.UUID, params ProfileUpdateParams) (models.Profile, error)

	DeleteProfileByAccountID(ctx context.Context, accountID uuid.UUID) error
}

type tombstoneProfileRepository interface {
	BuryProfile(ctx context.Context, accountID uuid.UUID) error
	ProfileIsBuried(ctx context.Context, accountID uuid.UUID) (bool, error)
}

type ProfileModule struct {
	profileRepo   profileRepository
	tombstoneRepo tombstoneProfileRepository
	tx            transactor
}

func (m *ProfileModule) Create(
	ctx context.Context,
	profile models.Profile,
) (models.Profile, error) {
	exists, err := m.profileRepo.ExistsProfileByAccountID(ctx, profile.AccountID)
	if err != nil {
		return models.Profile{}, err
	}
	if exists {
		return models.Profile{}, errx.ErrorProfileAlreadyExists.Raise(
			errors.New("profile with given account id already exists"),
		)
	}

	buried, err := m.tombstoneRepo.ProfileIsBuried(ctx, profile.AccountID)
	if err != nil {
		return models.Profile{}, err
	}
	if buried {
		return models.Profile{}, errx.ErrorProfileDeleted.Raise(
			errors.New("profile with given account id is deleted"),
		)
	}

	return m.profileRepo.CreateProfile(ctx, profile)
}

func (m *ProfileModule) GetByID(
	ctx context.Context,
	accountID uuid.UUID,
) (models.Profile, error) {
	return m.profileRepo.GetProfileByAccountID(ctx, accountID)
}

func (m *ProfileModule) GetByIDs(
	ctx context.Context,
	accountIDs []uuid.UUID,
) ([]models.Profile, error) {
	return m.profileRepo.GetProfilesByAccountIDs(ctx, accountIDs)
}

type ProfileUpdateParams struct {
	Username  string    `json:"username"`
	Pseudonym *string   `json:"pseudonym"`
	AvatarKey *string   `json:"avatar"`
	Version   int32     `json:"version"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (m *ProfileModule) Update(
	ctx context.Context,
	accountID uuid.UUID,
	params ProfileUpdateParams,
) (models.Profile, error) {
	profile, err := m.profileRepo.GetProfileByAccountID(ctx, accountID)
	if errors.Is(err, errx.ErrorProfileNotExists) {
		buried, err := m.tombstoneRepo.ProfileIsBuried(ctx, accountID)
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

	return m.profileRepo.UpdateProfile(ctx, accountID, params)
}

func (m *ProfileModule) Delete(
	ctx context.Context,
	accountID uuid.UUID,
) error {
	buried, err := m.tombstoneRepo.ProfileIsBuried(ctx, accountID)
	if err != nil {
		return err
	}
	if buried {
		return errx.ErrorProfileDeleted.Raise(
			fmt.Errorf("profile with account id %s is already deleted", accountID),
		)
	}

	return m.tx.Transaction(ctx, func(ctx context.Context) error {
		if err := m.tombstoneRepo.BuryProfile(ctx, accountID); err != nil {
			return err
		}

		return m.profileRepo.DeleteProfileByAccountID(ctx, accountID)
	})
}
