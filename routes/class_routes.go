package routes

import (
	"go-blog/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterClassRoutes(r *gin.Engine, ctrl *controllers.ClassController) {
	group := r.Group("/class")
	{
		group.POST("", ctrl.CreateClass)     // Remove trailing slash - should be "" not "/"
		group.GET("/get", ctrl.ListClasses)      // Remove trailing slash - should be "" not "/"
		group.GET("/:id", ctrl.GetClass)     // Get class by ID
	}
}