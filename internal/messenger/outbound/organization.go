package outbound

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/netbill/eventbox/headers"
	"github.com/netbill/organizations-svc/internal/core/models"
	"github.com/netbill/organizations-svc/internal/messenger/evtypes"
	"github.com/segmentio/kafka-go"
)

func (o *Outbound) WriteOrganizationCreated(
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

	_, err = o.outbox.WriteToOutbox(
		ctx,
		kafka.Message{
			Topic: evtypes.OrganizationsTopicV1,
			Key:   []byte(organization.ID.String()),
			Value: payload,
			Headers: []kafka.Header{
				{Key: headers.EventID, Value: []byte(uuid.New().String())},
				{Key: headers.EventType, Value: []byte(evtypes.OrganizationCreatedEvent)},
				{Key: headers.EventVersion, Value: []byte("1")},
				{Key: headers.Producer, Value: []byte(evtypes.OrganizationsSvcGroup)},
				{Key: headers.ContentType, Value: []byte("application/json")},
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create outbox event for organization created: %w", err)
	}

	return nil
}

func (o *Outbound) WriteOrganizationUpdated(
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

	_, err = o.outbox.WriteToOutbox(
		ctx,
		kafka.Message{
			Topic: evtypes.OrganizationsTopicV1,
			Key:   []byte(organization.ID.String()),
			Value: payload,
			Headers: []kafka.Header{
				{Key: headers.EventID, Value: []byte(uuid.New().String())},
				{Key: headers.EventType, Value: []byte(evtypes.OrganizationUpdatedEvent)},
				{Key: headers.EventVersion, Value: []byte("1")},
				{Key: headers.Producer, Value: []byte(evtypes.OrganizationsSvcGroup)},
				{Key: headers.ContentType, Value: []byte("application/json")},
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create outbox event for organization updated: %w", err)
	}

	return nil
}

func (o *Outbound) WriteOrganizationDeleted(
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

	_, err = o.outbox.WriteToOutbox(
		ctx,
		kafka.Message{
			Topic: evtypes.OrganizationsTopicV1,
			Key:   []byte(organization.ID.String()),
			Value: payload,
			Headers: []kafka.Header{
				{Key: headers.EventID, Value: []byte(uuid.New().String())},
				{Key: headers.EventType, Value: []byte(evtypes.OrganizationDeletedEvent)},
				{Key: headers.EventVersion, Value: []byte("1")},
				{Key: headers.Producer, Value: []byte(evtypes.OrganizationsSvcGroup)},
				{Key: headers.ContentType, Value: []byte("application/json")},
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create outbox event for organization deleted: %w", err)
	}

	return nil
}

func (o *Outbound) WriteOrganizationActivated(
	ctx context.Context,
	organization models.Organization,
) error {
	payload, err := json.Marshal(evtypes.OrganizationActivatedPayload{
		OrganizationID: organization.ID,
		ActivatedAt:    organization.UpdatedAt,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal organization activated payload: %w", err)
	}

	_, err = o.outbox.WriteToOutbox(
		ctx,
		kafka.Message{
			Topic: evtypes.OrganizationsTopicV1,
			Key:   []byte(organization.ID.String()),
			Value: payload,
			Headers: []kafka.Header{
				{Key: headers.EventID, Value: []byte(uuid.New().String())},
				{Key: headers.EventType, Value: []byte(evtypes.OrganizationActivatedEvent)},
				{Key: headers.EventVersion, Value: []byte("1")},
				{Key: headers.Producer, Value: []byte(evtypes.OrganizationsSvcGroup)},
				{Key: headers.ContentType, Value: []byte("application/json")},
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create outbox event for organization activated: %w", err)
	}

	return nil
}

func (o *Outbound) WriteOrganizationDeactivated(
	ctx context.Context,
	organization models.Organization,
) error {
	payload, err := json.Marshal(evtypes.OrganizationDeactivatedPayload{
		OrganizationID: organization.ID,
		DeactivatedAt:  organization.UpdatedAt,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal organization deactivated payload: %w", err)
	}

	_, err = o.outbox.WriteToOutbox(
		ctx,
		kafka.Message{
			Topic: evtypes.OrganizationsTopicV1,
			Key:   []byte(organization.ID.String()),
			Value: payload,
			Headers: []kafka.Header{
				{Key: headers.EventID, Value: []byte(uuid.New().String())},
				{Key: headers.EventType, Value: []byte(evtypes.OrganizationDeactivatedEvent)},
				{Key: headers.EventVersion, Value: []byte("1")},
				{Key: headers.Producer, Value: []byte(evtypes.OrganizationsSvcGroup)},
				{Key: headers.ContentType, Value: []byte("application/json")},
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create outbox event for organization deactivated: %w", err)
	}

	return nil
}
