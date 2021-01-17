package service

import "github.com/pkg/errors"

// Publisher defines the functions required to send messages to rabbit
type Publisher interface {
	Publish(queueName, exchange string, body []byte, headers map[string]interface{}) error
}

// Service handles sending messages to rabbit
type Service struct {
	publisher Publisher
}

// NewService will create and return a new service using the given dependancies
func NewService(publisher Publisher) *Service {
	return &Service{
		publisher: publisher,
	}
}

// SendMessage will send a given message to rabbit for requested number of times
func (p *Service) SendMessage(body []byte, headers map[string]interface{}, repeat int) error {
	if body == nil {
		return errors.New("no body provided")
	}

	for i := 0; i < repeat; i++ {
		err := p.publisher.Publish("test", "", body, headers)
		if err != nil {
			return errors.Wrapf(err, "error sending message %v", i)
		}
	}
	return nil
}
