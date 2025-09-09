package cassandra

import (
	"time"

	"github.com/EdanStasiuk/LiteCode/apps/backend/server/models"
)

// InsertUserSubmission writes the denormalized user submission table
func InsertUserSubmission(sub models.Submission) error {
	return Session.Query(
		`INSERT INTO submissions_by_user
		(user_id, submission_id, problem_id, status, runtime, memory, result, language, created_at)
         VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		sub.UserID, sub.SubmissionID, sub.ProblemID, sub.Status, sub.Runtime, sub.Memory, sub.Result, sub.Language, sub.CreatedAt,
	).Exec()
}

// InsertProblemSubmission writes the denormalized problem submission table
func InsertProblemSubmission(sub models.Submission) error {
	return Session.Query(
		`INSERT INTO submissions_by_problem
		(problem_id, submission_id, user_id, status, runtime, memory, result, language, created_at)
         VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		sub.ProblemID, sub.SubmissionID, sub.UserID, sub.Status, sub.Runtime, sub.Memory, sub.Result, sub.Language, sub.CreatedAt,
	).Exec()
}

// InsertProblemSubmissionForUser writes the denormalized problem+user table
func InsertProblemSubmissionForUser(sub models.Submission) error {
	return Session.Query(
		`INSERT INTO submissions_by_problem_and_user
		(problem_id, user_id, submission_id, status, runtime, memory, result, language, created_at)
         VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		sub.ProblemID, sub.UserID, sub.SubmissionID, sub.Status, sub.Runtime, sub.Memory, sub.Result, sub.Language, sub.CreatedAt,
	).Exec()
}

// InsertSubmissionCode writes the full code table
func InsertSubmissionCode(sub models.SubmissionCode) error {
	return Session.Query(
		`INSERT INTO submission_code (submission_id, code, language)
         VALUES (?, ?, ?)`,
		sub.SubmissionID, sub.Code, sub.Language,
	).Exec()
}

// InsertSubmission inserts a new submission into all relevant Cassandra tables
func InsertSubmission(userID, problemID, code, status, subID string) error {
	now := time.Now()

	if err := InsertUserSubmission(models.Submission{
		UserID:       userID,
		SubmissionID: subID,
		ProblemID:    problemID,
		Status:       status,
		CreatedAt:    now,
	}); err != nil {
		return err
	}

	if err := InsertProblemSubmission(models.Submission{
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

	if err := InsertProblemSubmissionForUser(models.Submission{
		ProblemID:    problemID,
		SubmissionID: subID,
		UserID:       userID,
		Status:       status,
		CreatedAt:    now,
	}); err != nil {
		return err
	}

	return nil
}
