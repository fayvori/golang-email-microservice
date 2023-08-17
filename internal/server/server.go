package server

import (
	log "github.com/sirupsen/logrus"
	config "go-email/config"
	logger "go-email/pkg/logger"

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

	// Init RabbitMQ connection
	rabbitConnection, err := rb.NewRabbitMQ(conf)
	cons := consumer.NewConsumer(rabbitConnection, mailer, conf)

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

	// Init database
	dbConn, err := db.NewDatabase(conf)
	repo := repository.NewRepository(dbConn)

	if err != nil {
		log.WithFields(log.Fields{
			"message": "Unable to connect to the Database",
		}).Infof("Cannot connect to the Database %s", err.Error())
	}

	// Init grpc-gateway rest
	go runGrpcRest()

	// Init metrics and Swagger docs
	go runMetrics()

	// Init grpc server
	runGrpc(mailer, repo)
}
