package rabbitmq

import (
	"auth-backend/internal/config"
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

type RabbitMQ struct {
	cfg        *config.RabbitConnectionConfig
	Connection *amqp.Connection
	Channel    *amqp.Channel
}

func NewRabbitMQ(config *config.RabbitConnectionConfig) *RabbitMQ {
	return &RabbitMQ{
		cfg:        config,
		Connection: nil,
		Channel:    nil,
	}
}

func (r *RabbitMQ) Connect() error {
	var err error
	r.Connection, err = amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d/", r.cfg.Username, r.cfg.Password, r.cfg.Host, r.cfg.Port))
	if err != nil {
		return err
	}
	r.Channel, err = r.Connection.Channel()
	if err != nil {
		return err
	}

	_, err = r.Channel.QueueDeclare("logger", false, false, false, false, nil)
	fmt.Println("Successfully connected and created queue")
	return err
}

func (r *RabbitMQ) Consume() {
	messages, err := r.Channel.Consume("logger", "", true, false, false, false, nil)

	if err != nil {
		log.Fatalf("failed to register a consumer. Error: %s", err)
	}

	for msg := range messages {
		log.Printf("saved to db: %s", string(msg.Body))
		// todo: log to mongo
	}
}

func (r *RabbitMQ) Close() {
	_ = r.Channel.Close()
	_ = r.Connection.Close()
}
