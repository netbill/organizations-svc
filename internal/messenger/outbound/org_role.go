package outbound

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/netbill/evebox/header"
	"github.com/netbill/organizations-svc/internal/core/models"
	"github.com/netbill/organizations-svc/internal/messenger/contracts"
	"github.com/segmentio/kafka-go"
)

func (o *Outbound) WriteOrgRoleCreated(
	ctx context.Context,
	role models.Role,
) error {
	payload, err := json.Marshal(contracts.OrgRoleCreatedPayload{
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

	_, err = o.outbox.CreateOutboxEvent(
		ctx,
		kafka.Message{
			Topic: contracts.OrganizationsTopicV1,
			Key:   []byte(role.OrganizationID.String()),
			Value: payload,
			Headers: []kafka.Header{
				{Key: header.EventID, Value: []byte(uuid.New().String())},
				{Key: header.EventType, Value: []byte(contracts.OrgRoleCreatedEvent)},
				{Key: header.EventVersion, Value: []byte("1")},
				{Key: header.Producer, Value: []byte(contracts.OrganizationsSvcGroup)},
				{Key: header.ContentType, Value: []byte("application/json")},
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create outbox event for org role created, cause: %w", err)
	}

	return nil
}

func (o *Outbound) WriteOrgRoleUpdated(
	ctx context.Context,
	role models.Role,
) error {
	payload, err := json.Marshal(contracts.OrgRoleUpdatedPayload{
		RoleID:      role.ID,
		Name:        role.Name,
		Description: role.Description,
		Color:       role.Color,
		UpdatedAt:   role.UpdatedAt,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal org role updated payload, cause: %w", err)
	}

	_, err = o.outbox.CreateOutboxEvent(
		ctx,
		kafka.Message{
			Topic: contracts.OrganizationsTopicV1,
			Key:   []byte(role.OrganizationID.String()),
			Value: payload,
			Headers: []kafka.Header{
				{Key: header.EventID, Value: []byte(uuid.New().String())},
				{Key: header.EventType, Value: []byte(contracts.OrgRoleUpdatedEvent)},
				{Key: header.EventVersion, Value: []byte("1")},
				{Key: header.Producer, Value: []byte(contracts.OrganizationsSvcGroup)},
				{Key: header.ContentType, Value: []byte("application/json")},
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create outbox event for org role updated, cause: %w", err)
	}

	return nil
}

func (o *Outbound) WriteOrgRoleDeleted(
	ctx context.Context,
	role models.Role,
) error {
	payload, err := json.Marshal(contracts.OrgRoleDeletedPayload{
		RoleID:    role.ID,
		DeletedAt: time.Now().UTC(),
	})
	if err != nil {
		return fmt.Errorf("failed to marshal org role deleted payload, cause: %w", err)
	}

	_, err = o.outbox.CreateOutboxEvent(
		ctx,
		kafka.Message{
			Topic: contracts.OrganizationsTopicV1,
			Key:   []byte(role.OrganizationID.String()),
			Value: payload,
			Headers: []kafka.Header{
				{Key: header.EventID, Value: []byte(uuid.New().String())},
				{Key: header.EventType, Value: []byte(contracts.OrgRoleDeletedEvent)},
				{Key: header.EventVersion, Value: []byte("1")},
				{Key: header.Producer, Value: []byte(contracts.OrganizationsSvcGroup)},
				{Key: header.ContentType, Value: []byte("application/json")},
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create outbox event for org role deleted, cause: %w", err)
	}

	return nil
}

func (o *Outbound) WriteOrgRolePermissionsUpdated(
	ctx context.Context,
	role models.Role,
	permissions models.OrgRolePermissionAccess,
) error {
	payload, err := json.Marshal(contracts.OrgRolePermissionsUpdatedPayload{
		RoleID:      role.ID,
		Permissions: permissions.ToMap(),
	})
	if err != nil {
		return fmt.Errorf("failed to marshal org role permissions updated payload, cause: %w", err)
	}

	_, err = o.outbox.CreateOutboxEvent(
		ctx,
		kafka.Message{
			Topic: contracts.OrganizationsTopicV1,
			Key:   []byte(role.OrganizationID.String()),
			Value: payload,
			Headers: []kafka.Header{
				{Key: header.EventID, Value: []byte(uuid.New().String())},
				{Key: header.EventType, Value: []byte(contracts.OrgRolePermissionsUpdatedEvent)},
				{Key: header.EventVersion, Value: []byte("1")},
				{Key: header.Producer, Value: []byte(contracts.OrganizationsSvcGroup)},
				{Key: header.ContentType, Value: []byte("application/json")},
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create outbox event for org role permissions updated, cause: %w", err)
	}

	return nil
}

func (o *Outbound) WriteOrgRolesRanksUpdated(
	ctx context.Context,
	organizationID uuid.UUID,
	ranks map[uuid.UUID]uint,
) error {
	payload, err := json.Marshal(contracts.OrgRolesRanksUpdatedPayload{
		OrganizationID: organizationID,
		Ranks:          ranks,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal org role ranks updated payload, cause: %w", err)
	}

	_, err = o.outbox.CreateOutboxEvent(
		ctx,
		kafka.Message{
			Topic: contracts.OrganizationsTopicV1,
			Key:   []byte(organizationID.String()),
			Value: payload,
			Headers: []kafka.Header{
				{Key: header.EventID, Value: []byte(uuid.New().String())},
				{Key: header.EventType, Value: []byte(contracts.OrgRolesRanksUpdatedEvent)},
				{Key: header.EventVersion, Value: []byte("1")},
				{Key: header.Producer, Value: []byte(contracts.OrganizationsSvcGroup)},
				{Key: header.ContentType, Value: []byte("application/json")},
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create outbox event for org roles ranks updated, cause: %w", err)
	}

	return nil
}
