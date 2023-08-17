package server

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"go-email/config"
	pb "go-email/pkg/proto/email-service"

	delivery "go-email/internal/delivery/grpc"
	"go-email/internal/mailer"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"

	"github.com/labstack/echo/v4"
	repository "go-email/internal/database"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var cfg = config.LoadConfigFromEnv()

func runGrpc(mailer *mailer.Mailer, repo *repository.Resository) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Grpc.Port))

	if err != nil {
		log.WithFields(log.Fields{
			"message": "Unable to start gRPC server",
		}).Errorf("Cannot start gRPC server %s", err.Error())
	}

	log.WithFields(log.Fields{
		"message": "gRPC server started successfully",
	}).Printf("gRPC server listening at port %d", cfg.Grpc.Port)

	s := grpc.NewServer()
	pb.RegisterMailerServiceServer(s, delivery.NewServer(cfg, mailer, repo))

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %s", err.Error())
	}
}

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

	log.WithFields(log.Fields{
		"message": "gRPC REST server started successfully",
	}).Printf("gRPC REST server listening at port %d", cfg.Gateway.Port)

	//nolint
	if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.Gateway.Port), mux); err != nil {
		panic(err)
	}
}

func runMetrics() {
	echoServer := echo.New()

	// Hide initial echo banner
	echoServer.HideBanner = true
	echoServer.HidePort = true

	log.WithFields(log.Fields{
		"message": "Metrics server started successfully",
	}).Printf("Metrics server listening at port %d", cfg.Metrics.Port)

	echoServer.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	log.Fatal(echoServer.Start(fmt.Sprintf(":%d", cfg.Metrics.Port)))
}
