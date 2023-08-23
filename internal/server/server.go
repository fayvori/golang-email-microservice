package server

import (
	"context"
	config "go-email/config"
	logger "go-email/pkg/logger"
	"go-email/pkg/tracer"
	"time"

	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"

	mailer "go-email/internal/mailer"
	mail "go-email/pkg/mailer"

	consumer "go-email/internal/delivery/rabbitmq"
	rb "go-email/pkg/rabbitmq"

	repository "go-email/internal/database"
	db "go-email/pkg/database"
)

func Run() {
	conf := config.LoadConfigFromEnv()

	// Init logger
	logger.InitLogger()

	// Setup for SMTP server
	d := mail.NewMailDialer(conf)
	mailer := mailer.NewMailer(d)

	// Init database
	dbConn, err := db.NewDatabase(conf)
	repo := repository.NewRepository(dbConn)

	if err != nil {
		log.WithFields(log.Fields{
			"message": "Unable to connect to the Database",
		}).Infof("Cannot connect to the Database %s", err.Error())
	}
	// Init RabbitMQ connection
	rabbitConnection, err := rb.NewRabbitMQ(conf)
	cons := consumer.NewConsumer(rabbitConnection, mailer, repo, conf)

	if err != nil {
		log.WithFields(log.Fields{
			"message": "Unable to connect to the RabbitMQ server",
		}).Infof("Cannot connect to the RabbitMQ %s", err.Error())
	}

	go func() {
		err := cons.Consume(conf.Rabbit.ConsumePool)

		if err != nil {
			log.WithFields(log.Fields{
				"message": "Unable to consume messages from RabbitMQ",
			}).Infof("Cannot receive messages from RabbitMQ %s", err.Error())
		}
	}()

	// Init tracer provider
	tracerProvider := tracer.TracerProvider(conf)

	otel.SetTracerProvider(tracerProvider)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Cleanly shutdown and flush telemetry when the application exits.
	defer func(ctx context.Context) {
		// Do not make the application hang when it is shutdown.
		ctx, cancel = context.WithTimeout(ctx, time.Second*5)
		defer cancel()

		if err := tracerProvider.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}
	}(ctx)

	// Init grpc-gateway rest
	go runGrpcRest()

	// Init grpc server
	go runGrpc(mailer, repo)

	// Init metrics and Swagger docs
	runMetricsAndHealth()
}
