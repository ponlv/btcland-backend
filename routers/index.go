package routers

import (
	"api/business/auth"
	"api/business/dashboard"
	"api/business/departments"
	"api/business/healthcheck"
	"api/business/profile"
	"api/business/teams"
	workconfirmations "api/business/work-confirmations"

	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {
	r.GET("healthcheck", healthcheck.Healthcheck())

	// Auth routes
	authRouter := r.Group("auth")
	auth.AddRouter(authRouter)

	// Work confirmations routes
	workConfirmationRouter := r.Group("work-confirmations")
	workconfirmations.Router(workConfirmationRouter)

	// Teams routes (for managers)
	teamsRouter := r.Group("teams")
	teams.Router(teamsRouter)

	// Departments routes (for leaders)
	departmentsRouter := r.Group("teams")
	departments.Router(departmentsRouter)

	// Dashboard routes (for leaders)
	dashboardRouter := r.Group("dashboard")
	dashboard.Router(dashboardRouter)

	// Profile routes (for all authenticated users)
	profileRouter := r.Group("profile")
	profile.Router(profileRouter)
}
