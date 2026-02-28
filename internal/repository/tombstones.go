package repository

import (
	"context"

	"github.com/google/uuid"
)

type TombstonesSql interface {
	BuryProfile(ctx context.Context, accountID uuid.UUID) error
	ProfileIsBuried(ctx context.Context, accountID uuid.UUID) (bool, error)

	BuryPlace(ctx context.Context, placeID uuid.UUID) error
	PlaceIsBuried(ctx context.Context, placeID uuid.UUID) (bool, error)

	BuryOrganization(ctx context.Context, orgID uuid.UUID) error
	OrganizationIsBuried(ctx context.Context, orgID uuid.UUID) (bool, error)
	BuryMember(ctx context.Context, accountID uuid.UUID) error
	MemberIsBuried(ctx context.Context, accountID uuid.UUID) (bool, error)
	BuryInvite(ctx context.Context, inviteID uuid.UUID) error
	InviteIsBuried(ctx context.Context, inviteID uuid.UUID) (bool, error)
}
