package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/models"
	"github.com/netbill/organizations-svc/internal/repository/pgdb"
)

func (s Service) UpsertProfile(ctx context.Context, profile models.Profile) (models.Profile, error) {
	row, err := s.profilesQ().Upsert(ctx, pgdb.ProfileUpsertInput{
		AccountID: profile.AccountID,
		Username:  profile.Username,
		Official:  profile.Official,
		Pseudonym: profile.Pseudonym,
	})
	if err != nil {
		return models.Profile{}, err
	}

	return Profile(row), nil
}

func (s Service) UpdateUsername(ctx context.Context, accountID uuid.UUID, username string) (models.Profile, error) {
	row, err := s.profilesQ().
		FilterByAccountID(accountID).
		UpdateUsername(username).
		Get(ctx)
	if err != nil {
		return models.Profile{}, err
	}

	return Profile(row), nil
}

func (s Service) GetProfileByAccountID(ctx context.Context, accountID uuid.UUID) (models.Profile, error) {
	row, err := s.profilesQ().FilterByAccountID(accountID).Get(ctx)
	if err != nil {
		return models.Profile{}, err
	}

	return Profile(row), nil
}

func (s Service) GetProfileByUsername(ctx context.Context, username string) (models.Profile, error) {
	row, err := s.profilesQ().FilterByUsername(username).Get(ctx)
	if err != nil {
		return models.Profile{}, err
	}

	return Profile(row), nil
}

func (s Service) DeleteProfileByAccountID(ctx context.Context, accountID uuid.UUID) error {
	return s.profilesQ().FilterByAccountID(accountID).Delete(ctx)
}

func Profile(row pgdb.Profile) models.Profile {
	return models.Profile{
		AccountID: row.AccountID,
		Username:  row.Username,
		Official:  row.Official,
		Pseudonym: row.Pseudonym,
	}
}
