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

func (p *Publisher) WriteOrgMemberCreated(
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
		Version:        member.Version,
		CreatedAt:      member.CreatedAt,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal org member created payload, cause: %w", err)
	}

	_, err = p.outbox.WriteOutboxEvent(ctx, eventbox.Message{
		ID:       uuid.New(),
		Type:     evtypes.OrgMemberCreatedEvent,
		Version:  1,
		Topic:    evtypes.OrgMembersTopicV1,
		Key:      member.ID.String(),
		Payload:  payload,
		Producer: p.identity,
	})

	if err != nil {
		return fmt.Errorf("failed to create sender event for org member created, cause: %w", err)
	}

	return nil
}

func (p *Publisher) WriteOrgMemberUpdated(
	ctx context.Context,
	member models.Member,
) error {
	payload, err := json.Marshal(evtypes.OrgMemberUpdatedPayload{
		MemberID:  member.ID,
		Position:  member.Position,
		Label:     member.Label,
		Version:   member.Version,
		UpdatedAt: member.UpdatedAt,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal org member updated payload, cause: %w", err)
	}

	_, err = p.outbox.WriteOutboxEvent(ctx, eventbox.Message{
		ID:       uuid.New(),
		Type:     evtypes.OrgMemberUpdatedEvent,
		Version:  1,
		Topic:    evtypes.OrgMembersTopicV1,
		Key:      member.ID.String(),
		Payload:  payload,
		Producer: p.identity,
	})
	if err != nil {
		return fmt.Errorf("failed to create sender event for org member updated, cause: %w", err)
	}

	return nil
}

func (p *Publisher) WriteOrgMemberDeleted(
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

	_, err = p.outbox.WriteOutboxEvent(ctx, eventbox.Message{
		ID:       uuid.New(),
		Type:     evtypes.OrgMemberDeletedEvent,
		Version:  1,
		Topic:    evtypes.OrgMembersTopicV1,
		Key:      memberID.String(),
		Payload:  payload,
		Producer: p.identity,
	})
	if err != nil {
		return fmt.Errorf("failed to create sender event for org member deleted, cause: %w", err)
	}

	return nil
}
