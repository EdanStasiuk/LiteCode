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
	r.GET("/problems/frontend/:frontendId", controllers.GetProblemByFrontendID(db))

	// Admin-only routes
	r.POST("/problems", middleware.AuthMiddleware(), middleware.AdminMiddleware(db), controllers.CreateProblem(db))
	r.PUT("/problems/:id", middleware.AuthMiddleware(), middleware.AdminMiddleware(db), controllers.UpdateProblem(db))
	r.DELETE("/problems/:id", middleware.AuthMiddleware(), middleware.AdminMiddleware(db), controllers.DeleteProblem(db))
	r.PATCH("/problems/:id/restore", middleware.AuthMiddleware(), middleware.AdminMiddleware(db), controllers.RestoreProblem(db))
	r.PUT("/problems/frontend/:frontendId", middleware.AuthMiddleware(), middleware.AdminMiddleware(db), controllers.UpdateProblemByFrontendID(db))
	r.DELETE("/problems/frontend/:frontendId", middleware.AuthMiddleware(), middleware.AdminMiddleware(db), controllers.DeleteProblemByFrontendID(db))
	r.PATCH("/problems/frontend/:frontendId/restore", middleware.AuthMiddleware(), middleware.AdminMiddleware(db), controllers.RestoreProblemByFrontendID(db))
}
