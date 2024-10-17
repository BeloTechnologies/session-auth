package models

import "github.com/BeloTechnologies/session-core/core_models/user_models"

type AuthResponse struct {
	Token    string           `json:"token"`
	UserData user_models.User `json:"user_data"`
}
