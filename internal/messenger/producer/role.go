package producer

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/netbill/kafkakit/header"
	"github.com/netbill/organizations-svc/internal/core/models"
	"github.com/netbill/organizations-svc/internal/messenger/contracts"
	"github.com/segmentio/kafka-go"
)

func (p Producer) WriteRoleCreated(
	ctx context.Context,
	role models.Role,
) error {
	payload, err := json.Marshal(contracts.RoleCreatedPayload{
		Role: role,
	})
	if err != nil {
		return err
	}

	_, err = p.outbox.CreateOutboxEvent(
		ctx,
		kafka.Message{
			Topic: contracts.RolesTopicV1,
			Key:   []byte(role.ID.String()),
			Value: payload,
			Headers: []kafka.Header{
				{Key: header.EventID, Value: []byte(uuid.New().String())},
				{Key: header.EventType, Value: []byte(contracts.RoleCreatedEvent)},
				{Key: header.EventVersion, Value: []byte("1")},
				{Key: header.Producer, Value: []byte(contracts.OrganizationsSvcGroup)},
				{Key: header.ContentType, Value: []byte("application/json")},
			},
		},
	)

	return err
}

func (p Producer) WriteRoleUpdated(
	ctx context.Context,
	rol models.Role,
) error {
	payload, err := json.Marshal(contracts.RoleUpdatedPayload{
		Role: rol,
	})
	if err != nil {
		return err
	}

	_, err = p.outbox.CreateOutboxEvent(
		ctx,
		kafka.Message{
			Topic: contracts.RolesTopicV1,
			Key:   []byte(rol.ID.String()),
			Value: payload,
			Headers: []kafka.Header{
				{Key: header.EventID, Value: []byte(uuid.New().String())},
				{Key: header.EventType, Value: []byte(contracts.RoleUpdatedEvent)},
				{Key: header.EventVersion, Value: []byte("1")},
				{Key: header.Producer, Value: []byte(contracts.OrganizationsSvcGroup)},
				{Key: header.ContentType, Value: []byte("application/json")},
			},
		},
	)

	return err
}

func (p Producer) WriteRolePermissionsUpdated(
	ctx context.Context,
	RoleID uuid.UUID,
	permissions map[models.Permission]bool,
) error {
	payload, err := json.Marshal(contracts.RolePermissionsUpdatedPayload{
		RoleID:      RoleID,
		Permissions: permissions,
	})
	if err != nil {
		return err
	}

	_, err = p.outbox.CreateOutboxEvent(
		ctx,
		kafka.Message{
			Topic: contracts.RolesTopicV1,
			Key:   []byte(RoleID.String()),
			Value: payload,
			Headers: []kafka.Header{
				{Key: header.EventID, Value: []byte(uuid.New().String())},
				{Key: header.EventType, Value: []byte(contracts.RolePermissionsUpdatedEvent)},
				{Key: header.EventVersion, Value: []byte("1")},
				{Key: header.Producer, Value: []byte(contracts.OrganizationsSvcGroup)},
				{Key: header.ContentType, Value: []byte("application/json")},
			},
		},
	)

	return err
}

func (p Producer) WriteRolesRanksUpdated(
	ctx context.Context,
	organizationID uuid.UUID,
	ranks map[uuid.UUID]uint,
) error {
	payload, err := json.Marshal(contracts.RolesRanksUpdatedPayload{
		OrganizationID: organizationID,
		Ranks:          ranks,
	})
	if err != nil {
		return err
	}

	_, err = p.outbox.CreateOutboxEvent(
		ctx,
		kafka.Message{
			Topic: contracts.RolesTopicV1,
			Key:   []byte(organizationID.String()),
			Value: payload,
			Headers: []kafka.Header{
				{Key: header.EventID, Value: []byte(uuid.New().String())},
				{Key: header.EventType, Value: []byte(contracts.RolesRanksUpdatedEvent)},
				{Key: header.EventVersion, Value: []byte("1")},
				{Key: header.Producer, Value: []byte(contracts.OrganizationsSvcGroup)},
				{Key: header.ContentType, Value: []byte("application/json")},
			},
		},
	)

	return err
}

func (p Producer) WriteRoleDeleted(
	ctx context.Context,
	role models.Role,
) error {
	payload, err := json.Marshal(contracts.RoleDeletedPayload{
		Role: role,
	})
	if err != nil {
		return err
	}

	_, err = p.outbox.CreateOutboxEvent(
		ctx,
		kafka.Message{
			Topic: contracts.RolesTopicV1,
			Key:   []byte(role.ID.String()),
			Value: payload,
			Headers: []kafka.Header{
				{Key: header.EventID, Value: []byte(uuid.New().String())},
				{Key: header.EventType, Value: []byte(contracts.RoleDeletedEvent)},
				{Key: header.EventVersion, Value: []byte("1")},
				{Key: header.Producer, Value: []byte(contracts.OrganizationsSvcGroup)},
				{Key: header.ContentType, Value: []byte("application/json")},
			},
		},
	)

	return err
}
