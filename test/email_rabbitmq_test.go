package test

import (
	"go-email/config"
	"testing"

	rb "go-email/pkg/rabbitmq"

	"github.com/stretchr/testify/require"
)

var cfg, _ = config.LoadConfigFromEnv()

func rabbitPublishEmailToQueueForTests() {

}

func TestEmailRabbitMQ_RabbitMQConnection(t *testing.T) {
	conn, err := rb.NewRabbitMQ(cfg)
	defer conn.Close()

	require.NoError(t, err)
	require.Nil(t, err)
}

// TODO:
func TestEmailRabbitMQ_RabbitMQReadMessagesFromQueue(t *testing.T) {
	conn, _ := rb.NewRabbitMQ(cfg)
	defer conn.Close()
}
