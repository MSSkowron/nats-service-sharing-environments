package main

import (
	"log/slog"
	"os"
	"os/signal"

	"github.com/nats-io/nats.go"
)

const url = "nats://admin:admin@nats3:4222"

func main() {
	nc, err := nats.Connect(url)
	if err != nil {
		panic(err)
	}
	defer nc.Close()

	js, err := nc.JetStream()
	if err != nil {
		panic(err)
	}

	// PUBLISHER
	if _, err := js.Subscribe("publishers.>", func(msg *nats.Msg) {
		slog.Info("PUBLISHER", "MESSAGE", string(msg.Data))
	}); err != nil {
		panic(err)
	}

	// AIR CONDTITIONER
	if _, err := js.QueueSubscribe("airconditionersturn", "airconditionersqueue", func(msg *nats.Msg) {
		slog.Info("AIRCONDITIONER", "MESSAGE", string(msg.Data))
	}); err != nil {
		panic(err)
	}
	if _, err := js.Subscribe("airconditionersset", func(msg *nats.Msg) {
		slog.Info("AIRCONDITIONER", "MESSAGE", string(msg.Data))
	}); err != nil {
		panic(err)
	}

	// FRIDGE
	if _, err := js.Subscribe("fridges.>", func(msg *nats.Msg) {
		slog.Info("FRIDGE", "MESSAGE", string(msg.Data))
	}); err != nil {
		panic(err)
	}

	// FURNANCE
	if _, err := nc.Subscribe("furnances.>", func(msg *nats.Msg) {
		slog.Info("FURNANCE", "MESSAGE", string(msg.Data))
	}); err != nil {
		panic(err)
	}

	// LIGHT
	if _, err := js.Subscribe("lights.>", func(msg *nats.Msg) {
		slog.Info("LIGHT", "MESSAGE", string(msg.Data))
	}); err != nil {
		panic(err)
	}

	nc.Flush()
	if err := nc.LastError(); err != nil {
		panic(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	_ = nc.Drain()

	slog.Info("Exiting...")
}
