package main

import (
	"fmt"
	"os"
	"os/signal"
	"github.com/nats-io/nats.go"
)

func main() {
	var urls = "nats://admin:admin@localhost:4222"
	var subject = "subject"

	nc, err := nats.Connect(urls)
	if err != nil {
		fmt.Println("Error during connection")
	}
	defer nc.Close()

	nc.Subscribe(subject, func(msg *nats.Msg) {
		fmt.Printf("For subject: %v, message: %v\n", subject, string(msg.Data))
	})
	nc.Flush()

	if nc.LastError() != nil {
		fmt.Printf("Error occured!\n")
	}

	fmt.Printf("Listening on subject: %s\n", subject)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	nc.Drain()
	fmt.Printf("Exiting")
}