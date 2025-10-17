package routes

import (
	"go-blog/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterMemberRoutes(router *gin.Engine, memberController *controllers.MemberController) {
	memberRoutes := router.Group("/members")
	{
		memberRoutes.POST("", memberController.CreateMember)
		memberRoutes.GET("", memberController.GetAllMembers)
		memberRoutes.GET("/:id", memberController.GetMemberByID)
		memberRoutes.PUT("/:id", memberController.UpdateMember)
		memberRoutes.DELETE("/:id", memberController.DeleteMember)
	}
}
