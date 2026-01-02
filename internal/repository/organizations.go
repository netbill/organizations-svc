package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/models"
	"github.com/netbill/organizations-svc/internal/core/modules/organization"
	"github.com/netbill/organizations-svc/internal/repository/pgdb"
	"github.com/netbill/pagi"
	"github.com/pkg/errors"
)

func (s Service) CreateOrganization(
	ctx context.Context,
	params organization.CreateParams,
) (models.Organization, error) {
	row, err := s.organizationsQ().Insert(ctx, pgdb.OrganizationsQInsertInput{
		Name: params.Name,
		Icon: params.Icon,
	})
	if err != nil {
		return models.Organization{}, err
	}

	return Organization(row), nil
}

func (s Service) UpdateOrganization(
	ctx context.Context,
	ID uuid.UUID,
	params organization.UpdateParams,
) (models.Organization, error) {
	q := s.organizationsQ().FilterByID(ID)
	if params.Name != nil {
		q = q.UpdateName(*params.Name)
	}
	if params.Icon != nil {
		q = q.UpdateIcon(*params.Icon)
	}

	row, err := q.UpdateOne(ctx)
	if err != nil {
		return models.Organization{}, err
	}

	return Organization(row), nil
}

func (s Service) UpdateOrganizationStatus(
	ctx context.Context,
	ID uuid.UUID,
	status string,
) (models.Organization, error) {
	row, err := s.organizationsQ().FilterByID(ID).UpdateStatus(status).UpdateOne(ctx)
	if err != nil {
		return models.Organization{}, err
	}

	return Organization(row), nil
}

func (s Service) UpdateOrganizationMaxRoles(
	ctx context.Context,
	ID uuid.UUID,
	maxRoles uint,
) (models.Organization, error) {
	row, err := s.organizationsQ().FilterByID(ID).UpdateMaxRoles(maxRoles).UpdateOne(ctx)
	if err != nil {
		return models.Organization{}, err
	}

	return Organization(row), nil
}

func (s Service) GetOrganizationByID(ctx context.Context, ID uuid.UUID) (models.Organization, error) {
	row, err := s.organizationsQ().FilterByID(ID).Get(ctx)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return models.Organization{}, nil
	case err != nil:
		return models.Organization{}, err
	}

	return Organization(row), nil
}

func (s Service) DeleteOrganization(ctx context.Context, ID uuid.UUID) error {
	return s.organizationsQ().FilterByID(ID).Delete(ctx)
}

func (s Service) GetOrganizations(
	ctx context.Context,
	filter organization.FilterParams,
	offset, limit uint,
) (pagi.Page[[]models.Organization], error) {
	q := s.organizationsQ()
	if filter.Name != nil {
		q = q.FilterNameLike(*filter.Name)
	}
	if filter.Status != nil {
		q = q.FilterByStatus(*filter.Status)
	}

	rows, err := q.Page(limit, offset).Select(ctx)
	if err != nil {
		return pagi.Page[[]models.Organization]{}, err
	}

	total, err := q.Count(ctx)
	if err != nil {
		return pagi.Page[[]models.Organization]{}, err
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

func (s Service) GetOrganizationsForUser(
	ctx context.Context,
	accountID uuid.UUID,
	limit, offset uint,
) (pagi.Page[[]models.Organization], error) {
	row, err := s.organizationsQ().FilterByID(accountID).Page(limit, offset).Select(ctx)
	if err != nil {
		return pagi.Page[[]models.Organization]{}, err
	}

	total, err := s.organizationsQ().FilterByID(accountID).Count(ctx)
	if err != nil {
		return pagi.Page[[]models.Organization]{}, err
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
		Icon:      db.Icon,
		MaxRoles:  db.MaxRoles,
		CreatedAt: db.CreatedAt,
		UpdatedAt: db.UpdatedAt,
	}
}
