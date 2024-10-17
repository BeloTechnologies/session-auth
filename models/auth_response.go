package models

type AuthResponse struct {
	Token          string `json:"token"`
	Username       string `json:"username"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Email          string `json:"email"`
	Phone          string `json:"phone"`
	CreatedAt      string `json:"created_at"`
	FollowersCount int    `json:"followers_count"`
	FollowingCount int    `json:"following_count"`
}
