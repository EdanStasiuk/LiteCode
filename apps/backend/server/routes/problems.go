package routes

import (
	"github.com/EdanStasiuk/LiteCode/apps/backend/server/controllers"
	"github.com/EdanStasiuk/LiteCode/apps/backend/server/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterProblemRoutes(r *gin.Engine, db *gorm.DB) {
	r.GET("/problems", controllers.GetProblems(db))
	r.GET("/problems/:id", controllers.GetProblemByID(db))

	// Admin-only routes
	r.POST("/problems", middleware.AuthMiddleware(), middleware.AdminMiddleware(db), controllers.CreateProblem(db))
	r.PUT("/problems/:id", middleware.AuthMiddleware(), middleware.AdminMiddleware(db), controllers.UpdateProblem(db))
	r.DELETE("/problems/:id", middleware.AuthMiddleware(), middleware.AdminMiddleware(db), controllers.DeleteProblem(db))
	r.PATCH("/problems/:id/restore", middleware.AuthMiddleware(), middleware.AdminMiddleware(db), controllers.RestoreProblem(db))
}
