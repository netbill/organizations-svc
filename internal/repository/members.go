package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/domain"
	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/modules/member"
	"github.com/netbill/restkit/pagi"
)

type OrganizationMemberRow struct {
	ID             uuid.UUID `db:"id"`
	AccountID      uuid.UUID `db:"account_id"`
	OrganizationID uuid.UUID `db:"organization_id"`
	Head           bool      `db:"head"`
	Position       *string   `db:"position,omitempty"`
	Label          *string   `db:"label,omitempty"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}

func (r OrganizationMemberRow) IsNil() bool {
	return r.ID == uuid.Nil
}

type OrganizationMemberWithProfileDataRow struct {
	OrganizationMemberRow
	Username  string  `db:"username"`
	Official  bool    `db:"official"`
	Pseudonym *string `db:"pseudonym,omitempty"`
	Icon      *string `db:"icon,omitempty"`
}

func (r OrganizationMemberWithProfileDataRow) IsNil() bool {
	return r.ID == uuid.Nil
}

func (r OrganizationMemberWithProfileDataRow) ToModel() domain.Member {
	return domain.Member{
		ID:             r.ID,
		AccountID:      r.AccountID,
		OrganizationID: r.OrganizationID,
		Head:           r.Head,
		Position:       r.Position,
		Label:          r.Label,
		Username:       r.Username,
		Official:       r.Official,
		Pseudonym:      r.Pseudonym,
		Icon:           r.Icon,
		CreatedAt:      r.CreatedAt,
		UpdatedAt:      r.UpdatedAt,
	}
}

type OrgMembersQ interface {
	New() OrgMembersQ
	Insert(ctx context.Context, params OrganizationMemberRow) (OrganizationMemberRow, error)

	FilterByID(id uuid.UUID) OrgMembersQ
	FilterByAccountID(accountID uuid.UUID) OrgMembersQ
	FilterByOrganizationID(organizationID uuid.UUID) OrgMembersQ
	FilterByUsername(username string) OrgMembersQ

	FilterRoleID(roleID uuid.UUID) OrgMembersQ
	FilterByPermissionID(permissionID uuid.UUID) OrgMembersQ
	FilterByRoleRankUp(rank uint) OrgMembersQ
	FilterByRoleRankDown(rank uint) OrgMembersQ

	FilterBestMatch(term string) OrgMembersQ
	FilterLikePseudonym(pseudonym string) OrgMembersQ
	FilterLikeUsername(username string) OrgMembersQ
	FilterLikeLabel(label string) OrgMembersQ
	FilterLikePosition(position string) OrgMembersQ
	FilterByHead(head bool) OrgMembersQ

	UpdatePosition(position *string) OrgMembersQ
	UpdateLabel(label *string) OrgMembersQ
	UpdateOne(ctx context.Context) (OrganizationMemberRow, error)
	UpdateMany(ctx context.Context) (int64, error)

	GetWithUserData(ctx context.Context) (OrganizationMemberWithProfileDataRow, error)
	SelectWithUserData(ctx context.Context) ([]OrganizationMemberWithProfileDataRow, error)

	Page(limit, offset uint) OrgMembersQ
	Count(ctx context.Context) (uint, error)
	Exists(ctx context.Context) (bool, error)
	Delete(ctx context.Context) error
}

func (r *Repository) CreateMember(
	ctx context.Context,
	accountID, organizationID uuid.UUID,
) (domain.Member, error) {
	row, err := r.OrgMembersSql.New().Insert(ctx, OrganizationMemberRow{
		AccountID:      accountID,
		OrganizationID: organizationID,
	})
	if err != nil {
		return domain.Member{}, fmt.Errorf("failed to create member, cause: %w", err)
	}

	return r.GetMember(ctx, row.ID)
}

func (r *Repository) CreateMemberHead(
	ctx context.Context,
	accountID, organizationID uuid.UUID,
) (domain.Member, error) {
	row, err := r.OrgMembersSql.New().Insert(ctx, OrganizationMemberRow{
		AccountID:      accountID,
		OrganizationID: organizationID,
		Head:           true,
	})
	if err != nil {
		return domain.Member{}, fmt.Errorf("failed to create member, cause: %w", err)
	}

	return r.GetMember(ctx, row.ID)
}

func (r *Repository) UpdateMember(
	ctx context.Context,
	ID uuid.UUID,
	params member.UpdateParams,
) (domain.Member, error) {
	row, err := r.OrgMembersSql.New().
		FilterByID(ID).
		UpdatePosition(params.Position).
		UpdateLabel(params.Label).
		UpdateOne(ctx)
	if err != nil {
		return domain.Member{}, fmt.Errorf("failed to update member, cause: %w", err)
	}
	if row.IsNil() {
		return domain.Member{}, errx.ErrorMemberNotFound.Raise(
			fmt.Errorf("member with ID %s not found", ID),
		)
	}

	return r.GetMember(ctx, row.ID)
}

func (r *Repository) GetMember(
	ctx context.Context,
	memberID uuid.UUID,
) (domain.Member, error) {
	row, err := r.OrgMembersSql.New().FilterByID(memberID).GetWithUserData(ctx)
	if err != nil {
		return domain.Member{}, fmt.Errorf("failed to get member, cause: %w", err)
	}
	if row.IsNil() {
		return domain.Member{}, errx.ErrorMemberNotFound.Raise(
			fmt.Errorf("member with ID %s not found", memberID),
		)
	}

	return row.ToModel(), nil
}

func (r *Repository) GetMemberByAccountAndOrganization(
	ctx context.Context,
	accountID, organizationID uuid.UUID,
) (domain.Member, error) {
	row, err := r.OrgMembersSql.New().
		FilterByAccountID(accountID).
		FilterByOrganizationID(organizationID).
		GetWithUserData(ctx)
	if err != nil {
		return domain.Member{}, fmt.Errorf("failed to get member by account and organization, cause: %w", err)
	}
	if row.IsNil() {
		return domain.Member{}, errx.ErrorMemberNotFound.Raise(
			fmt.Errorf("member with account ID %s and organization ID %s not found", accountID, organizationID),
		)
	}

	return row.ToModel(), nil
}

func (r *Repository) MemberExists(
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

func (r *Repository) GetMembers(
	ctx context.Context,
	filter member.FilterParams,
	limit uint,
	offset uint,
) (pagi.Page[[]domain.Member], error) {
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
	if filter.RoleID != nil {
		q = q.FilterRoleID(*filter.RoleID)
	}
	if filter.PermissionID != nil {
		q = q.FilterByPermissionID(*filter.PermissionID)
	}
	if filter.RoleRankUp != nil {
		q = q.FilterByRoleRankUp(*filter.RoleRankUp)
	}
	if filter.RoleRankDown != nil {
		q = q.FilterByRoleRankDown(*filter.RoleRankDown)
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

	if limit == 0 {
		limit = 10
	}

	rows, err := q.Page(limit, offset).SelectWithUserData(ctx)
	if err != nil {
		return pagi.Page[[]domain.Member]{}, fmt.Errorf("failed to get members, cause: %w", err)
	}

	total, err := q.Count(ctx)
	if err != nil {
		return pagi.Page[[]domain.Member]{}, fmt.Errorf("failed to count members, cause: %w", err)
	}

	collection := make([]domain.Member, 0, len(rows))
	for _, row := range rows {
		collection = append(collection, row.ToModel())
	}

	return pagi.Page[[]domain.Member]{
		Data:  collection,
		Page:  uint(offset/limit) + 1,
		Size:  uint(len(collection)),
		Total: total,
	}, nil
}

func (r *Repository) DeleteMember(
	ctx context.Context,
	memberID uuid.UUID,
) error {
	err := r.OrgMembersSql.New().FilterByID(memberID).Delete(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete member with ID %s, cause: %w", memberID, err)
	}

	return nil
}

func (r *Repository) DeleteMembersByAccountID(
	ctx context.Context,
	accountID uuid.UUID,
) error {
	err := r.OrgMembersSql.New().FilterByAccountID(accountID).Delete(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete members with account ID %s, cause: %w", accountID, err)
	}

	return nil
}

func (r *Repository) GetMemberMaxRole(
	ctx context.Context,
	memberID uuid.UUID,
) (domain.Role, error) {
	row, err := r.OrgRolesSql.New().
		FilterByMemberID(memberID).
		OrderByRoleRank(false). // DESC => max
		Get(ctx)
	if err != nil {
		return domain.Role{}, fmt.Errorf("failed to get member max role, cause: %w", err)
	}
	if row.IsNil() {
		return domain.Role{}, errx.ErrorRoleNotFound.Raise(
			fmt.Errorf("no roles found for member with ID %s", memberID),
		)
	}

	return row.ToModel(), nil
}

func (r *Repository) CheckMemberHavePermission(
	ctx context.Context,
	memberID uuid.UUID,
	permissionID uuid.UUID,
) (bool, error) {
	have, err := r.OrgMembersSql.New().
		FilterByID(memberID).
		FilterByPermissionID(permissionID).
		Exists(ctx)
	if err != nil {
		return false, fmt.Errorf("checking member permission: %w", err)
	}

	return have, nil
}
