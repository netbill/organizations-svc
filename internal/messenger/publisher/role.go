package publisher

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/netbill/eventbox"
	"github.com/netbill/evtypes"
	"github.com/netbill/organizations-svc/internal/core/models"
)

func (p *Publisher) WriteOrgRoleCreated(
	ctx context.Context,
	role models.Role,
) error {
	payload, err := json.Marshal(evtypes.OrgRoleCreatedPayload{
		RoleID:         role.ID,
		OrganizationID: role.OrganizationID,

		Name:        role.Name,
		Description: role.Description,
		Color:       role.Color,
		CreatedAt:   role.CreatedAt,
		Version:     role.Version,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal org role created payload, cause: %w", err)
	}

	_, err = p.outbox.WriteOutboxEvent(ctx, eventbox.Message{
		ID:       uuid.New(),
		Type:     evtypes.OrgRoleCreatedEvent,
		Version:  1,
		Topic:    evtypes.OrganizationsTopicV1,
		Key:      role.OrganizationID.String(),
		Payload:  payload,
		Producer: p.identity,
	})
	if err != nil {
		return fmt.Errorf("failed to create sender event for org role created, cause: %w", err)
	}

	return nil
}

func (p *Publisher) WriteOrgRoleUpdated(
	ctx context.Context,
	role models.Role,
) error {
	payload, err := json.Marshal(evtypes.OrgRoleUpdatedPayload{
		RoleID:         role.ID,
		OrganizationID: role.OrganizationID,
		Name:           role.Name,
		Description:    role.Description,
		Color:          role.Color,
		Version:        role.Version,
		UpdatedAt:      role.UpdatedAt,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal org role updated payload, cause: %w", err)
	}

	_, err = p.outbox.WriteOutboxEvent(ctx, eventbox.Message{
		ID:       uuid.New(),
		Type:     evtypes.OrgRoleUpdatedEvent,
		Version:  1,
		Topic:    evtypes.OrganizationsTopicV1,
		Key:      role.OrganizationID.String(),
		Payload:  payload,
		Producer: p.identity,
	})
	if err != nil {
		return fmt.Errorf("failed to create sender event for org role updated, cause: %w", err)
	}

	return nil
}

func (p *Publisher) WriteOrgRoleDeleted(
	ctx context.Context,
	role models.Role,
) error {
	payload, err := json.Marshal(evtypes.OrgRoleDeletedPayload{
		RoleID:         role.ID,
		OrganizationID: role.OrganizationID,
		DeletedAt:      time.Now().UTC(),
	})
	if err != nil {
		return fmt.Errorf("failed to marshal org role deleted payload, cause: %w", err)
	}

	_, err = p.outbox.WriteOutboxEvent(ctx, eventbox.Message{
		ID:       uuid.New(),
		Type:     evtypes.OrgRoleDeletedEvent,
		Version:  1,
		Topic:    evtypes.OrganizationsTopicV1,
		Key:      role.OrganizationID.String(),
		Payload:  payload,
		Producer: p.identity,
	})
	if err != nil {
		return fmt.Errorf("failed to create sender event for org role deleted, cause: %w", err)
	}

	return nil
}

func (p *Publisher) WriteOrgRolePermissionsUpdated(
	ctx context.Context,
	role models.Role,
	permissions []uuid.UUID,
	revision models.OrgRolePermissionsLinksRevision,
) error {
	payload, err := json.Marshal(evtypes.OrgRolePermissionsUpdatedPayload{
		RoleID:      role.ID,
		Permissions: permissions,

		PermissionsRevision:          revision.Revision,
		PermissionsRevisionUpdatedAt: revision.UpdatedAt,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal org role permissions updated payload, cause: %w", err)
	}

	_, err = p.outbox.WriteOutboxEvent(ctx, eventbox.Message{
		ID:       uuid.New(),
		Type:     evtypes.OrgRolePermissionsUpdatedEvent,
		Version:  1,
		Topic:    evtypes.OrganizationsTopicV1,
		Key:      role.OrganizationID.String(),
		Payload:  payload,
		Producer: p.identity,
	})
	if err != nil {
		return fmt.Errorf("failed to create sender event for org role permissions updated, cause: %w", err)
	}

	return nil
}

func (p *Publisher) WriteOrgRolesRanksUpdated(
	ctx context.Context,
	organizationID uuid.UUID,
	ranks []models.OrgRoleRank,
	revision models.OrgRoleRanksRevision,
) error {
	items := make([]evtypes.RankSnapshotItem, len(ranks))
	for i, r := range ranks {
		items[i] = evtypes.RankSnapshotItem{
			RoleID: r.RoleID,
			Rank:   r.Rank,
		}
	}

	payload, err := json.Marshal(evtypes.OrgRolesRanksUpdatedPayload{
		OrganizationID:         organizationID,
		Ranks:                  items,
		RanksRevision:          revision.Revision,
		RanksRevisionUpdatedAt: revision.UpdatedAt,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal org role ranks updated payload, cause: %w", err)
	}

	_, err = p.outbox.WriteOutboxEvent(ctx, eventbox.Message{
		ID:       uuid.New(),
		Type:     evtypes.OrgRolesRanksUpdatedEvent,
		Version:  1,
		Topic:    evtypes.OrganizationsTopicV1,
		Key:      organizationID.String(),
		Payload:  payload,
		Producer: p.identity,
	})
	if err != nil {
		return fmt.Errorf("failed to create sender event for org roles ranks updated, cause: %w", err)
	}

	return nil
}

func (p *Publisher) WriteOrgMemberRolesUpdated(
	ctx context.Context,
	memberID uuid.UUID,
	roles []uuid.UUID,
	revision models.OrgMemberRoleLinkRevision,
) error {
	payload, err := json.Marshal(evtypes.OrgMemberRolesUpdatedPayload{
		MemberID: memberID,
		Roles:    roles,

		MemberRolesRevision:          revision.Revision,
		MemberRolesRevisionUpdatedAt: revision.UpdatedAt,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal org member role updated payload, cause: %w", err)
	}

	_, err = p.outbox.WriteOutboxEvent(ctx, eventbox.Message{
		ID:       uuid.New(),
		Type:     evtypes.OrgMemberRolesUpdatedEvent,
		Version:  1,
		Topic:    evtypes.OrgMembersTopicV1,
		Key:      memberID.String(),
		Payload:  payload,
		Producer: p.identity,
	})
	if err != nil {
		return fmt.Errorf("failed to create sender event for org member role updated, cause: %w", err)
	}

	return nil
}
