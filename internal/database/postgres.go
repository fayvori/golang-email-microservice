package database

import (
	"go-email/internal/models"
	"gorm.io/gorm"
)

type MailRepository interface {
	CreateEmail(email *models.Email) error
}

type Resository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Resository {
	return &Resository{db: db}
}

func (r *Resository) CreateEmail(email *models.Email) error {
	emailModel := &models.EmailModel{
		From:        email.From,
		To:          email.GetToAsString(),
		ContentType: email.ContentType,
		Subject:     email.Subject,
		Body:        email.Body,
	}

	if result := r.db.Create(emailModel); result.Error != nil {
		return result.Error
	}

	return nil
}
