package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/models"
)

type OrgRoleRankRow struct {
	RoleID         uuid.UUID `db:"role_id"`
	OrganizationID uuid.UUID `db:"organization_id"`
	Rank           int32     `db:"rank"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (r OrgRoleRankRow) IsNil() bool {
	return r.RoleID == uuid.Nil
}

func (r OrgRoleRankRow) ToModel() models.OrgRoleRank {
	return models.OrgRoleRank{
		RoleID:         r.RoleID,
		OrganizationID: r.OrganizationID,
		Rank:           r.Rank,
		CreatedAt:      r.CreatedAt,
		UpdatedAt:      r.UpdatedAt,
	}
}

type OrgRoleRanksQ interface {
	New() OrgRoleRanksQ

	Insert(ctx context.Context, input OrgRoleRankRow) ([]OrgRoleRankRow, error)

	Get(ctx context.Context) (OrgRoleRankRow, error)
	Select(ctx context.Context) ([]OrgRoleRankRow, error)
	Exists(ctx context.Context) (bool, error)

	UpdateMany(ctx context.Context) (int64, error)
	UpdateOne(ctx context.Context) (OrgRoleRankRow, error)

	FilterRoleID(roleID uuid.UUID) OrgRoleRanksQ
	FilterOrganizationID(organizationID uuid.UUID) OrgRoleRanksQ
}

func (r *Repository) CreateRoleRank(ctx context.Context, roleID uuid.UUID, rank int32) ([]models.OrgRoleRank, error) {
	res, err := r.OrgRoleRanksSql.New().Insert(ctx, OrgRoleRankRow{
		RoleID: roleID,
		Rank:   rank,
	})
	if err != nil {
		return nil, err
	}

	var ranks []models.OrgRoleRank

	for _, r := range res {
		ranks = append(ranks, models.OrgRoleRank{
			RoleID:         r.RoleID,
			OrganizationID: r.OrganizationID,
			Rank:           r.Rank,
			CreatedAt:      r.CreatedAt,
			UpdatedAt:      r.UpdatedAt,
		})
	}

	return ranks, nil
}

func (r *Repository) GetRoleRank(ctx context.Context, roleID uuid.UUID) (models.OrgRoleRank, error) {
	res, err := r.OrgRoleRanksSql.New().FilterRoleID(roleID).Get(ctx)
	if err != nil {
		return models.OrgRoleRank{}, err
	}

	return res.ToModel(), nil
}

func (r *Repository) GetOrgRolesRanks(ctx context.Context, organizationID uuid.UUID) ([]models.OrgRoleRank, error) {
	res, err := r.OrgRoleRanksSql.New().FilterOrganizationID(organizationID).Select(ctx)
	if err != nil {
		return nil, err
	}

	var ranks []models.OrgRoleRank

	for _, r := range res {
		ranks = append(ranks, models.OrgRoleRank{
			RoleID:         r.RoleID,
			OrganizationID: r.OrganizationID,
			Rank:           r.Rank,
			CreatedAt:      r.CreatedAt,
			UpdatedAt:      r.UpdatedAt,
		})
	}

	return ranks, nil
}

func (r *Repository) UpdateRolesRanks(ctx context.Context, organizationID uuid.UUID, ranks map[int32]uuid.UUID) ([]models.OrgRoleRank, error) {
	
}

func (r *Repository) DeleteRoleRank(ctx context.Context, roleID uuid.UUID) ([]models.OrgRoleRank, error) {
	//TODO implement me
	panic("implement me")
}

// Revisions

func (r *Repository) CreateOrgRoleRankRevision(ctx context.Context, organizationID uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (r *Repository) LockOrgRoleRankRevision(ctx context.Context, organizationID uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (r *Repository) BumpOrgRoleRankRevision(ctx context.Context, organizationID uuid.UUID) (models.OrgRoleRanksRevision, error) {
	//TODO implement me
	panic("implement me")
}
