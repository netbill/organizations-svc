package outbound

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/netbill/evebox/header"
	"github.com/netbill/organizations-svc/internal/core/models"
	"github.com/netbill/organizations-svc/internal/messenger/contracts"
	"github.com/segmentio/kafka-go"
)

func (p Outbound) WriteOrgMemberCreated(
	ctx context.Context,
	member models.Member,
) error {
	payload, err := json.Marshal(contracts.OrgMemberCreatedPayload{
		MemberID:       member.ID,
		AccountID:      member.AccountID,
		OrganizationID: member.OrganizationID,
		Position:       member.Position,
		Label:          member.Label,
		CreatedAt:      member.CreatedAt,
	})
	if err != nil {
		return err
	}

	_, err = p.outbox.CreateOutboxEvent(
		ctx,
		kafka.Message{
			Topic: contracts.OrgMemberTopicV1,
			Key:   []byte(member.ID.String()),
			Value: payload,
			Headers: []kafka.Header{
				{Key: header.EventID, Value: []byte(uuid.New().String())},
				{Key: header.EventType, Value: []byte(contracts.OrgMemberCreatedEvent)},
				{Key: header.EventVersion, Value: []byte("1")},
				{Key: header.Producer, Value: []byte(contracts.OrganizationsSvcGroup)},
				{Key: header.ContentType, Value: []byte("application/json")},
			},
		},
	)

	return err
}

func (p Outbound) WriteOrgMemberUpdated(
	ctx context.Context,
	member models.Member,
) error {
	payload, err := json.Marshal(contracts.OrgMemberUpdatedPayload{
		MemberID:  member.ID,
		Position:  member.Position,
		Label:     member.Label,
		UpdatedAt: member.UpdatedAt,
	})
	if err != nil {
		return err
	}

	_, err = p.outbox.CreateOutboxEvent(
		ctx,
		kafka.Message{
			Topic: contracts.OrgMemberTopicV1,
			Key:   []byte(member.ID.String()),
			Value: payload,
			Headers: []kafka.Header{
				{Key: header.EventID, Value: []byte(uuid.New().String())},
				{Key: header.EventType, Value: []byte(contracts.OrgMemberUpdatedEvent)},
				{Key: header.EventVersion, Value: []byte("1")},
				{Key: header.Producer, Value: []byte(contracts.OrganizationsSvcGroup)},
				{Key: header.ContentType, Value: []byte("application/json")},
			},
		},
	)

	return err
}

func (p Outbound) WriteOrgMemberDeleted(
	ctx context.Context,
	memberID uuid.UUID,
) error {
	payload, err := json.Marshal(contracts.OrgMemberDeletedPayload{
		MemberID:  memberID,
		DeletedAt: time.Now().UTC(),
	})
	if err != nil {
		return err
	}

	_, err = p.outbox.CreateOutboxEvent(
		ctx,
		kafka.Message{
			Topic: contracts.OrgMemberTopicV1,
			Key:   []byte(memberID.String()),
			Value: payload,
			Headers: []kafka.Header{
				{Key: header.EventID, Value: []byte(uuid.New().String())},
				{Key: header.EventType, Value: []byte(contracts.OrgMemberDeletedEvent)},
				{Key: header.EventVersion, Value: []byte("1")},
				{Key: header.Producer, Value: []byte(contracts.OrganizationsSvcGroup)},
				{Key: header.ContentType, Value: []byte("application/json")},
			},
		},
	)

	return err
}

func (p Outbound) WriteOrgMemberRoleAdd(
	ctx context.Context,
	memberID uuid.UUID,
	roleID uuid.UUID,
) error {
	payload, err := json.Marshal(contracts.OrgMemberRoleAddedPayload{
		MemberID: memberID,
		RoleID:   roleID,
	})
	if err != nil {
		return err
	}

	_, err = p.outbox.CreateOutboxEvent(
		ctx,
		kafka.Message{
			Topic: contracts.OrgMemberTopicV1,
			Key:   []byte(memberID.String()),
			Value: payload,
			Headers: []kafka.Header{
				{Key: header.EventID, Value: []byte(uuid.New().String())},
				{Key: header.EventType, Value: []byte(contracts.OrgMemberRoleAddedEvent)},
				{Key: header.EventVersion, Value: []byte("1")},
				{Key: header.Producer, Value: []byte(contracts.OrganizationsSvcGroup)},
				{Key: header.ContentType, Value: []byte("application/json")},
			},
		},
	)

	return err
}

func (p Outbound) WriteOrgMemberRoleRemove(
	ctx context.Context,
	memberID uuid.UUID,
	roleID uuid.UUID,
) error {
	payload, err := json.Marshal(contracts.OrgMemberRoleRemovedPayload{
		MemberID: memberID,
		RoleID:   roleID,
	})
	if err != nil {
		return err
	}

	_, err = p.outbox.CreateOutboxEvent(
		ctx,
		kafka.Message{
			Topic: contracts.OrgMemberTopicV1,
			Key:   []byte(memberID.String()),
			Value: payload,
			Headers: []kafka.Header{
				{Key: header.EventID, Value: []byte(uuid.New().String())},
				{Key: header.EventType, Value: []byte(contracts.OrgMemberRoleRemovedEvent)},
				{Key: header.EventVersion, Value: []byte("1")},
				{Key: header.Producer, Value: []byte(contracts.OrganizationsSvcGroup)},
				{Key: header.ContentType, Value: []byte("application/json")},
			},
		},
	)

	return err
}
