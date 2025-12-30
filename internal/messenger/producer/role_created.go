package producer

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"github.com/umisto/cities-svc/internal/domain/models"
	"github.com/umisto/cities-svc/internal/messenger/contracts"
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
				{Key: "Producer", Value: []byte(CitiesSvcProducer)},
				{Key: "ContentType", Value: []byte("application/json")},
			},
		},
	)

	return err
}
