package routes

import (
	"github.com/EdanStasiuk/LiteCode/apps/backend/server/controllers"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterProblemRoutes(r *gin.Engine, db *gorm.DB) {
	r.GET("/problems/:id", controllers.GetProblemByID(db))
	r.GET("/problems", controllers.GetProblems(db))
	r.POST("/problems", controllers.CreateProblem(db))
	r.PUT("/problems/:id", controllers.UpdateProblem(db))
	r.DELETE("/problems/:id", controllers.DeleteProblem(db))
	r.POST("/problems/:id/restore", controllers.RestoreProblem(db))
}
