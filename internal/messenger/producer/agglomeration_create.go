package producer

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"github.com/umisto/cities-svc/internal/domain/models"
	"github.com/umisto/cities-svc/internal/messenger/contracts"
	"github.com/umisto/kafkakit/box"
	"github.com/umisto/kafkakit/header"
)

func (s Service) WriteAgglomerationCreated(
	ctx context.Context,
	agglomeration models.Agglomeration,
) error {
	payload, err := json.Marshal(contracts.AgglomerationCreatedPayload{
		Agglomeration: agglomeration,
	})
	if err != nil {
		return err
	}

	eventID := uuid.New()

	_, err = s.outbox.CreateOutboxEvent(
		ctx,
		box.OutboxStatusPending,
		kafka.Message{
			Topic: contracts.AgglomerationsTopicV1,
			Key:   []byte(agglomeration.ID.String()),
			Value: payload,
			Headers: []kafka.Header{
				{Key: header.EventID, Value: []byte(eventID.String())}, // Outbox will fill this
				{Key: header.EventType, Value: []byte(contracts.AgglomerationCreatedEvent)},
				{Key: header.EventVersion, Value: []byte("1")},
				{Key: header.Producer, Value: []byte(CitiesSvcProducer)},
				{Key: header.ContentType, Value: []byte("application/json")},
			},
		},
	)

	return err
}
