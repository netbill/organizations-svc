package consumer

import (
	"context"

	"github.com/segmentio/kafka-go"
	"github.com/umisto/agglomerations-svc/internal/messenger/contracts"
	"github.com/umisto/kafkakit/subscriber"
	"github.com/umisto/logium"
)

type Service struct {
	log       logium.Logger
	addr      []string
	callbacks callbacks
}

type callbacks interface {
	AccountCreated(ctx context.Context, event kafka.Message) error
	AccountDeleted(ctx context.Context, event kafka.Message) error
	AccountUsernameChanged(ctx context.Context, event kafka.Message) error
	ProfileUpdated(ctx context.Context, event kafka.Message) error
}

func New(log logium.Logger, addr []string, callbacks callbacks) *Service {
	return &Service{
		addr:      addr,
		log:       log,
		callbacks: callbacks,
	}
}

func (s Service) Run(ctx context.Context) {
	sub := subscriber.New(s.addr, contracts.AccountsTopicV1, contracts.agglomerationsSvcGroup)

	s.log.Info("starting events consumer", "addr", s.addr)

	go func() {
		err := sub.Consume(ctx, func(m kafka.Message) (subscriber.HandlerFunc, bool) {
			et, ok := subscriber.Header(m, "event_type")
			if !ok {
				return nil, false
			}

			switch et {
			case contracts.AccountCreatedEvent:
				return s.callbacks.AccountCreated, true
			case contracts.AccountDeletedEvent:
				return s.callbacks.AccountDeleted, true
			case contracts.AccountUsernameChangeEvent:
				return s.callbacks.AccountUsernameChanged, true
			case contracts.ProfileUpdatedEvent:
				return s.callbacks.ProfileUpdated, true
			default:
				return nil, false
			}
		})
		if err != nil {
			s.log.Warnf("accounts consumer stopped: %v", err)
		}
	}()
}
