package models

import "time"

// CreateUser represents a user document stored in the database
type CreateUser struct {
	ID        string    `bson:"_id,omitempty" json:"id"`
	Username  string    `bson:"username" json:"username" binding:"required"`
	FirstName string    `bson:"first_name" json:"first_name" binding:"required"`
	LastName  string    `bson:"last_name" json:"last_name" binding:"required"`
	Password  string    `bson:"password" json:"password" binding:"required"`
	Email     string    `bson:"email" json:"email" binding:"required"`
	Phone     string    `bson:"phone" json:"phone" binding:"required"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	PsqlID    int       `bson:"psql_id" json:"psql_id"`
}

// LoginUser represents a user document stored in the database
type LoginUser struct {
	Email    string `bson:"email" json:"email" binding:"required"`
	Password string `bson:"password" json:"password" binding:"required"`
}
