package provider

import (
	"fmt"
	"time"

	"github.com/MSSkowron/nats-service-sharing-environments/config"

	"github.com/nats-io/nats.go"
)

type NATSProvider struct {
	Config config.Config
	JS     nats.JetStreamContext
}

func NewNATSProvider(config config.Config) *NATSProvider {
	return &NATSProvider{
		Config: config,
	}
}

func (p *NATSProvider) Connect(addr string) error {
	nc, err := nats.Connect(addr)
	if err != nil {
		return fmt.Errorf("error connecting to node at %s: %v\n", addr, err)
	}

	js, _ := nc.JetStream()

	if p.Config.AddStream {
		if _, err := js.AddStream(&nats.StreamConfig{
			Name:     p.Config.Subject,
			Subjects: []string{p.Config.Subject},
		}); err != nil {
			return fmt.Errorf("error adding stream %w\n", err)
		}
	}

	p.JS = js

	return nil
}

func (p *NATSProvider) Publish(msg []byte) error {
	_, err := p.JS.Publish(p.Config.Subject, msg)
	return err
}

func (p *NATSProvider) Subscribe(since time.Time, handler MsgHandler) error {
	sub, err := p.JS.SubscribeSync(p.Config.Subject, nats.OrderedConsumer(), nats.StartTime(since))
	if err != nil {
		return fmt.Errorf("error subscribing to JetStream %v", err)
	}

	for {
		msg, err := sub.NextMsg(5 * time.Minute)
		if err != nil {
			return fmt.Errorf("error receiving message from JetStream: %w", err)

		}

		meta, err := msg.Metadata()
		if err != nil {
			return fmt.Errorf("retrieving message (not a jetstream?): %w", err)
		}

		if err := handler(msg.Data, meta.Timestamp); err != nil {
			return err
		}
	}
}
