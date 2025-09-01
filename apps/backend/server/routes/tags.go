package routes

import (
	"github.com/EdanStasiuk/LiteCode/apps/backend/server/controllers"
	"github.com/EdanStasiuk/LiteCode/apps/backend/server/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterTagRoutes(r *gin.Engine, db *gorm.DB) {
	tags := r.Group("/tags")
	{
		// Public route
		tags.GET("/", controllers.ListTags(db))

		// Admin-only routes
		tags.POST("/", middleware.AuthMiddleware(), middleware.AdminMiddleware(db), controllers.CreateTag(db))
		tags.DELETE("/:id", middleware.AuthMiddleware(), middleware.AdminMiddleware(db), controllers.DeleteTag(db))
	}
}
