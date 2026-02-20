package outbound

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/netbill/eventbox"
	"github.com/netbill/organizations-svc/internal/core/models"
	"github.com/netbill/organizations-svc/pkg/evtypes"
)

func (s *Sender) WriteOrgInviteCreated(
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

	_, err = s.outbox.WriteOutboxEvent(
		ctx,
		eventbox.Event{
			ID:       uuid.New(),
			Type:     evtypes.OrgInviteCreatedEvent,
			Version:  1,
			Topic:    evtypes.OrganizationsTopicV1,
			Key:      invite.OrganizationID.String(),
			Payload:  payload,
			Producer: s.identity,
		},
	)

	if err != nil {
		return fmt.Errorf("failed to create sender event for org invite created, cause: %w", err)
	}

	return nil
}

func (s *Sender) WriteOrgInviteAccepted(
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

	_, err = s.outbox.WriteOutboxEvent(ctx, eventbox.Event{
		ID:       uuid.New(),
		Type:     evtypes.OrgInviteAcceptedEvent,
		Version:  1,
		Topic:    evtypes.OrganizationsTopicV1,
		Key:      invite.OrganizationID.String(),
		Payload:  payload,
		Producer: s.identity,
	})
	if err != nil {
		return fmt.Errorf("failed to create sender event for org invite accepted, cause: %w", err)
	}

	return nil
}

func (s *Sender) WriteOrgInviteDeclined(
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

	_, err = s.outbox.WriteOutboxEvent(ctx, eventbox.Event{
		ID:       uuid.New(),
		Type:     evtypes.OrgInviteDeclinedEvent,
		Version:  1,
		Topic:    evtypes.OrganizationsTopicV1,
		Key:      invite.OrganizationID.String(),
		Payload:  payload,
		Producer: s.identity,
	})
	if err != nil {
		return fmt.Errorf("failed to create sender event for org invite declined, cause: %w", err)
	}

	return nil
}

func (s *Sender) WriteOrgInviteDeleted(
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

	_, err = s.outbox.WriteOutboxEvent(ctx, eventbox.Event{
		ID:       uuid.New(),
		Type:     evtypes.OrgInviteDeletedEvent,
		Version:  1,
		Topic:    evtypes.OrganizationsTopicV1,
		Key:      invite.OrganizationID.String(),
		Payload:  payload,
		Producer: s.identity,
	})
	if err != nil {
		return fmt.Errorf("failed to create sender event for org invite deleted, cause: %w", err)
	}

	return nil
}
