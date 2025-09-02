package models

import (
	"time"
)

// Denormalized Cassandra model for storing code
type SubmissionCode struct {
	SubmissionID string
	Code         string
	Language     string
}

type Submission struct {
	UserID       string
	SubmissionID string
	ProblemID    string
	Status       string
	Runtime      float64 // in seconds
	Memory       int64   // in bytes
	Language     string
	Result       string // e.g., "Accepted", "Wrong Answer"
	CreatedAt    time.Time
}
