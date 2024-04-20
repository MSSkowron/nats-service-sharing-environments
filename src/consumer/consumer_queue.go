package main

import (
	"os"
	"os/signal"
	"fmt"
	"github.com/nats-io/nats.go"
)

func main() {
	var urls = "nats://admin:admin@localhost:4222"
	var subject = "subject"
	var queue = "queue"

	nc, err := nats.Connect(urls)
	if err != nil {
		fmt.Println("Error during connection")
	}
	defer nc.Close()

	nc.QueueSubscribe(subject, queue, func(msg *nats.Msg) {
		fmt.Printf("For subject: %v, from queue: %v, message: %v\n", subject, queue, string(msg.Data))
	})
	nc.Flush()

	if nc.LastError() != nil {
		fmt.Printf("Error occured!\n")
	}

	fmt.Printf("Listening on subject %v and queue: %v\n", subject, queue)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	nc.Drain()
	fmt.Printf("Exiting")
}