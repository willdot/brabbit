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

// Request is a request to send a message to rabbit
type Request struct {
	Queue    string
	Exchange string
	Body     []byte
	Headers  map[string]interface{}
	Repeat   int
}

// SendMessage will send a given message to rabbit for requested number of times
func (p *Service) SendMessage(r Request) error {
	if r.Body == nil {
		return errors.New("no body provided")
	}

	for i := 0; i < r.Repeat; i++ {
		err := p.publisher.Publish(r.Queue, r.Exchange, r.Body, r.Headers)
		if err != nil {
			return errors.Wrapf(err, "error sending message %v", i)
		}
	}
	return nil
}
