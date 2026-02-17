package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/netbill/organizations-svc/internal/core/errx"
	"github.com/netbill/organizations-svc/internal/core/models"
	"github.com/netbill/organizations-svc/internal/core/modules/invite"
	"github.com/netbill/restkit/pagi"
)

type OrgInviteRow struct {
	ID             uuid.UUID `db:"id"`
	OrganizationID uuid.UUID `db:"organization_id"`
	AccountID      uuid.UUID `db:"account_id,omitempty"`
	Status         string    `db:"status"`
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
		ExpiresAt:      r.ExpiresAt,
		CreatedAt:      r.CreatedAt,
	}
}

type OrgInvitesQ interface {
	New() OrgInvitesQ
	Insert(ctx context.Context, input OrgInviteRow) (OrgInviteRow, error)

	Get(ctx context.Context) (OrgInviteRow, error)
	Select(ctx context.Context) ([]OrgInviteRow, error)

	UpdateMany(ctx context.Context) (int64, error)
	UpdateOne(ctx context.Context) (OrgInviteRow, error)

	UpdateStatus(status string) OrgInvitesQ

	FilterByID(id uuid.UUID) OrgInvitesQ
	FilterByOrganizationID(organizationID uuid.UUID) OrgInvitesQ
	FilterByAccountID(accountID uuid.UUID) OrgInvitesQ

	Page(limit, offset uint) OrgInvitesQ

	Delete(ctx context.Context) error
	Count(ctx context.Context) (uint, error)
}

func (r *Repository) CreateInvite(
	ctx context.Context,
	params invite.CreateParams,
) (models.Invite, error) {
	row, err := r.OrgInvitesSql.New().Insert(ctx, OrgInviteRow{
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

func (r *Repository) GetInvite(
	ctx context.Context,
	inviteID uuid.UUID,
) (models.Invite, error) {
	row, err := r.OrgInvitesSql.New().FilterByID(inviteID).Get(ctx)
	if err != nil {
		return models.Invite{}, fmt.Errorf("failed to get invite with ID %s, cause: %w", inviteID, err)
	}
	if row.IsNil() {
		return models.Invite{}, errx.ErrorInviteNotFound.Raise(
			fmt.Errorf("invite with ID %s not found", inviteID),
		)
	}

	return row.ToModel(), nil
}

func (r *Repository) UpdateInviteStatus(
	ctx context.Context,
	inviteID uuid.UUID,
	status string,
) (models.Invite, error) {
	row, err := r.OrgInvitesSql.New().FilterByID(inviteID).UpdateStatus(status).UpdateOne(ctx)
	if err != nil {
		return models.Invite{}, fmt.Errorf("failed to update invite status with ID %s, cause: %w", inviteID, err)
	}
	if row.IsNil() {
		return models.Invite{}, errx.ErrorInviteNotFound.Raise(
			fmt.Errorf("invite with ID %s not found", inviteID),
		)
	}

	return row.ToModel(), nil
}

func (r *Repository) DeleteInvite(
	ctx context.Context,
	inviteID uuid.UUID,
) error {
	err := r.OrgInvitesSql.New().FilterByID(inviteID).Delete(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete invite with inviteID %s, cause: %w", inviteID, err)
	}

	return nil
}

func (r *Repository) GetOrganizationInvites(
	ctx context.Context,
	organizationID uuid.UUID,
	limit, offset uint,
) (pagi.Page[[]models.Invite], error) {
	if limit == 0 {
		limit = 10
	}

	rows, err := r.OrgInvitesSql.New().
		FilterByOrganizationID(organizationID).
		Page(limit, offset).
		Select(ctx)
	if err != nil {
		return pagi.Page[[]models.Invite]{}, fmt.Errorf(
			"failed to get invites for organization ID %s, cause: %w", organizationID, err,
		)
	}

	total, err := r.OrgInvitesSql.New().
		FilterByOrganizationID(organizationID).
		Count(ctx)
	if err != nil {
		return pagi.Page[[]models.Invite]{}, fmt.Errorf(
			"failed to count invites for organization ID %s, cause: %w", organizationID, err,
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

func (r *Repository) GetAccountInvites(
	ctx context.Context,
	accountID uuid.UUID,
	limit, offset uint,
) (pagi.Page[[]models.Invite], error) {
	if limit == 0 {
		limit = 10
	}

	rows, err := r.OrgInvitesSql.New().
		FilterByAccountID(accountID).
		Page(limit, offset).
		Select(ctx)
	if err != nil {
		return pagi.Page[[]models.Invite]{}, fmt.Errorf(
			"failed to get invites for account ID %s, cause: %w", accountID, err,
		)
	}

	total, err := r.OrgInvitesSql.New().
		FilterByAccountID(accountID).
		Count(ctx)
	if err != nil {
		return pagi.Page[[]models.Invite]{}, fmt.Errorf(
			"failed to count invites for account ID %s, cause: %w", accountID, err,
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
