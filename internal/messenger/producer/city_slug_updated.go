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

func (s Service) WriteCitySlugUpdated(
	ctx context.Context,
	city models.City,
	oldSlug *string,
) error {
	payload, err := json.Marshal(contracts.CitySlugUpdatedPayload{
		City:    city,
		OldSlug: oldSlug,
	})
	if err != nil {
		return err
	}

	_, err = s.outbox.CreateOutboxEvent(
		ctx,
		box.OutboxStatusPending,
		kafka.Message{
			Topic: contracts.CitiesTopicV1,
			Key:   []byte(city.ID.String()),
			Value: payload,
			Headers: []kafka.Header{
				{Key: header.EventID, Value: []byte(uuid.New().String())}, // Outbox will fill this
				{Key: header.EventType, Value: []byte(contracts.CitySlugUpdatedEvent)},
				{Key: header.EventVersion, Value: []byte("1")},
				{Key: header.Producer, Value: []byte(contracts.CitiesSvcGroup)},
				{Key: header.ContentType, Value: []byte("application/json")},
			},
		},
	)

	return err
}
