package publisher

import (
	"time"

	"github.com/netbill/eventbox"
)

type WriterConfig struct {
	RequiredAcks string        `mapstructure:"required_acks"`
	Compression  string        `mapstructure:"compression"`
	Balancer     string        `mapstructure:"balancer"`
	BatchSize    int           `mapstructure:"batch_size"`
	BatchTimeout time.Duration `mapstructure:"batch_timeout"`
	DialTimeout  time.Duration `mapstructure:"dial_timeout"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
}

type Publisher struct {
	identity string
	outbox   eventbox.Outbox
	producer *eventbox.Producer
}

func New(
	identity string,
	outbox eventbox.Outbox,
	producer *eventbox.Producer,
) *Publisher {
	return &Publisher{
		identity: identity,
		outbox:   outbox,
		producer: producer,
	}
}
