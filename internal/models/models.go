package models

import (
	"strings"

	"gorm.io/gorm"
)

// email model for grpc and rabbitmq
type Email struct {
	From        string   `json:"from"`
	To          []string `json:"to"`
	ContentType string   `json:"contentType"`
	Subject     string   `json:"subject"`
	Body        string   `json:"body"`
}

func (e *Email) GetToAsString() string {
	return strings.Join(e.To, ",")
}

// email model for database
type EmailModel struct {
	gorm.Model
	From        string `json:"from"`
	To          string `json:"to"`
	ContentType string `json:"contentType"`
	Subject     string `json:"subject"`
	Body        string `json:"body"`
}

func (EmailModel) TableName() string {
	return "emails"
}
