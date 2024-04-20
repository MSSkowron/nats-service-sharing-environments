package main

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"os"
	"strconv"
)

func send_message(subject string, message string, nc *nats.Conn){
	nc.Publish(subject, []byte(message))

	nc.Flush()
	if nc.LastError() != nil {
		fmt.Printf("Error occured!\n")
		return
	}
	fmt.Printf("Message '%v' has been sent!\n", message)
}

func main() {
	var urls = "nats://admin:admin@localhost:4222" //nats.DefaultURL
	if len(os.Args) != 4 {
		fmt.Printf("Required three arguments!\n")
		fmt.Printf("Specify queue name and action for air conditioner [set_temperature, turn_conditioner] and value!\n")
		fmt.Printf("Temp value for set_temperature can be from 15 to 32\n")
		fmt.Printf("Time value for turn_conditioner can be from 1 to 30\n")
		os.Exit(1)
	}

	if os.Args[2] == "set_temperature" {
		if t, err := strconv.Atoi(string(os.Args[3])); err != nil || !(t <= 32 && t >= 15) {
			fmt.Printf("Temperature is not a proper value! Please provide an integer from 15 to 32!\n")
			os.Exit(1)
		}

	} else if os.Args[2] == "turn_conditioner" {
		if t, err := strconv.Atoi(string(os.Args[3])); err != nil || !(t <= 30 && t >= 1) {
			fmt.Printf("Time is not a proper value! Please provide an integer from 1 to 30!\n")
			os.Exit(1)
		}

	} else {
		fmt.Printf("Action for air conditioner is not proper! Please choose 'set_temperature' or 'turn_conditioner'\n")
		os.Exit(1)
	}

	var subject = os.Args[1] + "." + os.Args[2]

	nc, err := nats.Connect(urls)
	if err != nil {
		fmt.Println("Error during connection")
		os.Exit(1)
	}
	defer nc.Close()
	fmt.Printf("Sending on subject: %v\n", subject)

	send_message(subject, os.Args[3], nc)

}
