package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	ID           uint   `gorm:"primaryKey"`
	Username     string `gorm:"unique;not null"`
	Email        string `gorm:"unique;not null"`
	PasswordHash string `gorm:"not null"`
	CreatedAt    time.Time
	Submissions  []Submission
}

type Problem struct {
	ID          uint   `gorm:"primaryKey"`
	Title       string `gorm:"not null"`
	Slug        string `gorm:"unique;not null"`
	Difficulty  string `gorm:"not null"`
	Description string `gorm:"type:text"`
	Constraints string `gorm:"type:text"`
	CreatedAt   time.Time
	TestCases   []TestCase
	Tags        []Tag `gorm:"many2many:problem_tags;"`
	Submissions []Submission
}

type TestCase struct {
	ID        uint   `gorm:"primaryKey"`
	ProblemID uint   `gorm:"not null"`
	Input     string `gorm:"type:text"`
	Expected  string `gorm:"type:text"`
	IsSample  bool
}

type Submission struct {
	ID        uint   `gorm:"primaryKey"`
	UserID    uint   `gorm:"not null"`
	ProblemID uint   `gorm:"not null"`
	Language  string `gorm:"not null"`
	Code      string `gorm:"type:text"`
	Status    string
	RuntimeMs *int
	MemoryKb  *int
	CreatedAt time.Time
}

type Tag struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"unique;not null"`
}

type ProblemTag struct {
	ProblemID uint
	TagID     uint
}

func main() {
	dsn := "host=localhost user=litecode_user password=litecode_pass dbname=litecode port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to database: ", err)
	}

	// Migrate schema
	err = db.AutoMigrate(
		&User{},
		&Problem{},
		&TestCase{},
		&Submission{},
		&Tag{},
		&ProblemTag{},
	)
	if err != nil {
		log.Fatal("failed to migrate schema: ", err)
	}

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	fmt.Println("Listening on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("failed to start server: ", err)
	}
}
