package outbound

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/netbill/evebox/header"
	"github.com/netbill/organizations-svc/internal/core/models"
	"github.com/netbill/organizations-svc/internal/messenger/contracts"
	"github.com/segmentio/kafka-go"
)

func (p Outbound) WriteOrgRoleCreated(
	ctx context.Context,
	role models.Role,
) error {
	payload, err := json.Marshal(contracts.OrgRoleCreatedPayload{
		RoleID:         role.ID,
		OrganizationID: role.OrganizationID,
		Head:           role.Head,
		Rank:           role.Rank,
		Name:           role.Name,
		Description:    role.Description,
		Color:          role.Color,
		CreatedAt:      role.CreatedAt,
	})
	if err != nil {
		return err
	}

	_, err = p.outbox.CreateOutboxEvent(
		ctx,
		kafka.Message{
			Topic: contracts.OrganizationsTopicV1,
			Key:   []byte(role.OrganizationID.String()),
			Value: payload,
			Headers: []kafka.Header{
				{Key: header.EventID, Value: []byte(uuid.New().String())},
				{Key: header.EventType, Value: []byte(contracts.OrgRoleCreatedEvent)},
				{Key: header.EventVersion, Value: []byte("1")},
				{Key: header.Producer, Value: []byte(contracts.OrganizationsSvcGroup)},
				{Key: header.ContentType, Value: []byte("application/json")},
			},
		},
	)

	return err
}

func (p Outbound) WriteOrgRoleUpdated(
	ctx context.Context,
	role models.Role,
) error {
	payload, err := json.Marshal(contracts.OrgRoleUpdatedPayload{
		RoleID:      role.ID,
		Name:        role.Name,
		Description: role.Description,
		Color:       role.Color,
		UpdatedAt:   role.UpdatedAt,
	})
	if err != nil {
		return err
	}

	_, err = p.outbox.CreateOutboxEvent(
		ctx,
		kafka.Message{
			Topic: contracts.OrganizationsTopicV1,
			Key:   []byte(role.OrganizationID.String()),
			Value: payload,
			Headers: []kafka.Header{
				{Key: header.EventID, Value: []byte(uuid.New().String())},
				{Key: header.EventType, Value: []byte(contracts.OrgRoleUpdatedEvent)},
				{Key: header.EventVersion, Value: []byte("1")},
				{Key: header.Producer, Value: []byte(contracts.OrganizationsSvcGroup)},
				{Key: header.ContentType, Value: []byte("application/json")},
			},
		},
	)

	return err
}

func (p Outbound) WriteOrgRoleDeleted(
	ctx context.Context,
	role models.Role,
) error {
	payload, err := json.Marshal(contracts.OrgRoleDeletedPayload{
		RoleID:    role.ID,
		DeletedAt: time.Now().UTC(),
	})
	if err != nil {
		return err
	}

	_, err = p.outbox.CreateOutboxEvent(
		ctx,
		kafka.Message{
			Topic: contracts.OrganizationsTopicV1,
			Key:   []byte(role.OrganizationID.String()),
			Value: payload,
			Headers: []kafka.Header{
				{Key: header.EventID, Value: []byte(uuid.New().String())},
				{Key: header.EventType, Value: []byte(contracts.OrgRoleDeletedEvent)},
				{Key: header.EventVersion, Value: []byte("1")},
				{Key: header.Producer, Value: []byte(contracts.OrganizationsSvcGroup)},
				{Key: header.ContentType, Value: []byte("application/json")},
			},
		},
	)

	return err
}

func (p Outbound) WriteOrgRolePermissionsUpdated(
	ctx context.Context,
	role models.Role,
	permissions map[models.Permission]bool,
) error {
	per := make(map[string]bool)
	for k, v := range permissions {
		per[k.Code] = v
	}

	payload, err := json.Marshal(contracts.OrgRolePermissionsUpdatedPayload{
		RoleID:      role.ID,
		Permissions: per,
	})
	if err != nil {
		return err
	}

	_, err = p.outbox.CreateOutboxEvent(
		ctx,
		kafka.Message{
			Topic: contracts.OrganizationsTopicV1,
			Key:   []byte(role.OrganizationID.String()),
			Value: payload,
			Headers: []kafka.Header{
				{Key: header.EventID, Value: []byte(uuid.New().String())},
				{Key: header.EventType, Value: []byte(contracts.OrgRolePermissionsUpdatedEvent)},
				{Key: header.EventVersion, Value: []byte("1")},
				{Key: header.Producer, Value: []byte(contracts.OrganizationsSvcGroup)},
				{Key: header.ContentType, Value: []byte("application/json")},
			},
		},
	)

	return err
}

func (p Outbound) WriteOrgRolesRanksUpdated(
	ctx context.Context,
	organizationID uuid.UUID,
	ranks map[uuid.UUID]uint,
) error {
	payload, err := json.Marshal(contracts.RolesRanksUpdatedPayload{
		OrganizationID: organizationID,
		Ranks:          ranks,
	})
	if err != nil {
		return err
	}

	_, err = p.outbox.CreateOutboxEvent(
		ctx,
		kafka.Message{
			Topic: contracts.OrganizationsTopicV1,
			Key:   []byte(organizationID.String()),
			Value: payload,
			Headers: []kafka.Header{
				{Key: header.EventID, Value: []byte(uuid.New().String())},
				{Key: header.EventType, Value: []byte(contracts.OrgRolesRanksUpdatedEvent)},
				{Key: header.EventVersion, Value: []byte("1")},
				{Key: header.Producer, Value: []byte(contracts.OrganizationsSvcGroup)},
				{Key: header.ContentType, Value: []byte("application/json")},
			},
		},
	)

	return err
}
