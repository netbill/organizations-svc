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

func (p Producer) WriteInviteCreated(
	ctx context.Context,
	invite models.Invite,
) error {
	payload, err := json.Marshal(contracts.InviteCreatedPayload{
		Invite: invite,
	})
	if err != nil {
		return err
	}

	_, err = p.outbox.CreateOutboxEvent(
		ctx,
		kafka.Message{
			Topic: contracts.InvitesTopicV1,
			Key:   []byte(invite.ID.String()),
			Value: payload,
			Headers: []kafka.Header{
				{Key: header.EventID, Value: []byte(uuid.New().String())},
				{Key: header.EventType, Value: []byte(contracts.InviteCreatedEvent)},
				{Key: header.EventVersion, Value: []byte("1")},
				{Key: header.Producer, Value: []byte(contracts.AgglomerationsSvcGroup)},
				{Key: header.ContentType, Value: []byte("application/json")},
			},
		},
	)

	return err
}

func (p Producer) WriteInviteAccepted(
	ctx context.Context,
	invite models.Invite,
) error {
	payload, err := json.Marshal(contracts.InviteAcceptedPayload{
		Invite: invite,
	})
	if err != nil {
		return err
	}

	_, err = p.outbox.CreateOutboxEvent(
		ctx,
		kafka.Message{
			Topic: contracts.InvitesTopicV1,
			Key:   []byte(invite.ID.String()),
			Value: payload,
			Headers: []kafka.Header{
				{Key: header.EventID, Value: []byte(uuid.New().String())},
				{Key: header.EventType, Value: []byte(contracts.InviteAcceptedEvent)},
				{Key: header.EventVersion, Value: []byte("1")},
				{Key: header.Producer, Value: []byte(contracts.AgglomerationsSvcGroup)},
				{Key: header.ContentType, Value: []byte("application/json")},
			},
		},
	)

	return err
}

func (p Producer) WriteInviteDeclined(
	ctx context.Context,
	invite models.Invite,
) error {
	payload, err := json.Marshal(contracts.InviteAcceptedPayload{
		Invite: invite,
	})
	if err != nil {
		return err
	}

	_, err = p.outbox.CreateOutboxEvent(
		ctx,
		kafka.Message{
			Topic: contracts.InvitesTopicV1,
			Key:   []byte(invite.ID.String()),
			Value: payload,
			Headers: []kafka.Header{
				{Key: header.EventID, Value: []byte(uuid.New().String())},
				{Key: header.EventType, Value: []byte(contracts.InviteDeclinedEvent)},
				{Key: header.EventVersion, Value: []byte("1")},
				{Key: header.Producer, Value: []byte(contracts.AgglomerationsSvcGroup)},
				{Key: header.ContentType, Value: []byte("application/json")},
			},
		},
	)

	return err
}

func (p Producer) WriteInviteDeleted(
	ctx context.Context,
	invite models.Invite,
) error {
	payload, err := json.Marshal(contracts.InviteDeletedPayload{
		Invite: invite,
	})
	if err != nil {
		return err
	}

	_, err = p.outbox.CreateOutboxEvent(
		ctx,
		kafka.Message{
			Topic: contracts.InvitesTopicV1,
			Key:   []byte(invite.ID.String()),
			Value: payload,
			Headers: []kafka.Header{
				{Key: header.EventID, Value: []byte(uuid.New().String())},
				{Key: header.EventType, Value: []byte(contracts.InviteDeletedEvent)},
				{Key: header.EventVersion, Value: []byte("1")},
				{Key: header.Producer, Value: []byte(contracts.AgglomerationsSvcGroup)},
				{Key: header.ContentType, Value: []byte("application/json")},
			},
		},
	)

	return err
}
