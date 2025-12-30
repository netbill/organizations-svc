package callbacker

import (
	"context"

	"github.com/segmentio/kafka-go"
	"github.com/umisto/kafkakit/box"
)

func (s Service) AccountUsernameChanged(ctx context.Context, event kafka.Message) error {
	_, err := s.inbox.CreateInboxEvent(ctx, box.InboxStatusPending, event)
	if err != nil {
		s.log.Errorf("failed to upsert inbox event for account %s: %v", string(event.Key), err)
		return err
	}

	return nil
}
