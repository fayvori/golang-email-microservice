package test

import (
	"fmt"
	repo "go-email/internal/database"
	"go-email/internal/models"
	db "go-email/pkg/database"
	"log"
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

var dbConn *gorm.DB

func init() {
	db, err := db.NewDatabase(cfg)

	if err != nil {
		log.Fatal("cannot connect to the database for tests")
	}

	dbConn = db
}

func TestEmailRepository_TestPostgresCreateEmail(t *testing.T) {
	repo := repo.NewRepository(dbConn)

	email := &models.Email{
		From:        cfg.Smtp.User,
		To:          []string{"alexemailtestingtwo@yahoo.com"},
		ContentType: "text/plain",
		Subject:     "Testing",
		Body:        "Test email",
	}

	err := repo.CreateEmail(email)

	if err != nil {
		fmt.Println(err.Error())
	}

	require.NoError(t, err)
}
