package test

import (
	"go-email/internal/validator"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEmailValidation_InvalidEmail(t *testing.T) {
	emails := []struct {
		Email string
	}{
		{Email: "john.doegmail.com"},
		{Email: "test-:@email@yahoo"},
	}

	for _, email := range emails {
		t.Run(email.Email, func(t *testing.T) {
			err := validator.ValidateEmail(email.Email)

			require.False(t, err)
		})
	}
}

func TestEmailValidation_ValidEmail(t *testing.T) {
	emails := []struct {
		Email string
	}{
		{Email: "john.doe@gmail.com"},
		{Email: "test-email@yahoo.com"},
	}

	for _, email := range emails {
		t.Run(email.Email, func(t *testing.T) {
			err := validator.ValidateEmail(email.Email)

			assert.True(t, err)
		})
	}
}
