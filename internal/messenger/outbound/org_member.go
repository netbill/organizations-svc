package outbound

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/netbill/evebox/header"
	"github.com/netbill/organizations-svc/internal/core/models"
	"github.com/netbill/organizations-svc/internal/messenger/contracts"
	"github.com/segmentio/kafka-go"
)

func (o *Outbound) WriteOrgMemberCreated(
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
		return fmt.Errorf("failed to marshal org member created payload, cause: %w", err)
	}

	_, err = o.outbox.CreateOutboxEvent(
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

	if err != nil {
		return fmt.Errorf("failed to create outbox event for org member created, cause: %w", err)
	}

	return nil
}

func (o *Outbound) WriteOrgMemberUpdated(
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
		return fmt.Errorf("failed to marshal org member updated payload, cause: %w", err)
	}

	_, err = o.outbox.CreateOutboxEvent(
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
	if err != nil {
		return fmt.Errorf("failed to create outbox event for org member updated, cause: %w", err)
	}

	return nil
}

func (o *Outbound) WriteOrgMemberDeleted(
	ctx context.Context,
	memberID uuid.UUID,
) error {
	payload, err := json.Marshal(contracts.OrgMemberDeletedPayload{
		MemberID:  memberID,
		DeletedAt: time.Now().UTC(),
	})
	if err != nil {
		return fmt.Errorf("failed to marshal org member deleted payload, cause: %w", err)
	}

	_, err = o.outbox.CreateOutboxEvent(
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
	if err != nil {
		return fmt.Errorf("failed to create outbox event for org member deleted, cause: %w", err)
	}

	return nil
}

func (o *Outbound) WriteOrgMemberRoleAdd(
	ctx context.Context,
	memberID uuid.UUID,
	roleID uuid.UUID,
) error {
	payload, err := json.Marshal(contracts.OrgMemberRoleAddedPayload{
		MemberID: memberID,
		RoleID:   roleID,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal org member role added payload, cause: %w", err)
	}

	_, err = o.outbox.CreateOutboxEvent(
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
	if err != nil {
		return fmt.Errorf("failed to create outbox event for org member role added, cause: %w", err)
	}

	return nil
}

func (o *Outbound) WriteOrgMemberRoleRemove(
	ctx context.Context,
	memberID uuid.UUID,
	roleID uuid.UUID,
) error {
	payload, err := json.Marshal(contracts.OrgMemberRoleRemovedPayload{
		MemberID: memberID,
		RoleID:   roleID,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal org member role removed payload, cause: %w", err)
	}

	_, err = o.outbox.CreateOutboxEvent(
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
	if err != nil {
		return fmt.Errorf("failed to create outbox event for org member role removed, cause: %w", err)
	}

	return nil
}
