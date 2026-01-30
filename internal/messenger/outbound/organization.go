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

func (o Outbound) WriteOrganizationCreated(
	ctx context.Context,
	organization models.Organization,
) error {
	payload, err := json.Marshal(contracts.OrganizationCreatedPayload{
		OrganizationID: organization.ID,
		Status:         organization.Status,
		Name:           organization.Name,
		MaxRoles:       organization.MaxRoles,
		CreatedAt:      organization.CreatedAt,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal organization created payload: %w", err)
	}

	_, err = o.outbox.CreateOutboxEvent(
		ctx,
		kafka.Message{
			Topic: contracts.OrganizationsTopicV1,
			Key:   []byte(organization.ID.String()),
			Value: payload,
			Headers: []kafka.Header{
				{Key: header.EventID, Value: []byte(uuid.New().String())},
				{Key: header.EventType, Value: []byte(contracts.OrganizationCreatedEvent)},
				{Key: header.EventVersion, Value: []byte("1")},
				{Key: header.Producer, Value: []byte(contracts.OrganizationsSvcGroup)},
				{Key: header.ContentType, Value: []byte("application/json")},
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create outbox event for organization created: %w", err)
	}

	return nil
}

func (o Outbound) WriteOrganizationUpdated(
	ctx context.Context,
	organization models.Organization,
) error {
	payload, err := json.Marshal(contracts.OrganizationUpdatedPayload{
		OrganizationID: organization.ID,
		Status:         organization.Status,
		Name:           organization.Name,
		MaxRoles:       organization.MaxRoles,
		UpdatedAt:      organization.UpdatedAt,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal organization updated payload: %w", err)
	}

	_, err = o.outbox.CreateOutboxEvent(
		ctx,
		kafka.Message{
			Topic: contracts.OrganizationsTopicV1,
			Key:   []byte(organization.ID.String()),
			Value: payload,
			Headers: []kafka.Header{
				{Key: header.EventID, Value: []byte(uuid.New().String())},
				{Key: header.EventType, Value: []byte(contracts.OrganizationUpdatedEvent)},
				{Key: header.EventVersion, Value: []byte("1")},
				{Key: header.Producer, Value: []byte(contracts.OrganizationsSvcGroup)},
				{Key: header.ContentType, Value: []byte("application/json")},
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create outbox event for organization updated: %w", err)
	}

	return nil
}

func (o Outbound) WriteOrganizationDeleted(
	ctx context.Context,
	organization models.Organization,
) error {
	payload, err := json.Marshal(contracts.OrganizationDeletedPayload{
		OrganizationID: organization.ID,
		DeletedAt:      time.Now().UTC(),
	})
	if err != nil {
		return fmt.Errorf("failed to marshal organization deleted payload: %w", err)
	}

	_, err = o.outbox.CreateOutboxEvent(
		ctx,
		kafka.Message{
			Topic: contracts.OrganizationsTopicV1,
			Key:   []byte(organization.ID.String()),
			Value: payload,
			Headers: []kafka.Header{
				{Key: header.EventID, Value: []byte(uuid.New().String())},
				{Key: header.EventType, Value: []byte(contracts.OrganizationDeletedEvent)},
				{Key: header.EventVersion, Value: []byte("1")},
				{Key: header.Producer, Value: []byte(contracts.OrganizationsSvcGroup)},
				{Key: header.ContentType, Value: []byte("application/json")},
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create outbox event for organization deleted: %w", err)
	}

	return nil
}

func (o Outbound) WriteOrganizationActivated(
	ctx context.Context,
	organization models.Organization,
) error {
	payload, err := json.Marshal(contracts.OrganizationActivatedPayload{
		OrganizationID: organization.ID,
		ActivatedAt:    organization.UpdatedAt,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal organization activated payload: %w", err)
	}

	_, err = o.outbox.CreateOutboxEvent(
		ctx,
		kafka.Message{
			Topic: contracts.OrganizationsTopicV1,
			Key:   []byte(organization.ID.String()),
			Value: payload,
			Headers: []kafka.Header{
				{Key: header.EventID, Value: []byte(uuid.New().String())},
				{Key: header.EventType, Value: []byte(contracts.OrganizationActivatedEvent)},
				{Key: header.EventVersion, Value: []byte("1")},
				{Key: header.Producer, Value: []byte(contracts.OrganizationsSvcGroup)},
				{Key: header.ContentType, Value: []byte("application/json")},
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create outbox event for organization activated: %w", err)
	}

	return nil
}

func (o Outbound) WriteOrganizationDeactivated(
	ctx context.Context,
	organization models.Organization,
) error {
	payload, err := json.Marshal(contracts.OrganizationDeactivatedPayload{
		OrganizationID: organization.ID,
		DeactivatedAt:  organization.UpdatedAt,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal organization deactivated payload: %w", err)
	}

	_, err = o.outbox.CreateOutboxEvent(
		ctx,
		kafka.Message{
			Topic: contracts.OrganizationsTopicV1,
			Key:   []byte(organization.ID.String()),
			Value: payload,
			Headers: []kafka.Header{
				{Key: header.EventID, Value: []byte(uuid.New().String())},
				{Key: header.EventType, Value: []byte(contracts.OrganizationDeactivatedEvent)},
				{Key: header.EventVersion, Value: []byte("1")},
				{Key: header.Producer, Value: []byte(contracts.OrganizationsSvcGroup)},
				{Key: header.ContentType, Value: []byte("application/json")},
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create outbox event for organization deactivated: %w", err)
	}

	return nil
}
