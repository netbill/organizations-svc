package outbound

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/netbill/eventbox/headers"
	"github.com/netbill/organizations-svc/internal/core/models"
	"github.com/netbill/organizations-svc/pkg/evtypes"
	"github.com/segmentio/kafka-go"
)

func (o *Outbound) WriteOrgMemberCreated(
	ctx context.Context,
	member models.Member,
) error {
	payload, err := json.Marshal(evtypes.OrgMemberCreatedPayload{
		MemberID:       member.ID,
		AccountID:      member.AccountID,
		OrganizationID: member.OrganizationID,
		Head:           member.Head,
		Position:       member.Position,
		Label:          member.Label,
		CreatedAt:      member.CreatedAt,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal org member created payload, cause: %w", err)
	}

	_, err = o.outbox.WriteToOutbox(
		ctx,
		kafka.Message{
			Topic: evtypes.OrgMemberTopicV1,
			Key:   []byte(member.ID.String()),
			Value: payload,
			Headers: []kafka.Header{
				{Key: headers.EventID, Value: []byte(uuid.New().String())},
				{Key: headers.EventType, Value: []byte(evtypes.OrgMemberCreatedEvent)},
				{Key: headers.EventVersion, Value: []byte("1")},
				{Key: headers.Producer, Value: []byte(o.groupID)},
				{Key: headers.ContentType, Value: []byte("application/json")},
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
	payload, err := json.Marshal(evtypes.OrgMemberUpdatedPayload{
		MemberID:  member.ID,
		Position:  member.Position,
		Label:     member.Label,
		UpdatedAt: member.UpdatedAt,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal org member updated payload, cause: %w", err)
	}

	_, err = o.outbox.WriteToOutbox(
		ctx,
		kafka.Message{
			Topic: evtypes.OrgMemberTopicV1,
			Key:   []byte(member.ID.String()),
			Value: payload,
			Headers: []kafka.Header{
				{Key: headers.EventID, Value: []byte(uuid.New().String())},
				{Key: headers.EventType, Value: []byte(evtypes.OrgMemberUpdatedEvent)},
				{Key: headers.EventVersion, Value: []byte("1")},
				{Key: headers.Producer, Value: []byte(o.groupID)},
				{Key: headers.ContentType, Value: []byte("application/json")},
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
	payload, err := json.Marshal(evtypes.OrgMemberDeletedPayload{
		MemberID:  memberID,
		DeletedAt: time.Now().UTC(),
	})
	if err != nil {
		return fmt.Errorf("failed to marshal org member deleted payload, cause: %w", err)
	}

	_, err = o.outbox.WriteToOutbox(
		ctx,
		kafka.Message{
			Topic: evtypes.OrgMemberTopicV1,
			Key:   []byte(memberID.String()),
			Value: payload,
			Headers: []kafka.Header{
				{Key: headers.EventID, Value: []byte(uuid.New().String())},
				{Key: headers.EventType, Value: []byte(evtypes.OrgMemberDeletedEvent)},
				{Key: headers.EventVersion, Value: []byte("1")},
				{Key: headers.Producer, Value: []byte(o.groupID)},
				{Key: headers.ContentType, Value: []byte("application/json")},
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
	link models.OrgMemberRolesLink,
) error {
	payload, err := json.Marshal(evtypes.OrgMemberRoleAddedPayload{
		MemberID: link.MemberID,
		RoleID:   link.RoleID,
		AddedAt:  link.CreatedAt,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal org member role added payload, cause: %w", err)
	}

	_, err = o.outbox.WriteToOutbox(
		ctx,
		kafka.Message{
			Topic: evtypes.OrgMemberTopicV1,
			Key:   []byte(link.MemberID.String()),
			Value: payload,
			Headers: []kafka.Header{
				{Key: headers.EventID, Value: []byte(uuid.New().String())},
				{Key: headers.EventType, Value: []byte(evtypes.OrgMemberRoleAddedEvent)},
				{Key: headers.EventVersion, Value: []byte("1")},
				{Key: headers.Producer, Value: []byte(o.groupID)},
				{Key: headers.ContentType, Value: []byte("application/json")},
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
	payload, err := json.Marshal(evtypes.OrgMemberRoleRemovedPayload{
		MemberID: memberID,
		RoleID:   roleID,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal org member role removed payload, cause: %w", err)
	}

	_, err = o.outbox.WriteToOutbox(
		ctx,
		kafka.Message{
			Topic: evtypes.OrgMemberTopicV1,
			Key:   []byte(memberID.String()),
			Value: payload,
			Headers: []kafka.Header{
				{Key: headers.EventID, Value: []byte(uuid.New().String())},
				{Key: headers.EventType, Value: []byte(evtypes.OrgMemberRoleRemovedEvent)},
				{Key: headers.EventVersion, Value: []byte("1")},
				{Key: headers.Producer, Value: []byte(o.groupID)},
				{Key: headers.ContentType, Value: []byte("application/json")},
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create outbox event for org member role removed, cause: %w", err)
	}

	return nil
}
