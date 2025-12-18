package pubsub

import (
	"context"
	"fmt"
	"github.com/ipfs/go-ipfs-api"
)

type EventSubscriber struct {
	api *api.Shell
}

func NewEventSubscriber(apiURL string) *EventSubscriber {
	return &EventSubscriber{
		api: api.NewShell(apiURL),
	}
}

func (s *EventSubscriber) Subscribe(ctx context.Context, topic string, handler func([]byte)) error {
	sub, err := s.api.PubSubSubscribe(topic)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-sub:
				handler(msg.Data)
			}
		}
	}()

	return nil
}

func PublishEvent(apiURL string, topic string, data []byte) error {
	api := api.NewShell(apiURL)
	return api.PubSubPublish(topic, string(data))
}

