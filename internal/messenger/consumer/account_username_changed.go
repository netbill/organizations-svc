package consumer

import (
	"context"

	"github.com/segmentio/kafka-go"
)

func (c Consumer) AccountUsernameChanged(ctx context.Context, event kafka.Message) error {
	return c.inbox.Transaction(ctx, func(ctx context.Context) error {
		eventInBox, err := c.inbox.CreateInboxEvent(ctx, event)
		if err != nil {
			c.log.Errorf("failed to upsert inbox event for account %s: %v", string(event.Key), err)
			return err
		}

		if _, err = c.inbox.UpdateInboxEventStatus(ctx, eventInBox.ID, c.callbacks.AccountUsernameChanged(ctx, eventInBox)); err != nil {
			c.log.Errorf(
				"failed to update inbox event status for key %s, id: %s, error: %v", eventInBox.Key, eventInBox.ID, err,
			)
		}

		return nil
	})
}
