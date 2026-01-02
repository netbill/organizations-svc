package producer

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"github.com/umisto/agglomerations-svc/internal/domain/models"
	"github.com/umisto/agglomerations-svc/internal/messenger/contracts"
	"github.com/umisto/kafkakit/header"
)

func (p Producer) WriteAgglomerationCreated(
	ctx context.Context,
	agglomeration models.Agglomeration,
) error {
	payload, err := json.Marshal(contracts.AgglomerationCreatedPayload{
		Agglomeration: agglomeration,
	})
	if err != nil {
		return err
	}

	_, err = p.outbox.CreateOutboxEvent(
		ctx,
		kafka.Message{
			Topic: contracts.AgglomerationsTopicV1,
			Key:   []byte(agglomeration.ID.String()),
			Value: payload,
			Headers: []kafka.Header{
				{Key: header.EventID, Value: []byte(uuid.New().String())},
				{Key: header.EventType, Value: []byte(contracts.AgglomerationCreatedEvent)},
				{Key: header.EventVersion, Value: []byte("1")},
				{Key: header.Producer, Value: []byte(contracts.AgglomerationsSvcGroup)},
				{Key: header.ContentType, Value: []byte("application/json")},
			},
		},
	)

	return err
}

func (p Producer) WriteAgglomerationUpdated(
	ctx context.Context,
	agglomeration models.Agglomeration,
) error {
	payload, err := json.Marshal(contracts.AgglomerationUpdatedPayload{
		Agglomeration: agglomeration,
	})
	if err != nil {
		return err
	}

	_, err = p.outbox.CreateOutboxEvent(
		ctx,
		kafka.Message{
			Topic: contracts.AgglomerationsTopicV1,
			Key:   []byte(agglomeration.ID.String()),
			Value: payload,
			Headers: []kafka.Header{
				{Key: header.EventID, Value: []byte(uuid.New().String())},
				{Key: header.EventType, Value: []byte(contracts.AgglomerationUpdatedEvent)},
				{Key: header.EventVersion, Value: []byte("1")},
				{Key: header.Producer, Value: []byte(contracts.AgglomerationsSvcGroup)},
				{Key: header.ContentType, Value: []byte("application/json")},
			},
		},
	)

	return err
}

func (p Producer) WriteAgglomerationDeleted(
	ctx context.Context,
	agglomeration models.Agglomeration,
) error {
	payload, err := json.Marshal(contracts.AgglomerationDeletedPayload{
		Agglomeration: agglomeration,
	})
	if err != nil {
		return err
	}

	_, err = p.outbox.CreateOutboxEvent(
		ctx,
		kafka.Message{
			Topic: contracts.AgglomerationsTopicV1,
			Key:   []byte(agglomeration.ID.String()),
			Value: payload,
			Headers: []kafka.Header{
				{Key: header.EventID, Value: []byte(uuid.New().String())},
				{Key: header.EventType, Value: []byte(contracts.AgglomerationDeletedEvent)},
				{Key: header.EventVersion, Value: []byte("1")},
				{Key: header.Producer, Value: []byte(contracts.AgglomerationsSvcGroup)},
				{Key: header.ContentType, Value: []byte("application/json")},
			},
		},
	)

	return err
}

func (p Producer) WriteAgglomerationActivated(
	ctx context.Context,
	agglomeration models.Agglomeration,
) error {
	payload, err := json.Marshal(contracts.AgglomerationActivatedPayload{
		Agglomeration: agglomeration,
	})
	if err != nil {
		return err
	}

	_, err = p.outbox.CreateOutboxEvent(
		ctx,
		kafka.Message{
			Topic: contracts.AgglomerationsTopicV1,
			Key:   []byte(agglomeration.ID.String()),
			Value: payload,
			Headers: []kafka.Header{
				{Key: header.EventID, Value: []byte(uuid.New().String())},
				{Key: header.EventType, Value: []byte(contracts.AgglomerationActivatedEvent)},
				{Key: header.EventVersion, Value: []byte("1")},
				{Key: header.Producer, Value: []byte(contracts.AgglomerationsSvcGroup)},
				{Key: header.ContentType, Value: []byte("application/json")},
			},
		},
	)

	return err
}

func (p Producer) WriteAgglomerationDeactivated(
	ctx context.Context,
	agglomeration models.Agglomeration,
) error {
	payload, err := json.Marshal(contracts.AgglomerationDeactivatedPayload{
		Agglomeration: agglomeration,
	})
	if err != nil {
		return err
	}

	_, err = p.outbox.CreateOutboxEvent(
		ctx,
		kafka.Message{
			Topic: contracts.AgglomerationsTopicV1,
			Key:   []byte(agglomeration.ID.String()),
			Value: payload,
			Headers: []kafka.Header{
				{Key: header.EventID, Value: []byte(uuid.New().String())},
				{Key: header.EventType, Value: []byte(contracts.AgglomerationDeactivatedEvent)},
				{Key: header.EventVersion, Value: []byte("1")},
				{Key: header.Producer, Value: []byte(contracts.AgglomerationsSvcGroup)},
				{Key: header.ContentType, Value: []byte("application/json")},
			},
		},
	)

	return err
}

func (p Producer) WriteAgglomerationSuspended(
	ctx context.Context,
	agglomeration models.Agglomeration,
) error {
	payload, err := json.Marshal(contracts.AgglomerationSuspendedPayload{
		Agglomeration: agglomeration,
	})
	if err != nil {
		return err
	}

	_, err = p.outbox.CreateOutboxEvent(
		ctx,
		kafka.Message{
			Topic: contracts.AgglomerationsTopicV1,
			Key:   []byte(agglomeration.ID.String()),
			Value: payload,
			Headers: []kafka.Header{
				{Key: header.EventID, Value: []byte(uuid.New().String())},
				{Key: header.EventType, Value: []byte(contracts.AgglomerationSuspendedEvent)},
				{Key: header.EventVersion, Value: []byte("1")},
				{Key: header.Producer, Value: []byte(contracts.AgglomerationsSvcGroup)},
				{Key: header.ContentType, Value: []byte("application/json")},
			},
		},
	)

	return err
}
