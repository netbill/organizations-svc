package consumer

import (
	"context"

	"github.com/segmentio/kafka-go"
	"github.com/umisto/cities-svc/internal/messenger/contracts"
	"github.com/umisto/kafkakit/subscriber"
	"github.com/umisto/logium"
)

type Service struct {
	log       logium.Logger
	addr      []string
	callbacks callbacks
}

type callbacks interface {
	CreateAccount(ctx context.Context, event kafka.Message) error
	UpdateUsername(ctx context.Context, event kafka.Message) error
	UpdateProfile(ctx context.Context, event kafka.Message) error
}

func New(log logium.Logger, addr []string, callbacks callbacks) *Service {
	return &Service{
		addr:      addr,
		log:       log,
		callbacks: callbacks,
	}
}

func (s Service) Run(ctx context.Context) {
	sub := subscriber.New(s.addr, contracts.AccountsTopicV1, contracts.CitiesSvcGroup)

	s.log.Info("starting events consumer", "addr", s.addr)

	go func() {
		err := sub.Consume(ctx, func(m kafka.Message) (subscriber.HandlerFunc, bool) {
			et, ok := subscriber.Header(m, "event_type")
			if !ok {
				return nil, false
			}

			switch et {
			case contracts.AccountCreatedEvent:
				return s.callbacks.CreateAccount, true
			case contracts.AccountUsernameChangeEvent:
				return s.callbacks.UpdateUsername, true
			case contracts.ProfileUpdatedEvent:
				return s.callbacks.UpdateProfile, true
			default:
				return nil, false
			}
		})
		if err != nil {
			s.log.Warnf("accounts consumer stopped: %v", err)
		}
	}()
}
