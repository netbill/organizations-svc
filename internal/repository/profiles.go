package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/models"
	"github.com/netbill/organizations-svc/internal/core/modules/profile"
	"github.com/netbill/organizations-svc/internal/repository/pgdb"
)

func (s Service) CreateProfile(ctx context.Context, profile models.Profile) (models.Profile, error) {
	row, err := s.profilesQ(ctx).Insert(ctx, pgdb.ProfileInsertInput{
		AccountID:       profile.AccountID,
		Username:        profile.Username,
		Official:        profile.Official,
		Pseudonym:       profile.Pseudonym,
		SourceUpdatedAt: profile.UpdatedAt,
		SourceCreatedAt: profile.CreatedAt,
	})
	if err != nil {
		return models.Profile{}, err
	}

	return Profile(row), nil
}

func (s Service) UpdateProfile(ctx context.Context, accountID uuid.UUID, params profile.UpdateParams) (models.Profile, error) {
	row, err := s.profilesQ(ctx).
		FilterByAccountID(accountID).
		UpdateUsername(params.Username).
		UpdateOfficial(params.Official).
		UpdatePseudonym(params.Pseudonym).
		UpdateSourceUpdatedAt(params.UpdatedAt).
		UpdateOne(ctx)
	if err != nil {
		return models.Profile{}, err
	}

	return Profile(row), nil
}

func (s Service) GetProfileByAccountID(ctx context.Context, accountID uuid.UUID) (models.Profile, error) {
	row, err := s.profilesQ(ctx).FilterByAccountID(accountID).Get(ctx)
	if err != nil {
		return models.Profile{}, err
	}

	return Profile(row), nil
}

func (s Service) GetProfileByUsername(ctx context.Context, username string) (models.Profile, error) {
	row, err := s.profilesQ(ctx).FilterByUsername(username).Get(ctx)
	if err != nil {
		return models.Profile{}, err
	}

	return Profile(row), nil
}

func (s Service) DeleteProfileByAccountID(ctx context.Context, accountID uuid.UUID) error {
	return s.profilesQ(ctx).FilterByAccountID(accountID).Delete(ctx)
}

func Profile(row pgdb.Profile) models.Profile {
	return models.Profile{
		AccountID: row.AccountID,
		Username:  row.Username,
		Official:  row.Official,
		Pseudonym: row.Pseudonym,
		CreatedAt: row.SourceCreatedAt,
		UpdatedAt: row.SourceUpdatedAt,
	}
}
