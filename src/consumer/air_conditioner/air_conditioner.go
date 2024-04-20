package main

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"os"
	"os/signal"
	"reflect"
	"strconv"
	"sync"
	"time"
)

const mutexLocked = 1

func MutexLocked(m *sync.Mutex) bool {
	state := reflect.ValueOf(m).Elem().FieldByName("state")
	return state.Int()&mutexLocked == mutexLocked
}

var mutex = sync.Mutex{}

func main() {
	var urls = "nats://admin:admin@localhost:4222"
	if len(os.Args) != 2 {
		fmt.Printf("Specify queue name as an argument!\n")
		os.Exit(1)
	}

	var queue = os.Args[1]
	var subject_turn_conditioner = queue + ".turn_conditioner"
	var subject_set_temperature = queue + ".set_temperature"
	temperature := 0

	nc, err := nats.Connect(urls)
	if err != nil {
		fmt.Println("Error during connection")
	}
	defer nc.Close()

	fmt.Printf("Connection set successfully\n")
	fmt.Printf("Queue name: %v\n", queue)
	fmt.Printf("Available subjects: %v, %v\n", subject_turn_conditioner, subject_set_temperature)

	nc.QueueSubscribe(subject_turn_conditioner, queue, func(msg *nats.Msg) {
		fmt.Printf("For subject: %v, from queue: %v, message: %v\n", subject_turn_conditioner, queue, string(msg.Data))
		if t, err := strconv.Atoi(string(msg.Data)); err == nil && t <= 30 && t >= 1 {
			mutex.Lock()
			fmt.Printf("Execution time for air conditioner has been set to %v mins.\n", t)
			time.Sleep(time.Second * time.Duration(t))
			fmt.Printf("Execution end!\n")
			mutex.Unlock()
		} else {
			fmt.Printf("Invalid execution time from system! Execution time has not been set! It is not working!\n")
		}
	})
	nc.Subscribe(subject_set_temperature, func(msg *nats.Msg) {
		fmt.Printf("For subject: %v, message: %v\n", subject_set_temperature, string(msg.Data))
		if t, err := strconv.Atoi(string(msg.Data)); err == nil {
			if MutexLocked(&mutex) {
				fmt.Printf("Cannot set temperature when air conditoiner is working!\n")
			} else {
				temperature = t
				fmt.Printf("Temperature for air conditioner has been set to %v C°.\n", temperature)
			}
		} else {
			fmt.Printf("Invalid temperature from system! Temperature has not been set! It is still: %v C°\n", temperature)
		}
	})
	nc.Flush()

	if nc.LastError() != nil {
		fmt.Printf("Error occured!\n")
	}

	fmt.Printf("Listening on [%v] and [%v], queue group [%v]\n", subject_turn_conditioner, subject_set_temperature, queue)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	nc.Drain()
	fmt.Printf("Exiting")
}
