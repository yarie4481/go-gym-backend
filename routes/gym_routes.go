package routes

import (
	"go-blog/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterGymRoutes(r *gin.Engine, ctrl *controllers.GymController) {
	group := r.Group("/gymx")
	{
		group.POST("", ctrl.CreateGym)      // Remove trailing slash - should be "" not "/"
		group.GET("", ctrl.ListGyms)        // Remove trailing slash - should be "" not "/"
		group.GET("/:id", ctrl.GetGym)      // Get gym by ID
		group.PUT("/:id", ctrl.UpdateGym)   // Update gym
	}
}