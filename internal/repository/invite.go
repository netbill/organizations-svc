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

type OrgInviteRow struct {
	ID             uuid.UUID `db:"id"`
	OrganizationID uuid.UUID `db:"organization_id"`
	AccountID      uuid.UUID `db:"account_id,omitempty"`
	Status         string    `db:"status"`
	UpdatedAt      time.Time `db:"updated_at"`
	ExpiresAt      time.Time `db:"expires_at"`
	CreatedAt      time.Time `db:"created_at"`
}

func (r OrgInviteRow) IsNil() bool {
	return r.ID == uuid.Nil
}

func (r OrgInviteRow) ToModel() models.Invite {
	return models.Invite{
		ID:             r.ID,
		OrganizationID: r.OrganizationID,
		AccountID:      r.AccountID,
		Status:         r.Status,
		UpdatedAt:      r.UpdatedAt,
		ExpiresAt:      r.ExpiresAt,
		CreatedAt:      r.CreatedAt,
	}
}

type OrgInvitesQ interface {
	New() OrgInvitesQ
	Insert(ctx context.Context, input OrgInviteRow) (OrgInviteRow, error)

	Get(ctx context.Context) (OrgInviteRow, error)
	Select(ctx context.Context) ([]OrgInviteRow, error)
	Exists(ctx context.Context) (bool, error)

	UpdateOne(ctx context.Context) (OrgInviteRow, error)

	UpdateStatus(status string) OrgInvitesQ

	FilterByID(id uuid.UUID) OrgInvitesQ
	FilterByOrganizationID(organizationID uuid.UUID) OrgInvitesQ
	FilterByAccountID(accountID uuid.UUID) OrgInvitesQ
	FilterByStatus(status string) OrgInvitesQ
	FilterExpiresBefore(t time.Time) OrgInvitesQ
	FilterExpiresAfter(t time.Time) OrgInvitesQ

	Page(limit, offset uint) OrgInvitesQ

	Delete(ctx context.Context) error
	Count(ctx context.Context) (uint, error)
}

type InviteRepo struct {
	query OrgInvitesQ
}

func NewInviteRepo(orgInvitesSql OrgInvitesQ) *InviteRepo {
	return &InviteRepo{
		query: orgInvitesSql,
	}
}

func (r *InviteRepo) Create(
	ctx context.Context,
	params core.InviteCreateParams,
) (models.Invite, error) {
	row, err := r.query.New().Insert(ctx, OrgInviteRow{
		OrganizationID: params.OrganizationID,
		AccountID:      params.AccountID,
		ExpiresAt:      params.ExpiresAt,
	})
	if err != nil {
		return models.Invite{}, fmt.Errorf(
			"failed to create invite for organization ID %s and account ID %s cause: %w",
			params.OrganizationID, params.AccountID, err,
		)
	}

	return row.ToModel(), nil
}

func (r *InviteRepo) Get(
	ctx context.Context,
	inviteID uuid.UUID,
) (models.Invite, error) {
	row, err := r.query.New().FilterByID(inviteID).Get(ctx)
	if err != nil {
		return models.Invite{}, fmt.Errorf("failed to get invite with ID %s, cause: %w", inviteID, err)
	}
	if row.IsNil() {
		return models.Invite{}, errx.ErrorInviteNotExists.Raise(
			fmt.Errorf("invite with ID %s not found", inviteID),
		)
	}

	return row.ToModel(), nil
}

func (r *InviteRepo) UpdateStatus(
	ctx context.Context,
	inviteID uuid.UUID,
	status string,
) (models.Invite, error) {
	row, err := r.query.New().FilterByID(inviteID).UpdateStatus(status).UpdateOne(ctx)
	if err != nil {
		return models.Invite{}, fmt.Errorf("failed to update invite status with ID %s, cause: %w", inviteID, err)
	}
	if row.IsNil() {
		return models.Invite{}, errx.ErrorInviteNotExists.Raise(
			fmt.Errorf("invite with ID %s not found", inviteID),
		)
	}

	return row.ToModel(), nil
}

func (r *InviteRepo) Delete(
	ctx context.Context,
	inviteID uuid.UUID,
) error {
	err := r.query.New().FilterByID(inviteID).Delete(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete invite with inviteID %s, cause: %w", inviteID, err)
	}

	return nil
}

func (r *InviteRepo) GetList(
	ctx context.Context,
	params core.FilterInvitesParams,
	limit, offset uint,
) (pagi.Page[[]models.Invite], error) {
	if limit == 0 {
		limit = 20
	}
	if limit > 1000 {
		limit = 1000
	}

	q := r.query.New()
	if params.OrganizationID != nil {
		q = q.FilterByOrganizationID(*params.OrganizationID)
	}
	if params.AccountID != nil {
		q = q.FilterByAccountID(*params.AccountID)
	}

	rows, err := q.Page(limit, offset).Select(ctx)
	if err != nil {
		return pagi.Page[[]models.Invite]{}, fmt.Errorf(
			"failed to get invites for filter params %+v, cause: %w", params, err,
		)
	}

	total, err := q.Count(ctx)
	if err != nil {
		return pagi.Page[[]models.Invite]{}, fmt.Errorf(
			"failed to count invites for filter params %+v, cause: %w", params, err,
		)
	}

	res := make([]models.Invite, 0, len(rows))
	for _, row := range rows {
		res = append(res, row.ToModel())
	}

	return pagi.Page[[]models.Invite]{
		Data:  res,
		Page:  uint(offset/limit) + 1,
		Size:  uint(len(res)),
		Total: total,
	}, nil
}

func (r *InviteRepo) ExistActiveForAccountInOrg(
	ctx context.Context,
	accountID, organizationID uuid.UUID,
) (bool, error) {
	exists, err := r.query.New().
		FilterByAccountID(accountID).
		FilterByOrganizationID(organizationID).
		FilterExpiresAfter(time.Now().UTC()).
		Exists(ctx)
	if err != nil {
		return false, fmt.Errorf(
			"failed to check active invite existence for account ID %s and organization ID %s, cause: %w",
			accountID, organizationID, err,
		)
	}

	return exists, nil
}
