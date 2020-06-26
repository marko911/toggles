package message

import (
	"fmt"

	"github.com/nats-io/nats.go"
)

// Subscription holds ref to topic subscription
type Subscription interface {
	Unsubscribe() error
}

// Handler is a generic subscription callback
type Handler interface{}

// Service is a messaging service
type Service interface {
	Publish(subj string, v interface{}) error
	Subscribe(subj string, cb Handler) Subscription
	// ChanSubscribe(subj string, c chan Msg) error
}

// NatsClient is the nats client impl
type NatsClient interface {
	Publish(subj string, v interface{}) error
	Subscribe(subj string, cb nats.Handler) (*nats.Subscription, error)
	// ChanSubscribe(subj string, c chan *nats.Msg) (*nats.Subscription, error)
}

// Msg is a generic message
type Msg interface {
	Respond(data []byte) error // or a real type, whatever
}

type natsService struct {
	c NatsClient
}

// NewNatsService returns a messenger  instance
func NewNatsService(c NatsClient) Service {
	return &natsService{c}
}

func (s *natsService) Publish(subj string, data interface{}) error {
	err := s.c.Publish(subj, data)
	return err
}

func (s *natsService) Subscribe(subj string, cb Handler) Subscription {
	sub, _ := s.c.Subscribe(subj, cb)
	fmt.Println("subscribininggg!", subj)
	return sub
}
