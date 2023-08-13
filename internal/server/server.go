package server

import (
	"context"
	"fmt"
	config "go-email/config"
	delivery "go-email/internal/delivery/grpc"
	"go-email/internal/mailer"
	logger "go-email/pkg/logger"
	mail "go-email/pkg/mailer"
	pb "go-email/pkg/proto/email-service"
	rb "go-email/pkg/rabbitmq"
	"net"
	"net/http"

	echoSwagger "github.com/swaggo/echo-swagger"

	repository "go-email/internal/database"
	consumer "go-email/internal/delivery/rabbitmq"
	db "go-email/pkg/database"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"

	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var cfg, _ = config.LoadConfigFromEnv()

func runGrpcRest() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := pb.RegisterMailerServiceHandlerFromEndpoint(ctx, mux, fmt.Sprintf("localhost:%d", cfg.Grpc.Port), opts)

	if err != nil {
		panic(err)
	}

	log.Printf("server listening at %d", cfg.Gateway.Port)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.Gateway.Port), mux); err != nil {
		panic(err)
	}
}

func runMetricsAndSwagger() {
	e := echo.New()
	e.GET("/", echo.WrapHandler(promhttp.Handler()))

	// TODO:
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	log.Fatal(e.Start(fmt.Sprintf(":%d", cfg.Metrics.Port)))
}

func Run() {
	conf, err := config.LoadConfigFromEnv()

	if err != nil {
		log.Println("Error loading config")
	}

	// Init logger
	logger.InitLogger()

	// Settings for SMTP server
	d := mail.NewMailDialer(conf)
	mailer := mailer.NewMailer(d)

	// Init metrics and Swagger docs
	go runMetricsAndSwagger()

	// Init grpc-gateway rest
	go runGrpcRest()

	rabbitConnection, err := rb.NewRabbitMQ(conf)
	cons := consumer.NewConsumer(rabbitConnection, mailer, conf)

	if err != nil {
		log.Fatalf("Cannot connect to the rabbitmq %s\n", err.Error())
	}

	go func() {
		err := cons.Consume(conf.Rabbit.ConsumePool)

		if err != nil {
			log.Fatal(err.Error())
		}
	}()

	// Init database
	dbConn, err := db.NewDatabase(conf)
	repo := repository.NewRepository(dbConn)

	if err != nil {
		log.Fatal(err.Error())
	}

	// Implementing grpc server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", conf.Grpc.Port))
	if err != nil {
		log.Errorf("failed to start server %v", err)
	}

	log.WithFields(log.Fields{
		"port": lis.Addr(),
	}).Debug("Server started successfully")

	s := grpc.NewServer()
	pb.RegisterMailerServiceServer(s, delivery.NewServer(conf, mailer, repo))

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
