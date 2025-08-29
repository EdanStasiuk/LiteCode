package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/EdanStasiuk/LiteCode/apps/backend/server/models"
	"github.com/EdanStasiuk/LiteCode/apps/backend/server/routes"
	"github.com/EdanStasiuk/LiteCode/pkg/cassandra"
	"github.com/EdanStasiuk/LiteCode/pkg/redis"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load .env
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Postgres setup
	dsn := os.Getenv("NEON_DEV_DB_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Printf("failed to connect to db: %v\n", err)
		return
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

	// Register problem routes
	routes.RegisterProblemRoutes(r, db)

	fmt.Println("Listening on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("failed to start server: ", err)
	}
}
