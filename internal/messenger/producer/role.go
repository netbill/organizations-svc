package producer

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"github.com/umisto/agglomerations-svc/internal/domain/models"
	"github.com/umisto/agglomerations-svc/internal/messenger/contracts"
	"github.com/umisto/kafkakit/box"
)

func (s Service) WriteRoleCreated(
	ctx context.Context,
	role models.Role,
) error {
	payload, err := json.Marshal(contracts.RoleCreatedPayload{
		Role: role,
	})
	if err != nil {
		return err
	}

	_, err = s.outbox.CreateOutboxEvent(
		ctx,
		box.OutboxStatusPending,
		kafka.Message{
			Topic: contracts.RolesTopicV1,
			Key:   []byte(role.ID.String()),
			Value: payload,
			Headers: []kafka.Header{
				{Key: "EventID", Value: []byte(uuid.New().String())}, // Outbox will fill this
				{Key: "EventType", Value: []byte(contracts.RoleCreatedEvent)},
				{Key: "EventVersion", Value: []byte("1")},
				{Key: "Producer", Value: []byte(contracts.AgglomerationsSvcGroup)},
				{Key: "ContentType", Value: []byte("application/json")},
			},
		},
	)

	return err
}

func (s Service) WriteRoleUpdated(
	ctx context.Context,
	rol models.Role,
) error {
	payload, err := json.Marshal(contracts.RoleUpdatedPayload{
		Role: rol,
	})
	if err != nil {
		return err
	}

	_, err = s.outbox.CreateOutboxEvent(
		ctx,
		contracts.RolesTopicV1,
		kafka.Message{
			Topic: contracts.RolesTopicV1,
			Key:   []byte(rol.ID.String()),
			Value: payload,
			Headers: []kafka.Header{
				{Key: "EventID", Value: []byte(uuid.New().String())}, // Outbox will fill this
				{Key: "EventType", Value: []byte(contracts.RoleUpdatedEvent)},
				{Key: "EventVersion", Value: []byte("1")},
				{Key: "Producer", Value: []byte(contracts.AgglomerationsSvcGroup)},
				{Key: "ContentType", Value: []byte("application/json")},
			},
		},
	)

	return err
}

func (s Service) WriteRolePermissionsUpdated(
	ctx context.Context,
	RoleID uuid.UUID,
	permissions map[models.CodeRolePermission]bool,
) error {
	payload, err := json.Marshal(contracts.RolePermissionsUpdatedPayload{
		RoleID:      RoleID,
		Permissions: permissions,
	})
	if err != nil {
		return err
	}

	_, err = s.outbox.CreateOutboxEvent(
		ctx,
		contracts.RolesTopicV1,
		kafka.Message{
			Topic: contracts.RolesTopicV1,
			Key:   []byte(RoleID.String()),
			Value: payload,
			Headers: []kafka.Header{
				{Key: "EventID", Value: []byte(uuid.New().String())}, // Outbox will fill this
				{Key: "EventType", Value: []byte(contracts.RolePermissionsUpdatedEvent)},
				{Key: "EventVersion", Value: []byte("1")},
				{Key: "Producer", Value: []byte(contracts.AgglomerationsSvcGroup)},
				{Key: "ContentType", Value: []byte("application/json")},
			},
		},
	)

	return err
}

func (s Service) WriteRolesRanksUpdated(
	ctx context.Context,
	agglomerationID uuid.UUID,
	ranks map[uuid.UUID]uint,
) error {
	payload, err := json.Marshal(contracts.RolesRanksUpdatedPayload{
		AgglomerationID: agglomerationID,
		Ranks:           ranks,
	})
	if err != nil {
		return err
	}

	_, err = s.outbox.CreateOutboxEvent(
		ctx,
		contracts.RolesTopicV1,
		kafka.Message{
			Topic: contracts.RolesTopicV1,
			Key:   []byte(agglomerationID.String()),
			Value: payload,
			Headers: []kafka.Header{
				{Key: "EventID", Value: []byte(uuid.New().String())}, // Outbox will fill this
				{Key: "EventType", Value: []byte(contracts.RolesRanksUpdatedEvent)},
				{Key: "EventVersion", Value: []byte("1")},
				{Key: "Producer", Value: []byte(contracts.AgglomerationsSvcGroup)},
				{Key: "ContentType", Value: []byte("application/json")},
			},
		},
	)

	return err
}

func (s Service) WriteRoleDeleted(
	ctx context.Context,
	role models.Role,
) error {
	payload, err := json.Marshal(contracts.RoleDeletedPayload{
		Role: role,
	})
	if err != nil {
		return err
	}

	_, err = s.outbox.CreateOutboxEvent(
		ctx,
		contracts.RolesTopicV1,
		kafka.Message{
			Topic: contracts.RolesTopicV1,
			Key:   []byte(role.ID.String()),
			Value: payload,
			Headers: []kafka.Header{
				{Key: "EventID", Value: []byte(uuid.New().String())}, // Outbox will fill this
				{Key: "EventType", Value: []byte(contracts.RoleDeletedEvent)},
				{Key: "EventVersion", Value: []byte("1")},
				{Key: "Producer", Value: []byte(contracts.AgglomerationsSvcGroup)},
				{Key: "ContentType", Value: []byte("application/json")},
			},
		},
	)

	return err
}
