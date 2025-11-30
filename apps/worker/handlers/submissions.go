package handlers

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/EdanStasiuk/LiteCode/apps/worker/pkg/sandbox"
	"github.com/EdanStasiuk/LiteCode/pkg/models"
	rmq "github.com/EdanStasiuk/LiteCode/pkg/rabbitmq"
)

type SubmissionMessage struct {
	SubmissionID string `json:"submission_id"`
	UserID       string `json:"user_id"`
	ProblemID    string `json:"problem_id"`
	Code         string `json:"code"`
	Language     string `json:"language"`
}

// ConsumeSubmissionsRabbitMQ consumes messages from RabbitMQ submissions queue
func ConsumeSubmissionsRabbitMQ(consumer *rmq.Consumer) {
	for msg := range consumer.Msgs {
		var submission SubmissionMessage
		if err := json.Unmarshal(msg.Body, &submission); err != nil {
			log.Printf("Invalid message format: %v", err)
			msg.Nack(false, false) // reject the message, don't requeue
			continue
		}

		fmt.Printf("Processing submission %s for user %s\n", submission.SubmissionID, submission.UserID)

		// Run code in sandbox
		status, result, runtime, memory := sandbox.RunCode(submission.Code, submission.Language)

		res := models.SubmissionResult{
			SubmissionID: submission.SubmissionID,
			UserID:       submission.UserID,
			ProblemID:    submission.ProblemID,
			Status:       status,
			Result:       result,
			Runtime:      runtime,
			Memory:       memory,
		}

		// Publish to "submission-results" queue
		resultMsg, _ := json.Marshal(res)
		if err := rmq.ProduceMessage(resultMsg); err != nil {
			log.Printf("Failed to produce result message: %v", err)
			msg.Nack(false, true) // requeue the original message for retry
			continue
		}

		fmt.Printf("Submission %s processed successfully\n", submission.SubmissionID)
		msg.Ack(false) // acknowledge the message
	}
}
