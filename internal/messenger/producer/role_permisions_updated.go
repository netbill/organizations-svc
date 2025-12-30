package producer

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"github.com/umisto/cities-svc/internal/domain/models"
	"github.com/umisto/cities-svc/internal/messenger/contracts"
)

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
				{Key: "Producer", Value: []byte(CitiesSvcProducer)},
				{Key: "ContentType", Value: []byte("application/json")},
			},
		},
	)

	return err
}
