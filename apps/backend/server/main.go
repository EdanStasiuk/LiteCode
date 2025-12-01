package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/EdanStasiuk/LiteCode/apps/backend/server/routes"
	"github.com/EdanStasiuk/LiteCode/pkg/cassandra"
	"github.com/EdanStasiuk/LiteCode/pkg/models"
	rmq "github.com/EdanStasiuk/LiteCode/pkg/rabbitmq"
	"github.com/EdanStasiuk/LiteCode/pkg/redis"
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

	// RabbitMQ
	url := os.Getenv("RABBITMQ_URL")
	if url == "" {
		log.Fatal("RABBIMQ_URL not set")
	}

	// Producer
	if err := rmq.InitProducer(url, "submissions"); err != nil {
		log.Fatalf("Failed to initialize RabbitMQ producer: %v", err)
	}
	defer rmq.CloseProducer()
	log.Println("Backend RabbitMQ producer initialized")

	// Consumer
	consumer, err := rmq.NewConsumer(url, "submission-results")
	if err != nil {
		log.Fatal(err)
	}
	defer consumer.Close()
	log.Println("Backend RabbitMQ consumer initialized")

	go func() {
		for msg := range consumer.Msgs {
			var res models.SubmissionResult
			if err := json.Unmarshal(msg.Body, &res); err != nil {
				log.Printf("Invalid result message: %v", err)
				msg.Nack(false, false)
				continue
			}

			if err := cassandra.UpdateSubmissionResult(res); err != nil {
				log.Printf("Failed to update Cassandra for submission %s: %v", res.SubmissionID, err)
			}

			msg.Ack(false)
		}
	}()

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
