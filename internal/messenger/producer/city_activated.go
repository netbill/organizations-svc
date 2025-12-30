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

func (s Service) WriteCityActivated(
	ctx context.Context,
	city models.City,
) error {
	payload, err := json.Marshal(contracts.CityActivatedPayload{
		City: city,
	})
	if err != nil {
		return err
	}

	eventID := uuid.New()

	_, err = s.outbox.CreateOutboxEvent(
		ctx,
		box.OutboxStatusPending,
		kafka.Message{
			Topic: contracts.CitiesTopicV1,
			Key:   []byte(city.ID.String()),
			Value: payload,
			Headers: []kafka.Header{
				{Key: header.EventID, Value: []byte(eventID.String())}, // Outbox will fill this
				{Key: header.EventType, Value: []byte(contracts.CityActivatedEvent)},
				{Key: header.EventVersion, Value: []byte("1")},
				{Key: header.Producer, Value: []byte(CitiesSvcProducer)},
				{Key: header.ContentType, Value: []byte("application/json")},
			},
		},
	)

	return err
}
