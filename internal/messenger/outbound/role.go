package outbound

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/netbill/eventbox/headers"
	"github.com/netbill/organizations-svc/internal/core/models"
	"github.com/netbill/organizations-svc/internal/core/modules/role"
	"github.com/netbill/organizations-svc/pkg/evtypes"
	"github.com/segmentio/kafka-go"
)

func (o *Outbound) WriteOrgRoleCreated(
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

	_, err = o.outbox.WriteToOutbox(
		ctx,
		kafka.Message{
			Topic: evtypes.OrganizationsTopicV1,
			Key:   []byte(role.OrganizationID.String()),
			Value: payload,
			Headers: []kafka.Header{
				{Key: headers.EventID, Value: []byte(uuid.New().String())},
				{Key: headers.EventType, Value: []byte(evtypes.OrgRoleCreatedEvent)},
				{Key: headers.EventVersion, Value: []byte("1")},
				{Key: headers.Producer, Value: []byte(o.groupID)},
				{Key: headers.ContentType, Value: []byte("application/json")},
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

	_, err = o.outbox.WriteToOutbox(
		ctx,
		kafka.Message{
			Topic: evtypes.OrganizationsTopicV1,
			Key:   []byte(role.OrganizationID.String()),
			Value: payload,
			Headers: []kafka.Header{
				{Key: headers.EventID, Value: []byte(uuid.New().String())},
				{Key: headers.EventType, Value: []byte(evtypes.OrgRoleUpdatedEvent)},
				{Key: headers.EventVersion, Value: []byte("1")},
				{Key: headers.Producer, Value: []byte(o.groupID)},
				{Key: headers.ContentType, Value: []byte("application/json")},
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
	payload, err := json.Marshal(evtypes.OrgRoleDeletedPayload{
		RoleID:    role.ID,
		DeletedAt: time.Now().UTC(),
	})
	if err != nil {
		return fmt.Errorf("failed to marshal org role deleted payload, cause: %w", err)
	}

	_, err = o.outbox.WriteToOutbox(
		ctx,
		kafka.Message{
			Topic: evtypes.OrganizationsTopicV1,
			Key:   []byte(role.OrganizationID.String()),
			Value: payload,
			Headers: []kafka.Header{
				{Key: headers.EventID, Value: []byte(uuid.New().String())},
				{Key: headers.EventType, Value: []byte(evtypes.OrgRoleDeletedEvent)},
				{Key: headers.EventVersion, Value: []byte("1")},
				{Key: headers.Producer, Value: []byte(o.groupID)},
				{Key: headers.ContentType, Value: []byte("application/json")},
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
	permissions role.SetForRole,
) error {
	payload, err := json.Marshal(evtypes.OrgRolePermissionsUpdatedPayload{
		RoleID:      role.ID,
		Permissions: permissions,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal org role permissions updated payload, cause: %w", err)
	}

	_, err = o.outbox.WriteToOutbox(
		ctx,
		kafka.Message{
			Topic: evtypes.OrganizationsTopicV1,
			Key:   []byte(role.OrganizationID.String()),
			Value: payload,
			Headers: []kafka.Header{
				{Key: headers.EventID, Value: []byte(uuid.New().String())},
				{Key: headers.EventType, Value: []byte(evtypes.OrgRolePermissionsUpdatedEvent)},
				{Key: headers.EventVersion, Value: []byte("1")},
				{Key: headers.Producer, Value: []byte(o.groupID)},
				{Key: headers.ContentType, Value: []byte("application/json")},
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
	payload, err := json.Marshal(evtypes.OrgRolesRanksUpdatedPayload{
		OrganizationID: organizationID,
		Ranks:          ranks,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal org role ranks updated payload, cause: %w", err)
	}

	_, err = o.outbox.WriteToOutbox(
		ctx,
		kafka.Message{
			Topic: evtypes.OrganizationsTopicV1,
			Key:   []byte(organizationID.String()),
			Value: payload,
			Headers: []kafka.Header{
				{Key: headers.EventID, Value: []byte(uuid.New().String())},
				{Key: headers.EventType, Value: []byte(evtypes.OrgRolesRanksUpdatedEvent)},
				{Key: headers.EventVersion, Value: []byte("1")},
				{Key: headers.Producer, Value: []byte(o.groupID)},
				{Key: headers.ContentType, Value: []byte("application/json")},
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create outbox event for org roles ranks updated, cause: %w", err)
	}

	return nil
}
