package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/umisto/cities-svc/internal/domain/entity"
	"github.com/umisto/cities-svc/internal/repository/pgdb"
)

func (s Service) CreateProfile(ctx context.Context, profile entity.Profile) (entity.Profile, error) {
	row, err := s.profilesQ().Upsert(ctx, pgdb.ProfileUpsertInput{
		AccountID: profile.AccountID,
		Username:  profile.Username,
		Official:  profile.Official,
		Pseudonym: profile.Pseudonym,
	})
	if err != nil {
		return entity.Profile{}, err
	}

	return Profile(row), nil
}

func (s Service) GetProfileByAccountID(ctx context.Context, accountID uuid.UUID) (entity.Profile, error) {
	row, err := s.profilesQ().FilterByAccountID(accountID).Get(ctx)
	if err != nil {
		return entity.Profile{}, err
	}

	return Profile(row), nil
}

func (s Service) GetProfileByUsername(ctx context.Context, username string) (entity.Profile, error) {
	row, err := s.profilesQ().FilterByUsername(username).Get(ctx)
	if err != nil {
		return entity.Profile{}, err
	}

	return Profile(row), nil
}

func (s Service) DeleteProfileByAccountID(ctx context.Context, accountID uuid.UUID) error {
	return s.profilesQ().FilterByAccountID(accountID).Delete(ctx)
}

func Profile(row pgdb.Profile) entity.Profile {
	return entity.Profile{
		AccountID: row.AccountID,
		Username:  row.Username,
		Official:  row.Official,
		Pseudonym: row.Pseudonym,
	}
}
