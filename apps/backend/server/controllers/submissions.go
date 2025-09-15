package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/EdanStasiuk/LiteCode/pkg/cassandra"
	kafkaq "github.com/EdanStasiuk/LiteCode/pkg/kafka"
	"github.com/EdanStasiuk/LiteCode/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
)

// CreateSubmission handles POST /submissions
// Inserts a "pending" submission into Cassandra and enqueues it to Kafka
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

		// Ensure userID is present
		userIDValue, ok := c.Get("userID")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
			return
		}
		userID := userIDValue.(string)
		if userID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "userID cannot be empty"})
			return
		}

		// Generate submission UUID
		subID := gocql.TimeUUID().String()

		fmt.Printf("DEBUG: inserting submission for userID=%q subID=%q problemID=%q\n",
			userID, subID, body.ProblemID)

		// Insert into submissions_by_user
		if err := cassandra.Session.Query(
			`INSERT INTO submissions_by_user
			 (user_id, submission_id, problem_id, status, runtime, memory, result, language, created_at)
			 VALUES (?, ?, ?, ?, ?, ?, ?, ?, toTimestamp(now()))`,
			userID, subID, body.ProblemID, "pending", 0.0, int64(0), "", body.Language,
		).Exec(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to insert into submissions_by_user: " + err.Error()})
			return
		}

		// Insert into submissions_by_problem
		if err := cassandra.Session.Query(
			`INSERT INTO submissions_by_problem
			(problem_id, submission_id, user_id, status, runtime, memory, result, language, created_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, toTimestamp(now()))`,
			body.ProblemID, subID, userID, "pending", 0.0, int64(0), "", body.Language,
		).Exec(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to insert into submissions_by_problem: " + err.Error()})
			return
		}

		// Insert into submissions_by_problem_and_user
		if err := cassandra.Session.Query(
			`INSERT INTO submissions_by_problem_and_user
			(problem_id, user_id, submission_id, status, runtime, memory, result, language, created_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, toTimestamp(now()))`,
			body.ProblemID, userID, subID, "pending", 0.0, int64(0), "", body.Language,
		).Exec(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to insert into submissions_by_problem_and_user: " + err.Error()})
			return
		}

		// Insert into submission_code (to store the actual code)
		if err := cassandra.Session.Query(
			`INSERT INTO submission_code (submission_id, code) VALUES (?, ?)`,
			subID, body.Code,
		).Exec(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to insert into submission_code: " + err.Error()})
			return
		}

		// Build Kafka event
		event := struct {
			SubmissionID string `json:"submission_id"`
			UserID       string `json:"user_id"`
			ProblemID    string `json:"problem_id"`
			Code         string `json:"code"`
			Language     string `json:"language"`
		}{
			SubmissionID: subID,
			UserID:       userID,
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

		// Publish to Kafka
		if err := kafkaq.ProduceMessage(subID, eventBytes); err != nil {
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

		// Ensure userID is in context
		userIDValue, ok := c.Get("userID")
		if !ok || userIDValue == nil || userIDValue == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
			return
		}
		userID := userIDValue.(string)

		// Fetch code
		var code models.SubmissionCode
		if err := cassandra.Session.Query(
			`SELECT submission_id, code FROM submission_code WHERE submission_id = ?`,
			subID,
		).Consistency(gocql.One).Scan(&code.SubmissionID, &code.Code); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Submission not found"})
			return
		}

		// Fetch submission metadata for this user
		var userSub models.Submission
		if err := cassandra.Session.Query(
			`SELECT user_id, submission_id, problem_id, status, runtime, memory, result, language, created_at
			 FROM submissions_by_user
			 WHERE user_id = ? AND submission_id = ?`,
			userID, subID,
		).Consistency(gocql.One).Scan(
			&userSub.UserID, &userSub.SubmissionID, &userSub.ProblemID, &userSub.Status,
			&userSub.Runtime, &userSub.Memory, &userSub.Result, &userSub.Language, &userSub.CreatedAt,
		); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Submission not found or not owned by user"})
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
			"language":     userSub.Language,
			"code":         code.Code,
			"createdAt":    userSub.CreatedAt,
		})
	}
}

// GetUserSubmissions handles GET /users/submissions
// Returns all submissions made by the authenticated user.
func GetUserSubmissions() gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDValue, ok := c.Get("userID")
		if !ok || userIDValue == nil || userIDValue == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
			return
		}
		userID := userIDValue.(string)

		iter := cassandra.Session.Query(
			`SELECT submission_id, problem_id, status, runtime, memory, result, language, created_at
			 FROM submissions_by_user 
			 WHERE user_id = ?`,
			userID,
		).Iter()

		subs := make([]models.Submission, 0)
		var sub models.Submission
		for iter.Scan(&sub.SubmissionID, &sub.ProblemID, &sub.Status, &sub.Runtime, &sub.Memory, &sub.Result, &sub.Language, &sub.CreatedAt) {
			sub.UserID = userID
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

		userIDValue, ok := c.Get("userID")
		if !ok || userIDValue == nil || userIDValue == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
			return
		}
		userID := userIDValue.(string)

		iter := cassandra.Session.Query(
			`SELECT submission_id, status, runtime, memory, result, language, created_at
			 FROM submissions_by_problem_and_user
			 WHERE problem_id = ? AND user_id = ?`,
			problemID, userID,
		).Iter()

		subs := make([]models.Submission, 0)
		var sub models.Submission
		for iter.Scan(&sub.SubmissionID, &sub.Status, &sub.Runtime, &sub.Memory, &sub.Result, &sub.Language, &sub.CreatedAt) {
			sub.ProblemID = problemID
			sub.UserID = userID
			subs = append(subs, sub)
		}

		if err := iter.Close(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch submissions"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"submissions": subs})
	}
}
