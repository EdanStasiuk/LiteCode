package routes

import (
	"github.com/EdanStasiuk/LiteCode/apps/backend/server/controllers"
	"github.com/EdanStasiuk/LiteCode/apps/backend/server/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterAuthRoutes(r *gin.Engine, db *gorm.DB) {
	auth := r.Group("/auth")
	{
		auth.POST("/register", controllers.Register(db))
		auth.POST("/login", controllers.Login(db))
		auth.POST("/logout", middleware.AuthMiddleware(), controllers.Logout())
		auth.GET("/me", middleware.AuthMiddleware(), controllers.Me(db))
	}
}
