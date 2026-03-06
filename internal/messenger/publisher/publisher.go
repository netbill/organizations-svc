package publisher

import (
	"github.com/netbill/eventbox"
)

type Publisher struct {
	identity string
	outbox   eventbox.Outbox
}

func New(
	identity string,
	outbox eventbox.Outbox,
) *Publisher {
	return &Publisher{
		identity: identity,
		outbox:   outbox,
	}
}
