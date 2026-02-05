package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/errx"

	"github.com/netbill/organizations-svc/internal/core/models"
	"github.com/netbill/organizations-svc/internal/core/modules/profile"
)

type ProfileRow struct {
	AccountID uuid.UUID `db:"account_id"`
	Username  string    `db:"username"`
	Official  bool      `db:"official"`

	Pseudonym *string `db:"pseudonym,omitempty"`
	Avatar    *string `db:"avatar,omitempty"`

	SourceCreatedAt  time.Time `db:"source_created_at"`
	SourceUpdatedAt  time.Time `db:"source_updated_at"`
	ReplicaCreatedAt time.Time `db:"replica_created_at"`
	ReplicaUpdatedAt time.Time `db:"replica_updated_at"`
}

func (r ProfileRow) IsNil() bool {
	return r.AccountID == uuid.Nil
}

func (r ProfileRow) ToModel() models.Profile {
	return models.Profile{
		AccountID: r.AccountID,
		Username:  r.Username,
		Official:  r.Official,
		Pseudonym: r.Pseudonym,
		Avatar:    r.Avatar,
		CreatedAt: r.SourceCreatedAt,
		UpdatedAt: r.SourceUpdatedAt,
	}
}

type ProfilesQ interface {
	New() ProfilesQ
	Insert(ctx context.Context, input ProfileRow) (ProfileRow, error)

	Get(ctx context.Context) (ProfileRow, error)
	Select(ctx context.Context) ([]ProfileRow, error)

	UpdateMany(ctx context.Context) (int64, error)
	UpdateOne(ctx context.Context) (ProfileRow, error)

	UpdateUsername(username string) ProfilesQ
	UpdateOfficial(official bool) ProfilesQ
	UpdatePseudonym(v *string) ProfilesQ
	UpdateAvatar(avatar *string) ProfilesQ
	UpdateSourceUpdatedAt(v time.Time) ProfilesQ

	FilterByAccountID(accountID ...uuid.UUID) ProfilesQ
	FilterByUsername(username string) ProfilesQ

	Delete(ctx context.Context) error
}

func (r *Repository) CreateProfile(ctx context.Context, profile models.Profile) (models.Profile, error) {
	row, err := r.profilesQ().Insert(ctx, ProfileRow{
		AccountID:       profile.AccountID,
		Username:        profile.Username,
		Official:        profile.Official,
		Pseudonym:       profile.Pseudonym,
		Avatar:          profile.Avatar,
		SourceUpdatedAt: profile.UpdatedAt,
		SourceCreatedAt: profile.CreatedAt,
	})
	if err != nil {
		return models.Profile{}, fmt.Errorf("failed to create profile, cause: %w", err)
	}

	return row.ToModel(), nil
}

func (r *Repository) UpdateProfile(ctx context.Context, accountID uuid.UUID, params profile.UpdateParams) (models.Profile, error) {
	row, err := r.profilesQ().
		FilterByAccountID(accountID).
		UpdateUsername(params.Username).
		UpdateOfficial(params.Official).
		UpdatePseudonym(params.Pseudonym).
		UpdateAvatar(params.Avatar).
		UpdateSourceUpdatedAt(params.UpdatedAt).
		UpdateOne(ctx)
	if err != nil {
		return models.Profile{}, fmt.Errorf("failed to update profile, cause: %w", err)
	}

	return row.ToModel(), nil
}

func (r *Repository) GetProfileByAccountID(ctx context.Context, accountID uuid.UUID) (models.Profile, error) {
	row, err := r.profilesQ().FilterByAccountID(accountID).Get(ctx)
	if err != nil {
		return models.Profile{}, fmt.Errorf("failed to get profile by account ID, cause: %w", err)
	}
	if row.IsNil() {
		return models.Profile{}, errx.ErrorProfileNotFound.Raise(
			fmt.Errorf("profile with account ID %s not found", accountID),
		)
	}

	return row.ToModel(), nil
}

func (r *Repository) GetProfileByUsername(ctx context.Context, username string) (models.Profile, error) {
	row, err := r.profilesQ().FilterByUsername(username).Get(ctx)
	if err != nil {
		return models.Profile{}, fmt.Errorf("failed to get profile by username, cause: %w", err)
	}
	if row.IsNil() {
		return models.Profile{}, errx.ErrorProfileNotFound.Raise(
			fmt.Errorf("profile with username %s not found", username),
		)
	}

	return row.ToModel(), nil
}

func (r *Repository) DeleteProfileByAccountID(ctx context.Context, accountID uuid.UUID) error {
	err := r.profilesQ().FilterByAccountID(accountID).Delete(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete profile by account ID, cause: %w", err)
	}

	return nil
}
