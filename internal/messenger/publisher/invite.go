package publisher

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/netbill/eventbox"
	"github.com/netbill/evtypes"
	"github.com/netbill/organizations-svc/internal/core/models"
)

func (p *Publisher) WriteOrgInviteCreated(
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

	_, err = p.outbox.WriteOutboxEvent(ctx, eventbox.Message{
		ID:       uuid.New(),
		Type:     evtypes.OrgInviteCreatedEvent,
		Version:  1,
		Topic:    evtypes.OrgMembersTopicV1,
		Key:      invite.OrganizationID.String(),
		Payload:  payload,
		Producer: p.identity,
	})
	return err
}

func (p *Publisher) WriteOrgInviteAccepted(
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

	_, err = p.outbox.WriteOutboxEvent(ctx, eventbox.Message{
		ID:       uuid.New(),
		Type:     evtypes.OrgInviteAcceptedEvent,
		Version:  1,
		Topic:    evtypes.OrgMembersTopicV1,
		Key:      invite.OrganizationID.String(),
		Payload:  payload,
		Producer: p.identity,
	})
	return err
}

func (p *Publisher) WriteOrgInviteDeclined(
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

	_, err = p.outbox.WriteOutboxEvent(ctx, eventbox.Message{
		ID:       uuid.New(),
		Type:     evtypes.OrgInviteDeclinedEvent,
		Version:  1,
		Topic:    evtypes.OrgMembersTopicV1,
		Key:      invite.OrganizationID.String(),
		Payload:  payload,
		Producer: p.identity,
	})
	return err
}

func (p *Publisher) WriteOrgInviteCanceled(
	ctx context.Context,
	invite models.Invite,
) error {
	payload, err := json.Marshal(evtypes.OrgInviteCancelledPayload{
		InviteID:    invite.ID,
		CancelledAt: time.Now().UTC(),
	})
	if err != nil {
		return fmt.Errorf("failed to marshal org invite deleted payload, cause: %w", err)
	}

	_, err = p.outbox.WriteOutboxEvent(ctx, eventbox.Message{
		ID:       uuid.New(),
		Type:     evtypes.OrgInviteCancelledEvent,
		Version:  1,
		Topic:    evtypes.OrgMembersTopicV1,
		Key:      invite.OrganizationID.String(),
		Payload:  payload,
		Producer: p.identity,
	})
	return err
}
