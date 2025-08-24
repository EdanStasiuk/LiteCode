package models

import "time"

// Denormalized Cassandra model for storing code
type SubmissionCode struct {
	SubmissionID string
	Code         string
}

type UserSubmission struct {
	UserID       string
	SubmissionID string
	ProblemID    string
	Status       string
	CreatedAt    time.Time
}

type ProblemSubmission struct {
	ProblemID    string
	SubmissionID string
	UserID       string
	Status       string
	CreatedAt    time.Time
}
