package main

import (
	"os"
	"github.com/nats-io/nats.go"
	"fmt"
)


func main() {
	var urls = "nats://admin:admin@localhost:4222"
	var subject = "subject"
	var message = "message"
	var reply *string = nil

	nc, err := nats.Connect(urls)
	if err != nil {
		fmt.Println("Error during connection")
		os.Exit(1)
	}
	defer nc.Close()

	if reply != nil && *reply != "" {
		nc.PublishRequest(subject, *reply, []byte(message))
	} else {
		nc.Publish(subject, []byte(message))
	}

	nc.Flush()

	if nc.LastError() != nil {
		fmt.Printf("Error occured!\n")
		return
	}
	fmt.Printf("Message '%v' has been sent for subject: '%s'\n", message, subject)
}