package producer

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/netbill/kafkakit/header"
	"github.com/netbill/organizations-svc/internal/messenger/contracts"
	"github.com/segmentio/kafka-go"
)

func (p Producer) WriteMemberRoleAdd(
	ctx context.Context,
	memberID uuid.UUID,
	roleID uuid.UUID,
) error {
	payload, err := json.Marshal(contracts.MemberRoleAddedPayload{
		MemberID: memberID,
		RoleID:   roleID,
	})
	if err != nil {
		return err
	}

	_, err = p.outbox.CreateOutboxEvent(
		ctx,
		kafka.Message{
			Topic: contracts.MembersTopicV1,
			Key:   []byte(memberID.String()),
			Value: payload,
			Headers: []kafka.Header{
				{Key: header.EventID, Value: []byte(uuid.New().String())},
				{Key: header.EventType, Value: []byte(contracts.MemberRoleAddedEvent)},
				{Key: header.EventVersion, Value: []byte("1")},
				{Key: header.Producer, Value: []byte(contracts.OrganizationsSvcGroup)},
				{Key: header.ContentType, Value: []byte("application/json")},
			},
		},
	)

	return err
}

func (p Producer) WriteMemberRoleRemove(
	ctx context.Context,
	memberID uuid.UUID,
	roleID uuid.UUID,
) error {
	payload, err := json.Marshal(contracts.MemberRoleRemovedPayload{
		MemberID: memberID,
		RoleID:   roleID,
	})
	if err != nil {
		return err
	}

	_, err = p.outbox.CreateOutboxEvent(
		ctx,
		kafka.Message{
			Topic: contracts.MembersTopicV1,
			Key:   []byte(memberID.String()),
			Value: payload,
			Headers: []kafka.Header{
				{Key: header.EventID, Value: []byte(uuid.New().String())},
				{Key: header.EventType, Value: []byte(contracts.MemberRoleRemovedEvent)},
				{Key: header.EventVersion, Value: []byte("1")},
				{Key: header.Producer, Value: []byte(contracts.OrganizationsSvcGroup)},
				{Key: header.ContentType, Value: []byte("application/json")},
			},
		},
	)

	return err
}
