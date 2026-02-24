package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/netbill/organizations-svc/internal/core/models"
)

func (r *Repository) AddMemberRole(ctx context.Context, memberID, roleID uuid.UUID) ([]uuid.UUID, error) {
	//TODO implement me
	panic("implement me")
}

func (r *Repository) RemoveMemberRole(ctx context.Context, memberID, roleID uuid.UUID) ([]uuid.UUID, error) {
	//TODO implement me
	panic("implement me")
}

func (r *Repository) GetMemberMaxRoleRank(ctx context.Context, memberID uuid.UUID) (int32, error) {
	//TODO implement me
	panic("implement me")
}

func (r *Repository) CreateMemberRolesLinksRevision(ctx context.Context, memberID uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (r *Repository) LockMemberRolesLinksRevision(ctx context.Context, memberID uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (r *Repository) BumpMemberRolesLinksRevision(ctx context.Context, memberID uuid.UUID) (models.OrgMemberRoleLinkRevision, error) {
	//TODO implement me
	panic("implement me")
}
