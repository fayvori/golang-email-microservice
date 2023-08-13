package test

import (
	config "go-email/config"
	repository "go-email/internal/database"
	delivery "go-email/internal/delivery/grpc"
	"go-email/internal/mailer"
	db "go-email/pkg/database"
	mail "go-email/pkg/mailer"
	pb "go-email/pkg/proto/email-service"

	"context"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"net"
	"testing"
)

var conf, _ = config.LoadConfigFromEnv()

func server(ctx context.Context) (pb.MailerServiceClient, func()) {
	buffer := 1024 * 1024
	listener := bufconn.Listen(buffer)

	d := mail.NewMailDialer(conf)
	mailer := mailer.NewMailer(d)

	dbConn, _ := db.NewDatabase(conf)
	repo := repository.NewRepository(dbConn)

	s := grpc.NewServer()
	pb.RegisterMailerServiceServer(s, delivery.NewServer(conf,
		mailer,
		repo,
	))
	go func() {
		if err := s.Serve(listener); err != nil {
			panic(err)
		}
	}()

	conn, _ := grpc.DialContext(ctx, "", grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}), grpc.WithInsecure(), grpc.WithBlock())

	client := pb.NewMailerServiceClient(conn)

	return client, s.Stop
}

func TestEmailGrpc_SendEmails(t *testing.T) {
	ctx := context.Background()
	client, closer := server(ctx)

	defer closer()

	email := pb.EmailRequest{
		To:          []string{"alexemailtestingtwo@yahoo.com"},
		ContentType: "text/plain",
		Subject:     "grpc test",
		Body:        "testing",
	}

	_, err := client.SendEmails(ctx, &email)

	assert.NoError(t, err)
	assert.Nil(t, err)
}
