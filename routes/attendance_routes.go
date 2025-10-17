package routes

import (
	"go-blog/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterAttendanceRoutes(router *gin.Engine, c *controllers.AttendanceController) {
	group := router.Group("/attendance")
	{
		group.POST("/checkin", c.CheckIn)
		group.GET("/member/:member_id", c.GetMemberAttendance)
		group.GET("/all", c.GetAllAttendance)
	}
}
