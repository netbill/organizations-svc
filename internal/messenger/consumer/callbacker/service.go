package callbacker

import (
	"context"

	"github.com/segmentio/kafka-go"
	"github.com/umisto/kafkakit/box"
	"github.com/umisto/logium"
)

type Inbox interface {
	CreateInboxEvent(
		ctx context.Context,
		status string,
		message kafka.Message,
	) (box.InboxEvent, error)
}

type Service struct {
	inbox Inbox
	log   logium.Logger
}

func NewService(log logium.Logger, inbox Inbox) *Service {
	return &Service{
		inbox: inbox,
		log:   log,
	}
}
