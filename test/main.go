package main

import (
	"fmt"
	"log"

	"github.com/MSSkowron/nats-service-sharing-environments/config"
	"github.com/MSSkowron/nats-service-sharing-environments/provider"
)

func main() {
	config, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", config)

	chatter := provider.NewChatter(provider.NewNATSProvider(config), config)

	err = chatter.Run()
	if err != nil {
		log.Fatal(err)
	}
}
