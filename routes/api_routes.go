package routes

import (
	"go-blog/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, planController *controllers.PlanController, membershipController *controllers.MembershipController, paymentController *controllers.PaymentController) {
	api := r.Group("/api")

	// Plan routes
	api.POST("/plans", planController.CreatePlan)
	api.GET("/plans", planController.GetPlans)

	// Membership routes
	api.POST("/memberships", membershipController.CreateMembership)
	api.GET("/memberships/:memberID", membershipController.GetByMember)

	// Payment routes
	api.POST("/payments", paymentController.RecordPayment)
	api.GET("/payments/:memberID", paymentController.GetPayments)
	api.GET("/payments", paymentController.GetAllPayments) // <-- all payments

}
