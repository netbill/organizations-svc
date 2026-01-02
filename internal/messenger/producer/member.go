package producer

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/netbill/kafkakit/header"
	"github.com/netbill/organizations-svc/internal/domain/models"
	"github.com/netbill/organizations-svc/internal/messenger/contracts"
	"github.com/segmentio/kafka-go"
)

func (p Producer) WriteMemberCreated(
	ctx context.Context,
	member models.Member,
) error {
	payload, err := json.Marshal(contracts.MemberCreatedPayload{
		Member: member,
	})
	if err != nil {
		return err
	}

	_, err = p.outbox.CreateOutboxEvent(
		ctx,
		kafka.Message{
			Topic: contracts.MembersTopicV1,
			Key:   []byte(member.ID.String()),
			Value: payload,
			Headers: []kafka.Header{
				{Key: header.EventID, Value: []byte(uuid.New().String())},
				{Key: header.EventType, Value: []byte(contracts.MemberCreatedEvent)},
				{Key: header.EventVersion, Value: []byte("1")},
				{Key: header.Producer, Value: []byte(contracts.OrganizationsSvcGroup)},
				{Key: header.ContentType, Value: []byte("application/json")},
			},
		},
	)

	return err
}

func (p Producer) WriteMemberUpdated(
	ctx context.Context,
	member models.Member,
) error {
	payload, err := json.Marshal(contracts.MemberUpdatedPayload{
		Member: member,
	})
	if err != nil {
		return err
	}

	_, err = p.outbox.CreateOutboxEvent(
		ctx,
		kafka.Message{
			Topic: contracts.MembersTopicV1,
			Key:   []byte(member.ID.String()),
			Value: payload,
			Headers: []kafka.Header{
				{Key: header.EventID, Value: []byte(uuid.New().String())},
				{Key: header.EventType, Value: []byte(contracts.MemberUpdatedEvent)},
				{Key: header.EventVersion, Value: []byte("1")},
				{Key: header.Producer, Value: []byte(contracts.OrganizationsSvcGroup)},
				{Key: header.ContentType, Value: []byte("application/json")},
			},
		},
	)

	return err
}

func (p Producer) WriteMemberDeleted(
	ctx context.Context,
	member models.Member,
) error {
	payload, err := json.Marshal(contracts.MemberDeletedPayload{
		Member: member,
	})
	if err != nil {
		return err
	}

	_, err = p.outbox.CreateOutboxEvent(
		ctx,
		kafka.Message{
			Topic: contracts.MembersTopicV1,
			Key:   []byte(member.ID.String()),
			Value: payload,
			Headers: []kafka.Header{
				{Key: header.EventID, Value: []byte(uuid.New().String())},
				{Key: header.EventType, Value: []byte(contracts.MemberDeletedEvent)},
				{Key: header.EventVersion, Value: []byte("1")},
				{Key: header.Producer, Value: []byte(contracts.OrganizationsSvcGroup)},
				{Key: header.ContentType, Value: []byte("application/json")},
			},
		},
	)

	return err
}
