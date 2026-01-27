package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/models"
	"github.com/netbill/organizations-svc/internal/core/modules/organization"
	"github.com/netbill/organizations-svc/internal/repository/pgdb"
	"github.com/netbill/pagi"
	"github.com/pkg/errors"
)

func (r Repository) CreateOrganization(
	ctx context.Context,
	params organization.CreateParams,
) (models.Organization, error) {
	row, err := r.organizationsQ(ctx).Insert(ctx, pgdb.OrganizationsQInsertInput{
		Name: params.Name,
	})
	if err != nil {
		return models.Organization{}, fmt.Errorf(
			"failed to create organization, cause: %w", err,
		)
	}

	return Organization(row), nil
}

func (r Repository) UpdateOrganization(
	ctx context.Context,
	ID uuid.UUID,
	params organization.UpdateParams,
) (models.Organization, error) {
	q := r.organizationsQ(ctx).FilterByID(ID)
	if params.Name != nil {
		q = q.UpdateName(*params.Name)
	}

	row, err := q.UpdateOne(ctx)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return models.Organization{}, errx.ErrorOrganizationNotFound.Raise(
			fmt.Errorf("organization with ID %s not found", ID),
		)
	case err != nil:
		return models.Organization{}, fmt.Errorf(
			"failed to update organization with ID %s, cause: %w", ID, err,
		)
	}

	return Organization(row), nil
}

func (r Repository) UpdateOrganizationStatus(
	ctx context.Context,
	ID uuid.UUID,
	status string,
) (models.Organization, error) {
	row, err := r.organizationsQ(ctx).FilterByID(ID).UpdateStatus(status).UpdateOne(ctx)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return models.Organization{}, errx.ErrorOrganizationNotFound.Raise(
			fmt.Errorf("organization with ID %s not found", ID),
		)
	case err != nil:
		return models.Organization{}, fmt.Errorf(
			"failed to update organization status with ID %s, cause: %w", ID, err,
		)
	}

	return Organization(row), nil
}

func (r Repository) UpdateOrganizationMaxRoles(
	ctx context.Context,
	ID uuid.UUID,
	maxRoles uint,
) (models.Organization, error) {
	row, err := r.organizationsQ(ctx).FilterByID(ID).UpdateMaxRoles(maxRoles).UpdateOne(ctx)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return models.Organization{}, errx.ErrorOrganizationNotFound.Raise(
			fmt.Errorf("organization with ID %s not found", ID),
		)
	case err != nil:
		return models.Organization{}, fmt.Errorf(
			"failed to update organization max roles with ID %s, cause: %w", ID, err,
		)
	}

	return Organization(row), nil
}

func (r Repository) GetOrganizationByID(ctx context.Context, ID uuid.UUID) (models.Organization, error) {
	row, err := r.organizationsQ(ctx).FilterByID(ID).Get(ctx)
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return models.Organization{}, errx.ErrorOrganizationNotFound.Raise(
			fmt.Errorf("organization with ID %s not found", ID),
		)
	case err != nil:
		return models.Organization{}, fmt.Errorf(
			"failed to get organization with ID %s, cause: %w", ID, err,
		)
	}

	return Organization(row), nil
}

func (r Repository) DeleteOrganization(ctx context.Context, ID uuid.UUID) error {
	err := r.organizationsQ(ctx).FilterByID(ID).Delete(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete organization with ID %s, cause: %w", ID, err)
	}

	return nil
}

func (r Repository) GetOrganizations(
	ctx context.Context,
	filter organization.FilterParams,
	limit, offset uint,
) (pagi.Page[[]models.Organization], error) {
	if limit == 0 {
		limit = 10
	}

	q := r.organizationsQ(ctx)
	if filter.Name != nil {
		q = q.FilterNameLike(*filter.Name)
	}
	if filter.Status != nil {
		q = q.FilterByStatus(*filter.Status)
	}

	rows, err := q.Page(limit, offset).Select(ctx)
	if err != nil {
		return pagi.Page[[]models.Organization]{}, fmt.Errorf("failed to get organizations, cause: %w", err)
	}

	total, err := q.Count(ctx)
	if err != nil {
		return pagi.Page[[]models.Organization]{}, fmt.Errorf("failed to count organizations, cause: %w", err)
	}

	organizations := make([]models.Organization, len(rows))
	for i, row := range rows {
		organizations[i] = Organization(row)
	}

	return pagi.Page[[]models.Organization]{
		Data:  organizations,
		Page:  uint(offset/limit) + 1,
		Size:  uint(len(organizations)),
		Total: total,
	}, nil

}

func (r Repository) GetOrganizationsForUser(
	ctx context.Context,
	accountID uuid.UUID,
	limit, offset uint,
) (pagi.Page[[]models.Organization], error) {
	if limit == 0 {
		limit = 10
	}

	row, err := r.organizationsQ(ctx).FilterByAccountID(accountID).Page(limit, offset).Select(ctx)
	if err != nil {
		return pagi.Page[[]models.Organization]{}, fmt.Errorf(
			"failed to get organizations for accountID %s, cause: %w", accountID, err,
		)
	}

	total, err := r.organizationsQ(ctx).FilterByAccountID(accountID).Count(ctx)
	if err != nil {
		return pagi.Page[[]models.Organization]{}, fmt.Errorf(
			"failed to count organizations for accountID %s, cause: %w", accountID, err,
		)
	}

	organizations := make([]models.Organization, len(row))
	for i, r := range row {
		organizations[i] = Organization(r)
	}

	return pagi.Page[[]models.Organization]{
		Data:  organizations,
		Page:  uint(offset/limit) + 1,
		Size:  uint(len(organizations)),
		Total: total,
	}, nil
}

func Organization(db pgdb.Organization) models.Organization {
	return models.Organization{
		ID:        db.ID,
		Status:    db.Status,
		Name:      db.Name,
		MaxRoles:  db.MaxRoles,
		CreatedAt: db.CreatedAt,
		UpdatedAt: db.UpdatedAt,
	}
}
