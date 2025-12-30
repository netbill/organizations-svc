package consumer

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/umisto/cities-svc/internal/domain/models"
	"github.com/umisto/cities-svc/internal/messenger/contracts"
	"github.com/umisto/kafkakit/box"
	"github.com/umisto/logium"
)

type InboxWorker struct {
	log    logium.Logger
	inbox  inbox
	domain domain
}

type inbox interface {
	GetInboxEventByID(
		ctx context.Context,
		id uuid.UUID,
	) (box.InboxEvent, error)

	GetPendingInboxEvents(
		ctx context.Context,
		limit int32,
	) ([]box.InboxEvent, error)

	MarkInboxEventsAsProcessed(
		ctx context.Context,
		ids []uuid.UUID,
	) ([]box.InboxEvent, error)

	MarkInboxEventsAsFailed(
		ctx context.Context,
		ids []uuid.UUID,
	) ([]box.InboxEvent, error)

	MarkInboxEventsAsPending(
		ctx context.Context,
		ids []uuid.UUID,
		delay time.Duration,
	) ([]box.InboxEvent, error)
}

type domain interface {
	UpsertProfile(ctx context.Context, profile models.Profile) (models.Profile, error)
	UpdateUsername(ctx context.Context, accountID uuid.UUID, username string) (models.Profile, error)
	DeleteProfile(
		ctx context.Context,
		accountID uuid.UUID,
	) error
}

func NewInboxWorker(
	log logium.Logger,
	inbox inbox,
	domain domain,
) InboxWorker {
	return InboxWorker{
		log:    log,
		inbox:  inbox,
		domain: domain,
	}
}

const eventInboxRetryDelay = 1 * time.Minute

func (w InboxWorker) Run(ctx context.Context) {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}

		events, err := w.inbox.GetPendingInboxEvents(ctx, 10)
		if err != nil {
			w.log.Errorf("failed to get pending inbox events, cause: %v", err)
			continue
		}
		if len(events) == 0 {
			continue
		}

		var processed []uuid.UUID
		var delayed []uuid.UUID

		for _, event := range events {
			w.log.Infof("processing inbox event: %s, type %s", event.ID, event.Type)

			key, err := uuid.Parse(event.Key)
			if err != nil {
				w.log.Errorf("bad inbox event key, id: %s, key: %s, error: %v", event.ID, event.Key, err)
				processed = append(processed, event.ID)
				continue
			}

			switch event.Type {
			case contracts.AccountUsernameChangeEvent:
				var p contracts.AccountUsernameChangePayload
				if err = json.Unmarshal(event.Payload, &p); err != nil {
					w.log.Errorf("bad payload for %s, id: %s, error: %v", event.Type, event.ID, err)
					processed = append(processed, event.ID)
					continue
				}
				if _, err = w.domain.UpdateUsername(ctx, key, p.Account.Username); err != nil {
					w.log.Errorf("failed to update username, id: %s, error: %v", event.ID, err)
					delayed = append(delayed, event.ID)
					continue
				}

			case contracts.ProfileUpdatedEvent:
				var p contracts.ProfileUpdatedPayload
				if err = json.Unmarshal(event.Payload, &p); err != nil {
					w.log.Errorf("bad payload for %s, id: %s, error: %v", event.Type, event.ID, err)
					processed = append(processed, event.ID)
					continue
				}
				profile := models.Profile{
					AccountID: key,
					Username:  p.Profile.Username,
					Official:  p.Profile.Official,
					Pseudonym: p.Profile.Pseudonym,
				}
				if _, err = w.domain.UpsertProfile(ctx, profile); err != nil {
					w.log.Errorf("failed to upsert profile, id: %s, error: %v", event.ID, err)
					delayed = append(delayed, event.ID)
					continue
				}

			case contracts.AccountCreatedEvent:
				var p contracts.AccountCreatedPayload
				if err = json.Unmarshal(event.Payload, &p); err != nil {
					w.log.Errorf("bad payload for %s, id: %s, error: %v", event.Type, event.ID, err)
					processed = append(processed, event.ID)
					continue
				}
				profile := models.Profile{
					AccountID: key,
					Username:  p.Account.Username,
				}
				if _, err = w.domain.UpsertProfile(ctx, profile); err != nil {
					w.log.Errorf("failed to create profile, id: %s, error: %v", event.ID, err)
					delayed = append(delayed, event.ID)
					continue
				}

			case contracts.AccountDeletedEvent:
				var p contracts.AccountDeletedPayload
				if err = json.Unmarshal(event.Payload, &p); err != nil {
					w.log.Errorf("bad payload for %s, id: %s, error: %v", event.Type, event.ID, err)
					processed = append(processed, event.ID)
					continue
				}

				if err = w.domain.DeleteProfile(ctx, p.AccountID); err != nil {
					w.log.Errorf("failed to delete profile, id: %s, error: %v", event.ID, err)
					delayed = append(delayed, event.ID)
					continue
				}

			default:
				w.log.Errorf("unknown inbox event type: %s, id: %s", event.Type, event.ID)
			}

			processed = append(processed, event.ID)
		}

		if len(processed) > 0 {
			_, err = w.inbox.MarkInboxEventsAsProcessed(ctx, processed)
			if err != nil {
				w.log.Errorf("failed to mark inbox events as processed, ids: %v, error: %v", processed, err)
			}
		}

		if len(delayed) > 0 {
			_, err = w.inbox.MarkInboxEventsAsPending(ctx, delayed, eventInboxRetryDelay)
			if err != nil {
				w.log.Errorf("failed to mark inbox events as pending, ids: %v, error: %v", delayed, err)
			}
		}
	}
}
