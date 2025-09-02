package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/EdanStasiuk/LiteCode/apps/backend/server/models"
	"github.com/EdanStasiuk/LiteCode/apps/backend/server/pkg/cassandra"
	"github.com/EdanStasiuk/LiteCode/apps/backend/server/pkg/redis"
	"github.com/EdanStasiuk/LiteCode/apps/backend/server/routes"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Postgres setup
	dsn := os.Getenv("NEON_DEV_DB_URL")
	if dsn == "" {
		log.Fatal("NEON_DEV_DB_URL not set")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	fmt.Println("Postgres connected successfully")

	if err := db.AutoMigrate(
		&models.User{},
		&models.Problem{},
		&models.TestCase{},
		&models.Language{},
		&models.Tag{},
		&models.Hint{},
		&models.SimilarProblem{},
	); err != nil {
		log.Fatal("Failed to migrate schema:", err)
	}

	fmt.Println("Schema migration successful")

	// Cassandra setup
	if err := cassandra.Init(); err != nil {
		log.Fatal("Failed to connect to Cassandra:", err)
	}
	defer cassandra.Close()
	fmt.Println("Cassandra connected successfully")

	// Redis
	redis.InitRedis()
	defer redis.Rdb.Close()
	fmt.Println("Redis connected succesfully")

	// Gin routes
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	// Register routes
	routes.RegisterProblemRoutes(r, db)
	routes.RegisterAuthRoutes(r, db)
	routes.RegisterTagRoutes(r, db)
	routes.RegisterSubmissionRoutes(r)
	routes.RegisterUserRoutes(r, db)

	fmt.Println("Listening on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("failed to start server: ", err)
	}
}
