package delivery

import (
	"context"
	"errors"
	"go-email/config"
	repo "go-email/internal/database"
	"go-email/internal/mailer"
	"go-email/internal/models"
	"go-email/internal/validator"
	pb "go-email/pkg/proto/email-service"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	emailsSuccess = promauto.NewCounter(prometheus.CounterOpts{
		Name: "grpc_emails_sended_successfully_total",
		Help: "Successfully sended emails count",
	})

	emailsFailure = promauto.NewCounter(prometheus.CounterOpts{
		Name: "grpc_emails_sended_failure_total",
		Help: "Failed emails count",
	})

	emailsSavedSuccessfully = promauto.NewCounter(prometheus.CounterOpts{
		Name: "grpc_emails_saved_to_database_successfully_total",
		Help: "Count of successfully saved to database emails",
	})

	emailsSavedFailure = promauto.NewCounter(prometheus.CounterOpts{
		Name: "grpc_emails_saved_to_database_failure_total",
		Help: "Count of failure saved to database emails",
	})
)

type Server struct {
	pb.UnimplementedMailerServiceServer
	mailer *mailer.Mailer
	cfg    *config.Config
	repo   *repo.Resository
}

func NewServer(cfg *config.Config, mailer *mailer.Mailer, repo *repo.Resository) *Server {
	return &Server{cfg: cfg, mailer: mailer, repo: repo}
}

// context needed for implementing grpc interface, but linter sees it as unused
//nolint
func (s *Server) SendEmails(ctx context.Context, r *pb.EmailRequest) (*pb.EmailResponse, error) {
	email := &models.Email{
		From:        s.cfg.SMTP.User,
		To:          r.GetTo(),
		Body:        string(r.GetBody()),
		Subject:     r.GetSubject(),
		ContentType: r.GetContentType(),
	}

	for _, receiver := range r.GetTo() {
		if !validator.ValidateEmail(receiver) {
			return nil, errors.New("Unable to validate email")
		}
	}

	if err := s.mailer.SendEmails(email); err != nil {
		emailsFailure.Inc()
		return nil, err
	}

	if err := s.repo.CreateEmail(email); err != nil {
		emailsSavedFailure.Inc()
		return nil, err
	}

	emailsSuccess.Inc()
	emailsSavedSuccessfully.Inc()

	return &pb.EmailResponse{}, nil
}
