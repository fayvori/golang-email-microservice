package test

import (
	"go-email/config"
	"testing"

	rb "go-email/pkg/rabbitmq"

	"github.com/stretchr/testify/require"
)

var cfg = config.LoadConfigFromEnv()

func TestEmailRabbitMQ_RabbitMQConnection(t *testing.T) {
	conn, err := rb.NewRabbitMQ(cfg)
	//nolint
	defer conn.Close()

	require.NoError(t, err)
	require.Nil(t, err)
}
