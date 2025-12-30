package producer

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"github.com/umisto/cities-svc/internal/domain/models"
	"github.com/umisto/cities-svc/internal/messenger/contracts"
)

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
				{Key: "Producer", Value: []byte(CitiesSvcProducer)},
				{Key: "ContentType", Value: []byte("application/json")},
			},
		},
	)

	return err
}
