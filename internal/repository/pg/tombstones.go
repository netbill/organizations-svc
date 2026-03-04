package pg

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/repository"
	"github.com/netbill/pgdbx"
)

type tombstones struct {
	db *pgdbx.DB
}

func NewTombstonesQ(db *pgdbx.DB) repository.TombstoneRepo {
	return &tombstones{db: db}
}

func (t *tombstones) BuryProfile(ctx context.Context, accountID uuid.UUID) error {
	_, err := t.db.Exec(ctx, `
        INSERT INTO tombstones (entity_type, entity_id)
        VALUES ('profile', $1)
        ON CONFLICT (entity_type, entity_id) DO NOTHING
    `, accountID)
	if err != nil {
		return fmt.Errorf("burying profile: %w", err)
	}
	return nil
}

func (t *tombstones) ProfileIsBuried(ctx context.Context, accountID uuid.UUID) (bool, error) {
	var exists bool
	err := t.db.QueryRow(ctx, `
        SELECT EXISTS(
            SELECT 1 FROM tombstones WHERE entity_type = 'profile' AND entity_id = $1
        )
    `, accountID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("checking if profile is buried: %w", err)
	}
	return exists, nil
}

func (t *tombstones) BuryOrganization(ctx context.Context, orgID uuid.UUID) error {
	_, err := t.db.Exec(ctx, `
        INSERT INTO tombstones (entity_type, entity_id)
        SELECT 'organization', $1
        UNION ALL
        SELECT 'organization_member', om.id FROM organization_members om WHERE om.organization_id = $1
        UNION ALL
        SELECT 'organization_invite', oi.id FROM organization_invites oi WHERE oi.organization_id = $1
        ON CONFLICT (entity_type, entity_id) DO NOTHING
    `, orgID)
	if err != nil {
		return fmt.Errorf("burying organization: %w", err)
	}
	return nil
}

func (t *tombstones) OrganizationIsBuried(ctx context.Context, orgID uuid.UUID) (bool, error) {
	var exists bool
	err := t.db.QueryRow(ctx, `
        SELECT EXISTS(
            SELECT 1 FROM tombstones WHERE entity_type = 'organization' AND entity_id = $1
        )
    `, orgID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("checking if organization is buried: %w", err)
	}
	return exists, nil
}

func (t *tombstones) BuryMember(ctx context.Context, memberID uuid.UUID) error {
	_, err := t.db.Exec(ctx, `
        INSERT INTO tombstones (entity_type, entity_id)
        VALUES ('organization_member', $1)
        ON CONFLICT (entity_type, entity_id) DO NOTHING
    `, memberID)
	if err != nil {
		return fmt.Errorf("burying member: %w", err)
	}
	return nil
}

func (t *tombstones) MemberIsBuried(ctx context.Context, memberID uuid.UUID) (bool, error) {
	var exists bool
	err := t.db.QueryRow(ctx, `
        SELECT EXISTS(
            SELECT 1 FROM tombstones WHERE entity_type = 'organization_member' AND entity_id = $1
        )
    `, memberID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("checking if member is buried: %w", err)
	}
	return exists, nil
}

func (t *tombstones) BuryInvite(ctx context.Context, inviteID uuid.UUID) error {
	_, err := t.db.Exec(ctx, `
        INSERT INTO tombstones (entity_type, entity_id)
        VALUES ('organization_invite', $1)
        ON CONFLICT (entity_type, entity_id) DO NOTHING
    `, inviteID)
	if err != nil {
		return fmt.Errorf("burying invite: %w", err)
	}
	return nil
}

func (t *tombstones) InviteIsBuried(ctx context.Context, inviteID uuid.UUID) (bool, error) {
	var exists bool
	err := t.db.QueryRow(ctx, `
        SELECT EXISTS(
            SELECT 1 FROM tombstones WHERE entity_type = 'organization_invite' AND entity_id = $1
        )
    `, inviteID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("checking if invite is buried: %w", err)
	}
	return exists, nil
}

func (t *tombstones) BuryPlace(ctx context.Context, placeID uuid.UUID) error {
	_, err := t.db.Exec(ctx, `
        INSERT INTO tombstones (entity_type, entity_id)
        VALUES ('place', $1)
        ON CONFLICT (entity_type, entity_id) DO NOTHING
    `, placeID)
	if err != nil {
		return fmt.Errorf("burying place: %w", err)
	}
	return nil
}

func (t *tombstones) PlaceIsBuried(ctx context.Context, placeID uuid.UUID) (bool, error) {
	var exists bool
	err := t.db.QueryRow(ctx, `
        SELECT EXISTS(
            SELECT 1 FROM tombstones WHERE entity_type = 'place' AND entity_id = $1
        )
    `, placeID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("checking if place is buried: %w", err)
	}
	return exists, nil
}
