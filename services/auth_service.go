package services

import (
	"github.com/BeloTechnologies/session-core/core_models"
	"github.com/BeloTechnologies/session-core/core_models/auth_models"
)

// ValidateToken validates the token provided by the user
func ValidateToken(token auth_models.Token) (bool, *core_models.SessionError) {
	isValid := ValidateJwt(token.Token)
	if !isValid {
		return false, &core_models.SessionError{
			Message:     "Invalid token",
			Errors:      "Invalid token",
			Status:      401,
			Description: "The token provided is invalid. Please check the token and try again.",
		}
	}

	return true, nil
}
