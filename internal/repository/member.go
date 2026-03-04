package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/models"
	"github.com/netbill/restkit/pagi"
)

type MemberRepo struct {
	OrgMembersSql OrgMembersQ
}

type OrganizationMemberRow struct {
	ID             uuid.UUID `db:"id"`
	AccountID      uuid.UUID `db:"account_id"`
	OrganizationID uuid.UUID `db:"organization_id"`
	Head           bool      `db:"head"`
	Position       *string   `db:"position,omitempty"`
	Label          *string   `db:"label,omitempty"`
	Version        int32     `db:"version"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}

func (r OrganizationMemberRow) IsNil() bool {
	return r.ID == uuid.Nil
}

func (r OrganizationMemberRow) ToModel() models.Member {
	return models.Member{
		ID:             r.ID,
		AccountID:      r.AccountID,
		OrganizationID: r.OrganizationID,
		Head:           r.Head,
		Position:       r.Position,
		Label:          r.Label,
		Version:        r.Version,
		CreatedAt:      r.CreatedAt,
		UpdatedAt:      r.UpdatedAt,
	}
}

type OrgMembersQ interface {
	New() OrgMembersQ

	Insert(ctx context.Context, params OrganizationMemberRow) (OrganizationMemberRow, error)

	Get(ctx context.Context) (OrganizationMemberRow, error)
	Select(ctx context.Context) ([]OrganizationMemberRow, error)
	Exists(ctx context.Context) (bool, error)

	UpdateOne(ctx context.Context) (OrganizationMemberRow, error)

	UpdatePosition(position *string) OrgMembersQ
	UpdateLabel(label *string) OrgMembersQ

	Delete(ctx context.Context) error

	FilterByID(id uuid.UUID) OrgMembersQ
	FilterByAccountID(accountID uuid.UUID) OrgMembersQ
	FilterByOrganizationID(organizationID ...uuid.UUID) OrgMembersQ
	FilterByUsername(username string) OrgMembersQ

	FilterBestMatch(term string) OrgMembersQ
	FilterLikePseudonym(pseudonym string) OrgMembersQ
	FilterLikeUsername(username string) OrgMembersQ
	FilterLikeLabel(label string) OrgMembersQ
	FilterLikePosition(position string) OrgMembersQ
	FilterByHead(head bool) OrgMembersQ

	Page(limit, offset uint) OrgMembersQ
	Count(ctx context.Context) (uint, error)
}

func (r *MemberRepo) CreateMember(
	ctx context.Context,
	accountID, organizationID uuid.UUID,
) (models.Member, error) {
	row, err := r.OrgMembersSql.New().Insert(ctx, OrganizationMemberRow{
		AccountID:      accountID,
		OrganizationID: organizationID,
	})
	if err != nil {
		return models.Member{}, fmt.Errorf("failed to create member, cause: %w", err)
	}

	return r.GetMember(ctx, row.ID)
}

func (r *MemberRepo) CreateMemberHead(
	ctx context.Context,
	accountID, organizationID uuid.UUID,
) (models.Member, error) {
	row, err := r.OrgMembersSql.New().Insert(ctx, OrganizationMemberRow{
		AccountID:      accountID,
		OrganizationID: organizationID,
		Head:           true,
	})
	if err != nil {
		return models.Member{}, fmt.Errorf("failed to create member, cause: %w", err)
	}

	return r.GetMember(ctx, row.ID)
}

func (r *MemberRepo) UpdateMember(
	ctx context.Context,
	ID uuid.UUID,
	params core.MemberUpdateParams,
) (models.Member, error) {
	q := r.OrgMembersSql.New().
		FilterByID(ID)

	if params.Position != nil {
		q = q.UpdatePosition(params.Position)
	}
	if params.Label != nil {
		q = q.UpdateLabel(params.Label)
	}

	row, err := q.UpdateOne(ctx)
	if err != nil {
		return models.Member{}, fmt.Errorf("failed to update member, cause: %w", err)
	}
	if row.IsNil() {
		return models.Member{}, errx.ErrorMemberNotExists.Raise(
			fmt.Errorf("member with ID %s not found", ID),
		)
	}

	return r.GetMember(ctx, row.ID)
}

func (r *MemberRepo) GetMember(
	ctx context.Context,
	memberID uuid.UUID,
) (models.Member, error) {
	row, err := r.OrgMembersSql.New().FilterByID(memberID).Get(ctx)
	if err != nil {
		return models.Member{}, fmt.Errorf("failed to get member, cause: %w", err)
	}
	if row.IsNil() {
		return models.Member{}, errx.ErrorMemberNotExists.Raise(
			fmt.Errorf("member with ID %s not found", memberID),
		)
	}

	return row.ToModel(), nil
}

func (r *MemberRepo) GetMembersByAccountAndOrgs(
	ctx context.Context,
	accountID uuid.UUID,
	organizationIDs []uuid.UUID,
) ([]models.Member, error) {
	q := r.OrgMembersSql.New().
		FilterByAccountID(accountID).
		FilterByOrganizationID(organizationIDs...)

	rows, err := q.Select(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get members by account and organizations, cause: %w", err)
	}

	members := make([]models.Member, 0, len(rows))
	for _, row := range rows {
		members = append(members, row.ToModel())
	}

	return members, nil
}

func (r *MemberRepo) GetMemberByAccountAndOrganization(
	ctx context.Context,
	accountID, organizationID uuid.UUID,
) (models.Member, error) {
	row, err := r.OrgMembersSql.New().
		FilterByAccountID(accountID).
		FilterByOrganizationID(organizationID).
		Get(ctx)
	if err != nil {
		return models.Member{}, fmt.Errorf("failed to get member by account and organization, cause: %w", err)
	}
	if row.IsNil() {
		return models.Member{}, errx.ErrorMemberNotExists.Raise(
			fmt.Errorf("member with account ID %s and organization ID %s not found", accountID, organizationID),
		)
	}

	return row.ToModel(), nil
}

func (r *MemberRepo) MemberExists(
	ctx context.Context,
	accountID, organizationID uuid.UUID,
) (bool, error) {
	exists, err := r.OrgMembersSql.New().
		FilterByAccountID(accountID).
		FilterByOrganizationID(organizationID).
		Exists(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to check member existence, cause: %w", err)
	}

	return exists, nil
}

func (r *MemberRepo) GetMembers(
	ctx context.Context,
	filter core.MemberFilterParams,
	limit uint,
	offset uint,
) (pagi.Page[[]models.Member], error) {
	if limit == 0 {
		limit = 20
	}
	if limit > 1000 {
		limit = 1000
	}

	q := r.OrgMembersSql.New()
	if filter.OrganizationID != nil {
		q = q.FilterByOrganizationID(*filter.OrganizationID)
	}
	if filter.AccountID != nil {
		q = q.FilterByAccountID(*filter.AccountID)
	}
	if filter.Username != nil {
		q = q.FilterByUsername(*filter.Username)
	}
	if filter.BestMatch != nil {
		q = q.FilterBestMatch(*filter.BestMatch)
	}
	if filter.Label != nil {
		q = q.FilterLikeLabel(*filter.Label)
	}
	if filter.Position != nil {
		q = q.FilterLikePosition(*filter.Position)
	}
	if filter.Head != nil {
		q = q.FilterByHead(*filter.Head)
	}

	rows, err := q.Page(limit, offset).Select(ctx)
	if err != nil {
		return pagi.Page[[]models.Member]{}, fmt.Errorf("failed to get members, cause: %w", err)
	}

	total, err := q.Count(ctx)
	if err != nil {
		return pagi.Page[[]models.Member]{}, fmt.Errorf("failed to count members, cause: %w", err)
	}

	collection := make([]models.Member, 0, len(rows))
	for _, row := range rows {
		collection = append(collection, row.ToModel())
	}

	return pagi.Page[[]models.Member]{
		Data:  collection,
		Page:  uint(offset/limit) + 1,
		Size:  uint(len(collection)),
		Total: total,
	}, nil
}

func (r *MemberRepo) DeleteMember(
	ctx context.Context,
	memberID uuid.UUID,
) error {
	err := r.OrgMembersSql.New().FilterByID(memberID).Delete(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete member with ID %s, cause: %w", memberID, err)
	}

	return nil
}

func (r *MemberRepo) DeleteMembersByAccountID(
	ctx context.Context,
	accountID uuid.UUID,
) error {
	err := r.OrgMembersSql.New().FilterByAccountID(accountID).Delete(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete members with account ID %s, cause: %w", accountID, err)
	}

	return nil
}
