package rabbit

import (
	"log"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

type RabbitPublisher struct {
	connection *amqp.Connection
}

func NewRabbitPublisher() *RabbitPublisher {
	// create connection
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("failed to open connection %s", err.Error())
	}

	return &RabbitPublisher{
		connection: conn,
	}
}

// Shutdown will close the rabbit connection
func (r *RabbitPublisher) Shutdown() {
	r.connection.Close()
}

// Publish will send a given message onto a given queue or exchange
func (r *RabbitPublisher) Publish(queueName, exchange string, msg []byte, headers map[string]interface{}) error {
	// open a channel
	c, err := r.connection.Channel()
	if err != nil {
		log.Fatalf("failed to open channel: %s", err.Error())
	}
	defer c.Close()

	queue, err := c.QueueDeclarePassive(queueName, false, false, false, false, nil)
	if err != nil {
		log.Fatalf("failed to declare queue: %s", err.Error())
	}

	err = c.Publish(exchange, queueName, false, false, amqp.Publishing{
		Headers:     headers,
		ContentType: "application/json",
		Body:        msg,
	})

	if err != nil {
		return errors.Wrapf(err, "failed to publish message to exchange '%s' and queue '%s': %s", exchange, queue.Name, err.Error())
	}

	return nil
}
