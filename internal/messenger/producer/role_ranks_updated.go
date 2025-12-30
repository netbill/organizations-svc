package producer

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"github.com/umisto/cities-svc/internal/messenger/contracts"
)

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
				{Key: "Producer", Value: []byte(CitiesSvcProducer)},
				{Key: "ContentType", Value: []byte("application/json")},
			},
		},
	)

	return err
}
