package services

import (
	"github.com/BeloTechnologies/session-core/core_models/auth_models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidateToken(t *testing.T) {
	// Create a new token
	testToken, err := GenerateJwt("testID")
	if err != nil {
		t.Errorf("Error generating token: %v", err)
	}

	// Validate the token
	result, err := ValidateToken(auth_models.Token{
		Token: testToken,
	})

	assert.Nil(t, err)
	assert.Truef(t, result, "Expected token to be valid")
}

func TestInvalidToken(t *testing.T) {
	// Validate an invalid token
	result, err := ValidateToken(auth_models.Token{
		Token: "invalidToken",
	})

	assert.NotNil(t, err)
	assert.Falsef(t, result, "Expected token to be invalid")
}
