package provider

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/MSSkowron/nats-service-sharing-environments/config"
	"github.com/google/uuid"
)

type Message struct {
	ProviderID string `json:"providerID"`
	ID         int    `json:"id"`
	From       string `json:"name"`
	Text       string `json:"text"`
}

type MsgHandler func(data []byte, timestamp time.Time) error

type Provider interface {
	Connect(addr string) error
	Publish(message []byte) error
	Subscribe(since time.Time, handler MsgHandler) error
}

type IDProvider struct {
	id string
	Provider
}

type Chatter struct {
	provider *IDProvider
	config   config.Config
}

func NewChatter(p Provider, c config.Config) *Chatter {
	return &Chatter{
		provider: &IDProvider{
			id:       uuid.NewString(),
			Provider: p,
		},
		config: c,
	}
}

func (c *Chatter) Run() error {
	if err := c.provider.Connect(c.config.ServerAddress); err != nil {
		return err
	}

	pubC := make(chan error)
	subC := make(chan error)
	p, s := false, false

	if c.config.Publish {
		p = true

		go func() {
			pubC <- c.publish()
		}()
	}

	if c.config.Subscribe {
		s = true

		go func() {
			subC <- c.subscribe()
		}()
	}

	for {
		if !p && !s {
			break
		}

		select {
		case err := <-pubC:
			p = false

			if err != nil {
				return fmt.Errorf("publisher exited with error: %w", err)
			}

			log.Println("publisher finished")
		case err := <-subC:
			s = false

			if err != nil {
				return fmt.Errorf("subscriber exited with error: %w", err)
			}

			log.Println("subscriber finished")
		}
	}

	return nil
}

func (c *Chatter) publish() error {
	log.Printf("Publishing messages as %s\n", c.config.Username)
	if c.config.Interactive {
		log.Println("Running in interactive mode")
	}

	for id := 0; ; id++ {
		msg := Message{
			ID:   id,
			From: c.config.Username,
		}

		// read message from console
		if c.config.Interactive {
			reader := bufio.NewReader(os.Stdin)
			input, err := reader.ReadString('\n')
			if err != nil {
				return err
			}

			msg.Text = strings.TrimSuffix(input, "\n")
		}

		bytes, err := json.Marshal(msg)
		if err != nil {
			return err
		}

		if err := c.provider.Publish(bytes); err != nil {
			return err
		}

		if !c.config.Interactive {
			fmt.Printf("PUB %s: %d %s\n", c.config.Username, msg.ID, msg.Text)
			time.Sleep(time.Duration(c.config.Timeout) * time.Millisecond)
		}
	}
}

func (c *Chatter) subscribe() error {
	since, err := c.readTimestamp()
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Printf("error reading timestamp: %v", err)
	}

	handler := func(data []byte, timestamp time.Time) error {
		msg := Message{}

		if err := json.Unmarshal(data, &msg); err != nil {
			return fmt.Errorf("error unmarshalling message: %w", err)
		}

		if !c.config.ListenSelf && (msg.ProviderID == c.provider.id) {
			return nil
		}

		if err := c.saveTimestamp(timestamp); err != nil {
			log.Printf("error saving timestamp: %v", err)
		}

		if c.config.Interactive {
			fmt.Printf("	%s: %s\n", msg.From, msg.Text)
		} else {
			fmt.Printf("	SUB %s: %d %s\n", msg.From, msg.ID, msg.Text)
		}

		return nil
	}

	log.Printf("Subscribing to messages since %s", since.String())

	return c.provider.Subscribe(since, handler)
}

func (c *Chatter) readTimestamp() (time.Time, error) {
	bytes, err := os.ReadFile(c.fn())
	if err != nil {
		return time.Time{}, err
	}

	nanos, err := strconv.ParseInt(string(bytes), 10, 64)
	if err != nil {
		return time.Time{}, err
	}

	return time.Unix(0, nanos), nil
}

func (c *Chatter) saveTimestamp(t time.Time) error {
	t = t.Add(1)
	return os.WriteFile(c.fn(), []byte(strconv.Itoa(int(t.UnixNano()))), 0777)
}

func (c *Chatter) fn() string {
	return fmt.Sprintf("state_%s.txt", c.config.Username)
}
