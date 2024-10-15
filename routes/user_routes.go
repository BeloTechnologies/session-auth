package routes

import (
	"session-auth/controllers"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// UserRoutes sets up the user-related routes.
func UserRoutes(r *gin.Engine, db *mongo.Client) {
	userGroup := r.Group("/auth")
	{
		userGroup.POST("/create/", controllers.CreateUser(db))
		userGroup.POST("/login/", controllers.LoginUser(db))
	}
}
