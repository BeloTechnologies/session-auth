package routes

import (
	"session-auth/controllers"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// UserRoutes sets up the user-related routes.
func UserRoutes(r *gin.Engine, db *mongo.Client) {
	userGroup := r.Group("/users")
	{
		userGroup.POST("/create/", controllers.CreateUser(db))
		// Add more routes here (e.g., GET, PUT, DELETE)
	}
}
