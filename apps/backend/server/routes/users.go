package routes

import (
	"github.com/EdanStasiuk/LiteCode/apps/backend/server/controllers"
	"github.com/EdanStasiuk/LiteCode/apps/backend/server/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterUserRoutes(r *gin.Engine, db *gorm.DB) {
	users := r.Group("/users")
	users.Use(middleware.AuthMiddleware())

	users.GET("/:id", controllers.GetUserProfile(db))
	users.PUT("/:id", controllers.UpdateUserProfile(db))
	users.GET("/:id/stats", controllers.GetUserStats(db))
}
