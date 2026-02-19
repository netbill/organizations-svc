package sender

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

func (s *Sender) WriteOrgMemberCreated(
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

	_, err = s.outbox.WriteOutboxEvent(
		ctx,
		eventbox.Event{
			ID:       uuid.New(),
			Type:     evtypes.OrgMemberCreatedEvent,
			Version:  1,
			Topic:    evtypes.OrgMemberTopicV1,
			Key:      member.ID.String(),
			Payload:  payload,
			Producer: s.identity,
		},
	)

	if err != nil {
		return fmt.Errorf("failed to create sender event for org member created, cause: %w", err)
	}

	return nil
}

func (s *Sender) WriteOrgMemberUpdated(
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

	_, err = s.outbox.WriteOutboxEvent(
		ctx,
		eventbox.Event{
			ID:       uuid.New(),
			Type:     evtypes.OrgMemberUpdatedEvent,
			Version:  1,
			Topic:    evtypes.OrgMemberTopicV1,
			Key:      member.ID.String(),
			Payload:  payload,
			Producer: s.identity,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create sender event for org member updated, cause: %w", err)
	}

	return nil
}

func (s *Sender) WriteOrgMemberDeleted(
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

	_, err = s.outbox.WriteOutboxEvent(
		ctx,
		eventbox.Event{
			ID:       uuid.New(),
			Type:     evtypes.OrgMemberDeletedEvent,
			Version:  1,
			Topic:    evtypes.OrgMemberTopicV1,
			Key:      memberID.String(),
			Payload:  payload,
			Producer: s.identity,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create sender event for org member deleted, cause: %w", err)
	}

	return nil
}

func (s *Sender) WriteOrgMemberRoleAdd(
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

	_, err = s.outbox.WriteOutboxEvent(
		ctx,
		eventbox.Event{
			ID:       uuid.New(),
			Type:     evtypes.OrgMemberRoleAddedEvent,
			Version:  1,
			Topic:    evtypes.OrgMemberTopicV1,
			Key:      link.MemberID.String(),
			Payload:  payload,
			Producer: s.identity,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create sender event for org member role added, cause: %w", err)
	}

	return nil
}

func (s *Sender) WriteOrgMemberRoleRemove(
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

	_, err = s.outbox.WriteOutboxEvent(
		ctx,
		eventbox.Event{
			ID:       uuid.New(),
			Type:     evtypes.OrgMemberRoleRemovedEvent,
			Version:  1,
			Topic:    evtypes.OrgMemberTopicV1,
			Key:      memberID.String(),
			Payload:  payload,
			Producer: s.identity,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create sender event for org member role removed, cause: %w", err)
	}

	return nil
}
