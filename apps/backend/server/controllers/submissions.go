package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/EdanStasiuk/LiteCode/apps/backend/server/models"
	"github.com/EdanStasiuk/LiteCode/pkg/cassandra"
	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"github.com/segmentio/kafka-go"
)

var kafkaWriter *kafka.Writer

func InitKafkaProducer(broker string, topic string) {
	kafkaWriter = &kafka.Writer{
		Addr:     kafka.TCP(broker),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
}

// CreateSubmission handles POST /submissions
// Creates a new submission for the authenticated user with problem ID, code, and language.
func CreateSubmission() gin.HandlerFunc {
	return func(c *gin.Context) {
		var body struct {
			ProblemID string `json:"problemId" binding:"required"`
			Code      string `json:"code" binding:"required"`
			Language  string `json:"language" binding:"required"`
		}

		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userID, _ := c.Get("userID")
		subID := gocql.TimeUUID().String()

		// Insert submission into Cassandra with "pending" status
		err := cassandra.InsertSubmission(userID.(string), body.ProblemID, body.Code, "pending", subID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Build submission event
		event := struct {
			SubmissionID string `json:"submission_id"`
			UserID       string `json:"user_id"`
			ProblemID    string `json:"problem_id"`
			Code         string `json:"code"`
			Language     string `json:"language"`
		}{
			SubmissionID: subID,
			UserID:       userID.(string),
			ProblemID:    body.ProblemID,
			Code:         body.Code,
			Language:     body.Language,
		}

		// Serialize to JSON
		eventBytes, err := json.Marshal(event)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to serialize submission event"})
			return
		}

		// Publish submission event to Kafka
		msg := kafka.Message{
			Key:   []byte(subID),
			Value: eventBytes,
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := kafkaWriter.WriteMessages(ctx, msg); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to enqueue submission"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message":      "Submission created successfully",
			"submissionId": subID,
		})
	}
}

// GetSubmissionByID handles GET /submissions/:id
// Returns a submission's details and code for the authenticated user by submission ID.
func GetSubmissionByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		subID := c.Param("id")
		userID, _ := c.Get("userID")

		var code models.SubmissionCode
		if err := cassandra.Session.Query(
			`SELECT submission_id, code FROM submission_code WHERE submission_id = ?`,
			subID,
		).Consistency(gocql.One).Scan(&code.SubmissionID, &code.Code); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Submission not found"})
			return
		}

		var userSub models.Submission
		if err := cassandra.Session.Query(
			`SELECT user_id, submission_id, problem_id, status, runtime, memory, result, created_at
			 FROM submissions_by_user
			 WHERE user_id = ? AND submission_id = ?`,
			userID.(string), subID,
		).Consistency(gocql.One).Scan(
			&userSub.UserID, &userSub.SubmissionID, &userSub.ProblemID, &userSub.Status,
			&userSub.Runtime, &userSub.Memory, &userSub.Result, &userSub.CreatedAt,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch submission metadata"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"submissionId": code.SubmissionID,
			"userId":       userSub.UserID,
			"problemId":    userSub.ProblemID,
			"status":       userSub.Status,
			"runtime":      userSub.Runtime,
			"memory":       userSub.Memory,
			"result":       userSub.Result,
			"code":         code.Code,
			"createdAt":    userSub.CreatedAt,
		})
	}
}

// GetUserSubmissions handles GET /users/:id/submissions
// Returns all submissions made by the authenticated user.
func GetUserSubmissions() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("userID")

		iter := cassandra.Session.Query(
			`SELECT submission_id, problem_id, status, runtime, memory, result, language, created_at
			 FROM submissions_by_user 
			 WHERE user_id = ?`,
			userID.(string),
		).Iter()

		var subs []models.Submission
		var sub models.Submission
		for iter.Scan(&sub.SubmissionID, &sub.ProblemID, &sub.Status, &sub.Runtime, &sub.Memory, &sub.Result, &sub.Language, &sub.CreatedAt) {
			sub.UserID = userID.(string)
			subs = append(subs, sub)
		}
		if err := iter.Close(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch submissions"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"submissions": subs})
	}
}

// GetProblemSubmissions handles GET /problems/:id/submissions
// Returns all submissions for a specific problem by the authenticated user.
func GetProblemSubmissions() gin.HandlerFunc {
	return func(c *gin.Context) {
		problemID := c.Param("id")
		userID := c.MustGet("userID")

		iter := cassandra.Session.Query(
			`SELECT submission_id, status, runtime, memory, result, language, created_at
			 FROM submissions_by_problem_and_user
			 WHERE problem_id = ? AND user_id = ?`,
			problemID, userID.(string),
		).Iter()

		var subs []models.Submission
		var sub models.Submission
		for iter.Scan(&sub.SubmissionID, &sub.Status, &sub.Runtime, &sub.Memory, &sub.Result, &sub.Language, &sub.CreatedAt) {
			sub.ProblemID = problemID
			sub.UserID = userID.(string)
			subs = append(subs, sub)
		}

		if err := iter.Close(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch submissions"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"submissions": subs})
	}
}
