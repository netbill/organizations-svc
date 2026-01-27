package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/models"
	"github.com/netbill/organizations-svc/internal/core/modules/profile"
	"github.com/netbill/organizations-svc/internal/repository/pgdb"
)

func (r Repository) CreateProfile(ctx context.Context, profile models.Profile) (models.Profile, error) {
	row, err := r.profilesQ(ctx).Insert(ctx, pgdb.ProfileInsertInput{
		AccountID:       profile.AccountID,
		Username:        profile.Username,
		Official:        profile.Official,
		Pseudonym:       profile.Pseudonym,
		SourceUpdatedAt: profile.UpdatedAt,
		SourceCreatedAt: profile.CreatedAt,
	})
	if err != nil {
		return models.Profile{}, fmt.Errorf("failed to create profile, cause: %w", err)
	}

	return Profile(row), nil
}

func (r Repository) UpdateProfile(ctx context.Context, accountID uuid.UUID, params profile.UpdateParams) (models.Profile, error) {
	row, err := r.profilesQ(ctx).
		FilterByAccountID(accountID).
		UpdateUsername(params.Username).
		UpdateOfficial(params.Official).
		UpdatePseudonym(params.Pseudonym).
		UpdateSourceUpdatedAt(params.UpdatedAt).
		UpdateOne(ctx)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return models.Profile{}, errx.ErrorProfileNotFound.Raise(
			fmt.Errorf("profile with accountID %s not found", accountID),
		)
	case err != nil:
		return models.Profile{}, fmt.Errorf("failed to update profile with accountID %s, cause: %w", accountID, err)
	}

	return Profile(row), nil
}

func (r Repository) GetProfileByAccountID(ctx context.Context, accountID uuid.UUID) (models.Profile, error) {
	row, err := r.profilesQ(ctx).FilterByAccountID(accountID).Get(ctx)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return models.Profile{}, errx.ErrorProfileNotFound.Raise(
			fmt.Errorf("profile with accountID %s not found", accountID),
		)
	case err != nil:
		return models.Profile{}, fmt.Errorf("failed to get profile with accountID %s, cause: %w", accountID, err)
	}

	return Profile(row), nil
}

func (r Repository) GetProfileByUsername(ctx context.Context, username string) (models.Profile, error) {
	row, err := r.profilesQ(ctx).FilterByUsername(username).Get(ctx)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return models.Profile{}, errx.ErrorProfileNotFound.Raise(
			fmt.Errorf("profile with username %s not found", username),
		)
	case err != nil:
		return models.Profile{}, fmt.Errorf("failed to get profile with username %s, cause: %w", username, err)
	}

	return Profile(row), nil
}

func (r Repository) DeleteProfileByAccountID(ctx context.Context, accountID uuid.UUID) error {
	err := r.profilesQ(ctx).FilterByAccountID(accountID).Delete(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete profile with accountID %s, cause: %w", accountID, err)
	}

	return nil
}

func Profile(row pgdb.Profile) models.Profile {
	res := models.Profile{
		AccountID: row.AccountID,
		Username:  row.Username,
		Official:  row.Official,
		CreatedAt: row.SourceCreatedAt,
		UpdatedAt: row.SourceUpdatedAt,
	}
	if row.Pseudonym.Valid {
		res.Pseudonym = &row.Pseudonym.String
	}
	if row.Avatar.Valid {
		res.Avatar = &row.Avatar.String
	}

	return res
}
