package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/EdanStasiuk/LiteCode/apps/backend/server/models"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDatabase() *gorm.DB {
	dsn := os.Getenv("NEON_DEV_DB_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	return db
}

func main() {
	// Load environment variables
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connect to Neon DB using GORM
	dsn := os.Getenv("NEON_DATABASE_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto migrate schemas
	if err := db.AutoMigrate(
		&models.User{},
		&models.Problem{},
		&models.TestCase{},
		&models.Submission{},
		&models.Tag{},
		&models.ProblemTag{},
	); err != nil {
		log.Fatal("Failed to migrate schema:", err)
	}

	fmt.Println("Schema migration successful!")

	// For running a server locally
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	fmt.Println("Listening on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("failed to start server: ", err)
	}
}
