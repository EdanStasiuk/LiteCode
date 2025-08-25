package cassandra

import (
	"time"

	"github.com/EdanStasiuk/LiteCode/apps/backend/server/models"
	"github.com/gocql/gocql"
)

// InsertUserSubmission writes the denormalized user submission table
func InsertUserSubmission(sub models.UserSubmission) error {
	return Session.Query(
		`INSERT INTO submissions_by_user (user_id, submission_id, problem_id, status, created_at)
         VALUES (?, ?, ?, ?, ?)`,
		sub.UserID, sub.SubmissionID, sub.ProblemID, sub.Status, sub.CreatedAt,
	).Exec()
}

// InsertProblemSubmission writes the denormalized problem submission table
func InsertProblemSubmission(sub models.ProblemSubmission) error {
	return Session.Query(
		`INSERT INTO submissions_by_problem (problem_id, submission_id, user_id, status, created_at)
         VALUES (?, ?, ?, ?, ?)`,
		sub.ProblemID, sub.SubmissionID, sub.UserID, sub.Status, sub.CreatedAt,
	).Exec()
}

// InsertSubmissionCode writes the full code table
func InsertSubmissionCode(sub models.SubmissionCode) error {
	return Session.Query(
		`INSERT INTO submission_code (submission_id, code)
         VALUES (?, ?)`,
		sub.SubmissionID, sub.Code,
	).Exec()
}

// InsertSubmission inserts a new submission into all three Cassandra tables
func InsertSubmission(userID, problemID, code, status string) error {
	subID := gocql.TimeUUID().String()
	now := time.Now()

	if err := InsertUserSubmission(models.UserSubmission{
		UserID:       userID,
		SubmissionID: subID,
		ProblemID:    problemID,
		Status:       status,
		CreatedAt:    now,
	}); err != nil {
		return err
	}

	if err := InsertProblemSubmission(models.ProblemSubmission{
		ProblemID:    problemID,
		SubmissionID: subID,
		UserID:       userID,
		Status:       status,
		CreatedAt:    now,
	}); err != nil {
		return err
	}

	if err := InsertSubmissionCode(models.SubmissionCode{
		SubmissionID: subID,
		Code:         code,
	}); err != nil {
		return err
	}

	return nil
}
