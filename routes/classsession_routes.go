package routes

import (
	"go-blog/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterClassSessionRoutes(r *gin.Engine, ctrl *controllers.ClassSessionController) {
	group := r.Group("/classsession")
	{
		group.POST("/", ctrl.CreateSession)      // Create a session
		group.GET("/", ctrl.ListSessions)        // List all sessions
		group.GET("/:id", ctrl.GetSession)       // Get session by ID
	}
}
