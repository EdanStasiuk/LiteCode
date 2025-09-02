package routes

import (
	"github.com/EdanStasiuk/LiteCode/apps/backend/server/controllers"
	"github.com/EdanStasiuk/LiteCode/apps/backend/server/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterSubmissionRoutes(r *gin.Engine) {
	sub := r.Group("/submissions")
	sub.Use(middleware.AuthMiddleware())

	sub.POST("/", controllers.CreateSubmission())
	sub.GET("/:id", controllers.GetSubmissionByID())

	// User-specific submissions
	r.GET("/users/:id/submissions", middleware.AuthMiddleware(), controllers.GetUserSubmissions())

	// Problem-specific submissions (only own submissions)
	r.GET("/problems/:id/submissions", middleware.AuthMiddleware(), controllers.GetProblemSubmissions())
}
