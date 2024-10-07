package models

import "time"

// User represents a user document stored in the database
type User struct {
	ID        string    `bson:"_id,omitempty" json:"id"`
	Username  string    `bson:"username" json:"username" binding:"required"`
	Password  string    `bson:"password" json:"password" binding:"required"`
	Email     string    `bson:"email" json:"email" binding:"required"`
	Phone     string    `bson:"phone" json:"phone" binding:"required"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
}
