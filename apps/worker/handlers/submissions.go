package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/EdanStasiuk/LiteCode/apps/backend/server/models"
	"github.com/EdanStasiuk/LiteCode/pkg/cassandra"
	"github.com/segmentio/kafka-go"
)

type SubmissionMessage struct {
	SubmissionID string `json:"submissionId"`
	UserID       string `json:"userId"`
	ProblemID    string `json:"problemId"`
	Code         string `json:"code"`
	Language     string `json:"language"`
}

func ConsumeSubmissions(reader *kafka.Reader) {
	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Error reading message: %v", err)
			continue
		}

		var msg SubmissionMessage
		if err := json.Unmarshal(m.Value, &msg); err != nil {
			log.Printf("Invalid message format: %v", err)
			continue
		}

		fmt.Printf("Processing submission %s for user %s\n", msg.SubmissionID, msg.UserID)

		// TODO: Run code execution inside Docker sandbox
		result := "Accepted"
		runtime := 0.123
		memory := int64(2048)

		res := models.SubmissionResult{
			SubmissionID: msg.SubmissionID,
			UserID:       msg.UserID,
			ProblemID:    msg.ProblemID,
			Status:       result,
			Runtime:      runtime,
			Memory:       memory,
		}

		err = cassandra.UpdateSubmissionResult(res)

		if err != nil {
			log.Printf("Failed to update submission result: %v", err)
			continue
		}

		fmt.Printf("Submission %s processed successfully\n", msg.SubmissionID)
	}
}
