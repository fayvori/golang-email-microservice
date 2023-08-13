package rabbitmq

import (
	"encoding/json"
	"go-email/config"
	"go-email/internal/mailer"
	"go-email/internal/models"
	"log"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	amqp "github.com/rabbitmq/amqp091-go"
)

// Load config from env variable
var cfg, _ = config.LoadConfigFromEnv()

// Custom metrics for `Prometheus`
var (
	messagesCousumedSuccessfully = promauto.NewCounter(prometheus.CounterOpts{
		Name: "rabbitmq_emails_sended_successfully_total",
		Help: "Count of successfully sended emails througth rabbitmq",
	})

	messagesConsumedFailure = promauto.NewCounter(prometheus.CounterOpts{
		Name: "rabbitmq_emails_sended_failure_total",
		Help: "Count of failure sended emails througth rabbitmq",
	})
)

type Consumer struct {
	conn   *amqp.Connection
	mailer *mailer.Mailer
	cfg    *config.Config
}

func NewConsumer(conn *amqp.Connection, mailer *mailer.Mailer, cfg *config.Config) *Consumer {
	return &Consumer{conn: conn, mailer: mailer, cfg: cfg}
}

// Function for creating new channel in `RabbitMQ`
func (c *Consumer) createChannel() (*amqp.Channel, error) {
	ch, err := c.conn.Channel()

	if err != nil {
		return nil, err
	}

	_, err = ch.QueueDeclare(
		c.cfg.Rabbit.QueueName, // name
		false,                  // durable
		false,                  // delete when unused
		false,                  // exclusive
		false,                  // no-wait
		nil,                    // arguments
	)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return ch, nil
}

func (c *Consumer) Consume(poolSize int) error {
	ch, err := c.createChannel()

	if err != nil {
		log.Fatal(err)
	}

	defer ch.Close()

	var forever chan struct{}

	messages, err := ch.Consume(
		c.cfg.Rabbit.QueueName, // queue
		"",                     // consumer
		true,                   // auto-ack
		false,                  // exclusive
		false,                  // no-local
		false,                  // no-wait
		nil,                    // args
	)

	for i := 0; i < poolSize; i++ {
		for msg := range messages {
			email := &models.Email{}
			err := json.Unmarshal(msg.Body, &email)

			// Set root email from cfg
			email.From = cfg.Smtp.User

			if err != nil {
				log.Fatal(err)
			}

			if err := c.mailer.SendEmails(email); err != nil {
				messagesConsumedFailure.Inc()
				return err
			}

			messagesCousumedSuccessfully.Inc()
		}
	}
	<-forever

	messagesConsumedFailure.Inc()

	return nil
}
