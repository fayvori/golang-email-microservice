package rabbitmq

import (
	"context"
	"encoding/json"
	"go-email/config"
	"go-email/internal/mailer"
	"go-email/internal/models"
	"go-email/pkg/constants"

	repo "go-email/internal/database"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel"
)

// Load config from env variable
var cfg = config.LoadConfigFromEnv()

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
	repo   *repo.Resository
	cfg    *config.Config
}

func NewConsumer(conn *amqp.Connection, mailer *mailer.Mailer, repo *repo.Resository, cfg *config.Config) *Consumer {
	return &Consumer{conn: conn, mailer: mailer, repo: repo, cfg: cfg}
}

// Function for creating new channel in `RabbitMQ`
func (c *Consumer) createChannel() (*amqp.Channel, error) {
	channel, err := c.conn.Channel()

	if err != nil {
		return nil, err
	}

	_, err = channel.QueueDeclare(
		c.cfg.Rabbit.QueueName, // name
		false,                  // durable
		false,                  // delete when unused
		false,                  // exclusive
		false,                  // no-wait
		nil,                    // arguments
	)

	if err != nil {
		return nil, err
	}

	return channel, nil
}

func (c *Consumer) Consume(poolSize int) error {
	channel, err := c.createChannel()

	if err != nil {
		return err
	}

	defer channel.Close()

	var forever chan struct{}

	messages, err := channel.Consume(
		c.cfg.Rabbit.QueueName, // queue
		"",                     // consumer
		true,                   // auto-ack
		false,                  // exclusive
		false,                  // no-local
		false,                  // no-wait
		nil,                    // args
	)

	if err != nil {
		return err
	}

	tp := otel.Tracer(constants.TRACER_NAME)

	for i := 0; i < poolSize; i++ {
		for msg := range messages {
			rootContext, rootSpan := tp.Start(context.Background(), "emails-rabbit-root")

			email := &models.Email{}
			err := json.Unmarshal(msg.Body, &email)

			// Set root email from cfg
			email.From = cfg.SMTP.User

			if err != nil {
				return err
			}

			_, span := tp.Start(rootContext, "emails-rabbit-send-emails")

			if err := c.mailer.SendEmails(email); err != nil {
				messagesConsumedFailure.Inc()
				return err
			}

			span.End()

			_, span = tp.Start(rootContext, "emails-rabbit-save-emails")

			if err := c.repo.CreateEmail(email); err != nil {
				messagesConsumedFailure.Inc()
				return err
			}

			span.End()
			rootSpan.End()

			messagesCousumedSuccessfully.Inc()
		}
	}
	<-forever

	messagesConsumedFailure.Inc()

	return nil
}
