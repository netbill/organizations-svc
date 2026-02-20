package outbound

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/netbill/eventbox"
	"github.com/netbill/organizations-svc/internal/core/models"
	"github.com/netbill/organizations-svc/internal/core/modules/role"
	"github.com/netbill/organizations-svc/pkg/evtypes"
)

func (s *Sender) WriteOrgRoleCreated(
	ctx context.Context,
	role models.Role,
) error {
	payload, err := json.Marshal(evtypes.OrgRoleCreatedPayload{
		RoleID:         role.ID,
		OrganizationID: role.OrganizationID,
		Rank:           role.Rank,
		Name:           role.Name,
		Description:    role.Description,
		Color:          role.Color,
		CreatedAt:      role.CreatedAt,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal org role created payload, cause: %w", err)
	}

	_, err = s.outbox.WriteOutboxEvent(
		ctx,
		eventbox.Event{
			ID:       uuid.New(),
			Type:     evtypes.OrgRoleCreatedEvent,
			Version:  1,
			Topic:    evtypes.OrganizationsTopicV1,
			Key:      role.OrganizationID.String(),
			Payload:  payload,
			Producer: s.identity,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create sender event for org role created, cause: %w", err)
	}

	return nil
}

func (s *Sender) WriteOrgRoleUpdated(
	ctx context.Context,
	role models.Role,
) error {
	payload, err := json.Marshal(evtypes.OrgRoleUpdatedPayload{
		RoleID:      role.ID,
		Name:        role.Name,
		Description: role.Description,
		Color:       role.Color,
		UpdatedAt:   role.UpdatedAt,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal org role updated payload, cause: %w", err)
	}

	_, err = s.outbox.WriteOutboxEvent(
		ctx,
		eventbox.Event{
			ID:       uuid.New(),
			Type:     evtypes.OrgRoleUpdatedEvent,
			Version:  1,
			Topic:    evtypes.OrganizationsTopicV1,
			Key:      role.OrganizationID.String(),
			Payload:  payload,
			Producer: s.identity,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create sender event for org role updated, cause: %w", err)
	}

	return nil
}

func (s *Sender) WriteOrgRoleDeleted(
	ctx context.Context,
	role models.Role,
) error {
	payload, err := json.Marshal(evtypes.OrgRoleDeletedPayload{
		RoleID:    role.ID,
		DeletedAt: time.Now().UTC(),
	})
	if err != nil {
		return fmt.Errorf("failed to marshal org role deleted payload, cause: %w", err)
	}

	_, err = s.outbox.WriteOutboxEvent(
		ctx,
		eventbox.Event{
			ID:       uuid.New(),
			Type:     evtypes.OrgRoleDeletedEvent,
			Version:  1,
			Topic:    evtypes.OrganizationsTopicV1,
			Key:      role.OrganizationID.String(),
			Payload:  payload,
			Producer: s.identity,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create sender event for org role deleted, cause: %w", err)
	}

	return nil
}

func (s *Sender) WriteOrgRolePermissionsUpdated(
	ctx context.Context,
	role models.Role,
	permissions role.SetPermissions,
) error {
	payload, err := json.Marshal(evtypes.OrgRolePermissionsUpdatedPayload{
		RoleID:      role.ID,
		Permissions: permissions,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal org role permissions updated payload, cause: %w", err)
	}

	_, err = s.outbox.WriteOutboxEvent(
		ctx,
		eventbox.Event{
			ID:       uuid.New(),
			Type:     evtypes.OrgRolePermissionsUpdatedEvent,
			Version:  1,
			Topic:    evtypes.OrganizationsTopicV1,
			Key:      role.OrganizationID.String(),
			Payload:  payload,
			Producer: s.identity,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create sender event for org role permissions updated, cause: %w", err)
	}

	return nil
}

func (s *Sender) WriteOrgRolesRanksUpdated(
	ctx context.Context,
	organizationID uuid.UUID,
	ranks map[uuid.UUID]uint,
) error {
	payload, err := json.Marshal(evtypes.OrgRolesRanksUpdatedPayload{
		OrganizationID: organizationID,
		Ranks:          ranks,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal org role ranks updated payload, cause: %w", err)
	}

	_, err = s.outbox.WriteOutboxEvent(
		ctx,
		eventbox.Event{
			ID:       uuid.New(),
			Type:     evtypes.OrgRolesRanksUpdatedEvent,
			Version:  1,
			Topic:    evtypes.OrganizationsTopicV1,
			Key:      organizationID.String(),
			Payload:  payload,
			Producer: s.identity,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create sender event for org roles ranks updated, cause: %w", err)
	}

	return nil
}
