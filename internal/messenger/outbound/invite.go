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

func (o *Outbound) WriteOrgInviteCreated(
	ctx context.Context,
	invite models.Invite,
) error {
	payload, err := json.Marshal(evtypes.OrgInviteCreatedPayload{
		InviteID:       invite.ID,
		OrganizationID: invite.OrganizationID,
		AccountID:      invite.AccountID,
		Status:         invite.Status,
		CreatedAt:      invite.CreatedAt,
		ExpiresAt:      invite.ExpiresAt,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal org invite created payload, cause: %w", err)
	}

	_, err = o.outbox.WriteToOutbox(
		ctx,
		kafka.Message{
			Topic: evtypes.OrganizationsTopicV1,
			Key:   []byte(invite.OrganizationID.String()),
			Value: payload,
			Headers: []kafka.Header{
				{Key: headers.EventID, Value: []byte(uuid.New().String())},
				{Key: headers.EventType, Value: []byte(evtypes.OrgInviteCreatedEvent)},
				{Key: headers.EventVersion, Value: []byte("1")},
				{Key: headers.Producer, Value: []byte(o.groupID)},
				{Key: headers.ContentType, Value: []byte("application/json")},
			},
		},
	)

	if err != nil {
		return fmt.Errorf("failed to create outbox event for org invite created, cause: %w", err)
	}

	return nil
}

func (o *Outbound) WriteOrgInviteAccepted(
	ctx context.Context,
	invite models.Invite,
) error {
	payload, err := json.Marshal(evtypes.OrgInviteAcceptedPayload{
		InviteID:   invite.ID,
		AcceptedAt: time.Now().UTC(),
	})
	if err != nil {
		return fmt.Errorf("failed to marshal org invite accepted payload, cause: %w", err)
	}

	_, err = o.outbox.WriteToOutbox(
		ctx,
		kafka.Message{
			Topic: evtypes.OrganizationsTopicV1,
			Key:   []byte(invite.OrganizationID.String()),
			Value: payload,
			Headers: []kafka.Header{
				{Key: headers.EventID, Value: []byte(uuid.New().String())},
				{Key: headers.EventType, Value: []byte(evtypes.OrgInviteAcceptedEvent)},
				{Key: headers.EventVersion, Value: []byte("1")},
				{Key: headers.Producer, Value: []byte(o.groupID)},
				{Key: headers.ContentType, Value: []byte("application/json")},
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create outbox event for org invite accepted, cause: %w", err)
	}

	return nil
}

func (o *Outbound) WriteOrgInviteDeclined(
	ctx context.Context,
	invite models.Invite,
) error {
	payload, err := json.Marshal(evtypes.OrgInviteDeclinedPayload{
		InviteID:   invite.ID,
		DeclinedAt: time.Now().UTC(),
	})
	if err != nil {
		return fmt.Errorf("failed to marshal org invite declined payload, cause: %w", err)
	}

	_, err = o.outbox.WriteToOutbox(
		ctx,
		kafka.Message{
			Topic: evtypes.OrganizationsTopicV1,
			Key:   []byte(invite.OrganizationID.String()),
			Value: payload,
			Headers: []kafka.Header{
				{Key: headers.EventID, Value: []byte(uuid.New().String())},
				{Key: headers.EventType, Value: []byte(evtypes.OrgInviteDeclinedEvent)},
				{Key: headers.EventVersion, Value: []byte("1")},
				{Key: headers.Producer, Value: []byte(o.groupID)},
				{Key: headers.ContentType, Value: []byte("application/json")},
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create outbox event for org invite declined, cause: %w", err)
	}

	return nil
}

func (o *Outbound) WriteOrgInviteDeleted(
	ctx context.Context,
	invite models.Invite,
) error {
	payload, err := json.Marshal(evtypes.OrgInviteDeletedPayload{
		InvitedID: invite.ID,
		DeletedAt: time.Now().UTC(),
	})
	if err != nil {
		return fmt.Errorf("failed to marshal org invite deleted payload, cause: %w", err)
	}

	_, err = o.outbox.WriteToOutbox(
		ctx,
		kafka.Message{
			Topic: evtypes.OrganizationsTopicV1,
			Key:   []byte(invite.OrganizationID.String()),
			Value: payload,
			Headers: []kafka.Header{
				{Key: headers.EventID, Value: []byte(uuid.New().String())},
				{Key: headers.EventType, Value: []byte(evtypes.OrgInviteDeletedEvent)},
				{Key: headers.EventVersion, Value: []byte("1")},
				{Key: headers.Producer, Value: []byte(o.groupID)},
				{Key: headers.ContentType, Value: []byte("application/json")},
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create outbox event for org invite deleted, cause: %w", err)
	}

	return nil
}
