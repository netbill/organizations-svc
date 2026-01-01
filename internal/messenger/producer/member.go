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

func (s Service) WriteMemberCreated(
	ctx context.Context,
	member models.Member,
) error {
	payload, err := json.Marshal(contracts.MemberCreatedPayload{
		Member: member,
	})
	if err != nil {
		return err
	}

	_, err = s.outbox.CreateOutboxEvent(
		ctx,
		box.OutboxStatusPending,
		kafka.Message{
			Topic: contracts.MembersTopicV1,
			Key:   []byte(member.ID.String()),
			Value: payload,
			Headers: []kafka.Header{
				{Key: "EventID", Value: []byte(uuid.New().String())}, // Outbox will fill this
				{Key: "EventType", Value: []byte(contracts.MemberCreatedEvent)},
				{Key: "EventVersion", Value: []byte("1")},
				{Key: "Producer", Value: []byte(contracts.AgglomerationsSvcGroup)},
				{Key: "ContentType", Value: []byte("application/json")},
			},
		},
	)

	return err
}

func (s Service) WriteMemberUpdated(
	ctx context.Context,
	member models.Member,
) error {
	payload, err := json.Marshal(contracts.MemberUpdatedPayload{
		Member: member,
	})
	if err != nil {
		return err
	}

	_, err = s.outbox.CreateOutboxEvent(
		ctx,
		box.OutboxStatusPending,
		kafka.Message{
			Topic: contracts.MembersTopicV1,
			Key:   []byte(member.ID.String()),
			Value: payload,
			Headers: []kafka.Header{
				{Key: "EventID", Value: []byte(uuid.New().String())}, // Outbox will fill this
				{Key: "EventType", Value: []byte(contracts.MemberUpdatedEvent)},
				{Key: "EventVersion", Value: []byte("1")},
				{Key: "Producer", Value: []byte(contracts.AgglomerationsSvcGroup)},
				{Key: "ContentType", Value: []byte("application/json")},
			},
		},
	)

	return err
}

func (s Service) WriteMemberDeleted(
	ctx context.Context,
	member models.Member,
) error {
	payload, err := json.Marshal(contracts.MemberDeletedPayload{
		Member: member,
	})
	if err != nil {
		return err
	}

	_, err = s.outbox.CreateOutboxEvent(
		ctx,
		box.OutboxStatusPending,
		kafka.Message{
			Topic: contracts.MembersTopicV1,
			Key:   []byte(member.ID.String()),
			Value: payload,
			Headers: []kafka.Header{
				{Key: "EventID", Value: []byte(uuid.New().String())}, // Outbox will fill this
				{Key: "EventType", Value: []byte(contracts.MemberDeletedEvent)},
				{Key: "EventVersion", Value: []byte("1")},
				{Key: "Producer", Value: []byte(contracts.AgglomerationsSvcGroup)},
				{Key: "ContentType", Value: []byte("application/json")},
			},
		},
	)

	return err
}
