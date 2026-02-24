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

func (p *Publisher) WriteOrganizationCreated(
	ctx context.Context,
	organization models.Organization,
) error {
	payload, err := json.Marshal(evtypes.OrganizationCreatedPayload{
		OrganizationID: organization.ID,
		Status:         organization.Status,
		Name:           organization.Name,
		MaxRoles:       organization.MaxRoles,
		CreatedAt:      organization.CreatedAt,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal organization created payload: %w", err)
	}

	_, err = p.outbox.WriteOutboxEvent(ctx, eventbox.Message{
		ID:       uuid.New(),
		Type:     evtypes.OrganizationCreatedEvent,
		Version:  1,
		Topic:    evtypes.OrganizationsTopicV1,
		Key:      organization.ID.String(),
		Payload:  payload,
		Producer: p.identity,
	})
	if err != nil {
		return fmt.Errorf("failed to create sender event for organization created: %w", err)
	}

	return nil
}

func (p *Publisher) WriteOrganizationUpdated(
	ctx context.Context,
	organization models.Organization,
) error {
	payload, err := json.Marshal(evtypes.OrganizationUpdatedPayload{
		OrganizationID: organization.ID,
		Status:         organization.Status,
		Name:           organization.Name,
		MaxRoles:       organization.MaxRoles,
		UpdatedAt:      organization.UpdatedAt,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal organization updated payload: %w", err)
	}

	_, err = p.outbox.WriteOutboxEvent(ctx, eventbox.Message{
		ID:       uuid.New(),
		Type:     evtypes.OrganizationUpdatedEvent,
		Version:  1,
		Topic:    evtypes.OrganizationsTopicV1,
		Key:      organization.ID.String(),
		Payload:  payload,
		Producer: p.identity,
	})
	if err != nil {
		return fmt.Errorf("failed to create sender event for organization updated: %w", err)
	}

	return nil
}

func (p *Publisher) WriteOrganizationDeleted(
	ctx context.Context,
	organization models.Organization,
) error {
	payload, err := json.Marshal(evtypes.OrganizationDeletedPayload{
		OrganizationID: organization.ID,
		DeletedAt:      time.Now().UTC(),
	})
	if err != nil {
		return fmt.Errorf("failed to marshal organization deleted payload: %w", err)
	}

	_, err = p.outbox.WriteOutboxEvent(ctx, eventbox.Message{
		ID:       uuid.New(),
		Type:     evtypes.OrganizationDeletedEvent,
		Version:  1,
		Topic:    evtypes.OrganizationsTopicV1,
		Key:      organization.ID.String(),
		Payload:  payload,
		Producer: p.identity,
	})
	if err != nil {
		return fmt.Errorf("failed to create sender event for organization deleted: %w", err)
	}

	return nil
}
