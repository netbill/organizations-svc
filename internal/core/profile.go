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

type profileRepo interface {
	Create(ctx context.Context, profile models.Profile) (models.Profile, error)

	GetByID(ctx context.Context, accountID uuid.UUID) (models.Profile, error)
	GetByUsername(ctx context.Context, username string) (models.Profile, error)
	ExistsByUsername(ctx context.Context, username string) (bool, error)

	GetListByIDs(ctx context.Context, accountIDs []uuid.UUID) ([]models.Profile, error)
	ExistsByAccountID(ctx context.Context, accountID uuid.UUID) (bool, error)

	Update(ctx context.Context, accountID uuid.UUID, params ProfileUpdateParams) (models.Profile, error)
	Delete(ctx context.Context, accountID uuid.UUID) error
}

type profileTombstone interface {
	BuryProfile(ctx context.Context, accountID uuid.UUID) error
	ProfileIsBuried(ctx context.Context, accountID uuid.UUID) (bool, error)
}

type ProfileModule struct {
	repo      profileRepo
	tombstone profileTombstone
	tx        transactor
}

type ProfileDeps struct {
	Repo      profileRepo
	Tombstone profileTombstone
	Tx        transactor
}

func NewProfileModule(deps ProfileDeps) *ProfileModule {
	return &ProfileModule{
		repo:      deps.Repo,
		tombstone: deps.Tombstone,
		tx:        deps.Tx,
	}
}

func (m *ProfileModule) Create(
	ctx context.Context,
	profile models.Profile,
) (models.Profile, error) {
	exists, err := m.repo.ExistsByAccountID(ctx, profile.AccountID)
	if err != nil {
		return models.Profile{}, err
	}
	if exists {
		return models.Profile{}, errx.ErrorProfileAlreadyExists.Raise(
			errors.New("profile with given account id already exists"),
		)
	}

	buried, err := m.tombstone.ProfileIsBuried(ctx, profile.AccountID)
	if err != nil {
		return models.Profile{}, err
	}
	if buried {
		return models.Profile{}, errx.ErrorProfileDeleted.Raise(
			errors.New("profile with given account id is deleted"),
		)
	}

	return m.repo.Create(ctx, profile)
}

func (m *ProfileModule) GetByID(
	ctx context.Context,
	accountID uuid.UUID,
) (models.Profile, error) {
	return m.repo.GetByID(ctx, accountID)
}

func (m *ProfileModule) GetByIDs(
	ctx context.Context,
	accountIDs []uuid.UUID,
) ([]models.Profile, error) {
	return m.repo.GetListByIDs(ctx, accountIDs)
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
	profile, err := m.repo.GetByID(ctx, accountID)
	if errors.Is(err, errx.ErrorProfileNotExists) {
		buried, err := m.tombstone.ProfileIsBuried(ctx, accountID)
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

	return m.repo.Update(ctx, accountID, params)
}

func (m *ProfileModule) Delete(
	ctx context.Context,
	accountID uuid.UUID,
) error {
	buried, err := m.tombstone.ProfileIsBuried(ctx, accountID)
	if err != nil {
		return err
	}
	if buried {
		return errx.ErrorProfileDeleted.Raise(
			fmt.Errorf("profile with account id %s is already deleted", accountID),
		)
	}

	return m.tx.Transaction(ctx, func(ctx context.Context) error {
		if err := m.tombstone.BuryProfile(ctx, accountID); err != nil {
			return err
		}

		return m.repo.Delete(ctx, accountID)
	})
}
