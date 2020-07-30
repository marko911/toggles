package message

import (
	"github.com/nats-io/nats.go"
)

// Subscription holds ref to topic subscription
type Subscription interface {
	Unsubscribe() error
}

// Service is a messaging service
type Service interface {
	Publish(subj string, data []byte) error
	Subscribe(subj string, cb func(m interface{})) Subscription
	ChanSubscribe(subj string, c chan Msg) error
}

// NatsClient is the nats client impl
type NatsClient interface {
	Publish(subj string, data []byte) error
	Subscribe(subj string, cb nats.MsgHandler) (*nats.Subscription, error)
	ChanSubscribe(subj string, c chan *nats.Msg) (*nats.Subscription, error)
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

func (s *natsService) Publish(subj string, data []byte) error {
	err := s.c.Publish(subj, data)
	return err
}

func (s *natsService) Subscribe(subj string, cb func(m interface{})) Subscription {
	sub, _ := s.c.Subscribe(subj, func(m *nats.Msg) {
		cb(m)
	})
	return sub
}

func (s *natsService) ChanSubscribe(subj string, c chan Msg) error {
	ch := make(chan *nats.Msg, 64)
	go func() {
		defer close(ch)
		for {
			msg, ok := <-ch
			if !ok {
				break
			}
			c <- msg
		}
	}()
	_, err := s.c.ChanSubscribe(subj, ch)
	return err
}
