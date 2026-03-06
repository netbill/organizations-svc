package profile

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/errx"
	"github.com/netbill/organizations-svc/internal/models"
)

type Service struct {
	repo      repo
	tombstone tombstone
	tx        transactor
}

type ServiceDeps struct {
	Repo      repo
	Tombstone tombstone
	Tx        transactor
}

func NewProfileModule(deps ServiceDeps) *Service {
	return &Service{
		repo:      deps.Repo,
		tombstone: deps.Tombstone,
		tx:        deps.Tx,
	}
}

type CreateParams struct {
	AccountID uuid.UUID `json:"account_id"`
	Username  string    `json:"username"`
	Pseudonym *string   `json:"pseudonym,omitempty"`
	AvatarKey *string   `json:"avatar_key,omitempty"`

	Version   int32     `json:"version"`
	CreatedAt time.Time `json:"created_at"`
}

func (s *Service) Create(
	ctx context.Context,
	profile CreateParams,
) (models.Profile, error) {
	exists, err := s.repo.ExistsByAccountID(ctx, profile.AccountID)
	if err != nil {
		return models.Profile{}, err
	}
	if exists {
		return models.Profile{}, errx.ErrorProfileAlreadyExists.Raise(
			errors.New("profile with given account id already exists"),
		)
	}

	return s.repo.Create(ctx, profile)
}

func (s *Service) GetByID(
	ctx context.Context,
	accountID uuid.UUID,
) (models.Profile, error) {
	buried, err := s.tombstone.ProfileIsBuried(ctx, accountID)
	if err != nil {
		return models.Profile{}, err
	}
	if buried {
		return models.Profile{}, errx.ErrorProfileDeleted.Raise(
			errors.New("profile with given account id is deleted"),
		)
	}

	return s.repo.GetByID(ctx, accountID)
}

func (s *Service) GetByIDs(
	ctx context.Context,
	accountIDs []uuid.UUID,
) ([]models.Profile, error) {
	return s.repo.GetListByIDs(ctx, accountIDs)
}

type UpdateParams struct {
	Username  string    `json:"username"`
	Pseudonym *string   `json:"pseudonym"`
	AvatarKey *string   `json:"avatar"`
	Version   int32     `json:"version"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (s *Service) Update(
	ctx context.Context,
	accountID uuid.UUID,
	params UpdateParams,
) (models.Profile, error) {
	profile, err := s.GetByID(ctx, accountID)
	if err != nil {
		return models.Profile{}, err
	}
	if profile.Version >= params.Version {
		return profile, nil
	}

	return s.repo.Update(ctx, accountID, params)
}

func (s *Service) Delete(
	ctx context.Context,
	accountID uuid.UUID,
) error {
	buried, err := s.tombstone.ProfileIsBuried(ctx, accountID)
	if err != nil {
		return err
	}
	if buried {
		return errx.ErrorProfileDeleted.Raise(
			fmt.Errorf("profile with account id %s is already deleted", accountID),
		)
	}

	return s.tx.Transaction(ctx, func(ctx context.Context) error {
		if err := s.tombstone.BuryProfile(ctx, accountID); err != nil {
			return err
		}

		return s.repo.Delete(ctx, accountID)
	})
}
