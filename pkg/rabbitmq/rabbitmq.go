package rabbitmq

import (
	"fmt"
	"go-email/config"
	"strconv"

	amqp "github.com/rabbitmq/amqp091-go"
)

func NewRabbitMQ(cfg *config.Config) (*amqp.Connection, error) {
	// net.JoinHostPort is not come handy here, because RabbitMQ host consists with user, passowrd, host, port
	//nolint
	connAddr := fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		cfg.Rabbit.User,
		cfg.Rabbit.Password,
		cfg.Rabbit.Host,
		strconv.Itoa(cfg.Rabbit.Port),
	)

	conn, err := amqp.Dial(connAddr)

	if err != nil {
		return nil, err
	}

	return conn, nil
}
