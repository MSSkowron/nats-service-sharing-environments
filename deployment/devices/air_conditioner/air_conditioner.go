package main

import (
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
)

const (
	url         = "nats://admin:admin@nats3:4222"
	queueName   = "air_conditioners_queue"
	subjectTurn = "turn_conditioner"
	subjectSet  = "set_temperature"
	defaultTemp = 20
	maxTemp     = 30
	minTemp     = 0
)

var (
	temperature = defaultTemp
	mu          sync.Mutex
)

func main() {
	name := os.Getenv("AIR_CONDITIONER_NAME")

	nc, err := nats.Connect(url)
	if err != nil {
		panic(err)
	}
	defer nc.Close()

	js, err := nc.JetStream()
	if err != nil {
		panic(err)
	}

	if _, err := js.QueueSubscribe(subjectTurn, queueName, handleTurn); err != nil {
		panic(err)
	}
	if _, err := js.QueueSubscribe(subjectSet, queueName, handleSet); err != nil {
		panic(err)
	}

	nc.Flush()
	if err := nc.LastError(); err != nil {
		panic(err)
	}

	slog.Info("Waiting for messages...", "name", name, "subjects", "subject_turn_conditioner, subject_set_temperature", "queue", queueName)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	nc.Drain()

	slog.Info("Exiting...")
}

func handleTurn(msg *nats.Msg) {
	slog.Info("Received turn message", "message", string(msg.Data))

	if t, err := strconv.Atoi(string(msg.Data)); err == nil && t <= 30 && t >= 1 {
		mu.Lock()
		slog.Info("Execution time for air conditioner has been set", "time", t, "unit", "mins")
		time.Sleep(time.Second * time.Duration(t))
		slog.Info("Execution end!")
		mu.Unlock()
	} else {
		slog.Info("Invalid execution time from system! Execution time has not been set!")
	}
}

func handleSet(msg *nats.Msg) {
	slog.Info("Received set message", "message", string(msg.Data))

	if t, err := strconv.Atoi(string(msg.Data)); err == nil {
		if mu.TryLock() {
			temperature = t
			slog.Info("Temperature for air conditioner has been set", "temperature", temperature, "unit", "C°")
			mu.Unlock()
		} else {
			slog.Info("Cannot set temperature when air conditioner is working!")
		}
	} else {
		slog.Info("Invalid temperature from system! Temperature has not been set!", "temperature", temperature, "unit", "C°")
	}
}
