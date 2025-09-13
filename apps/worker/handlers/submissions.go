package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/EdanStasiuk/LiteCode/apps/backend/server/models"
	"github.com/EdanStasiuk/LiteCode/apps/worker/pkg/sandbox"
	kpkg "github.com/EdanStasiuk/LiteCode/pkg/kafka"
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

		status, result, runtime, memory := sandbox.RunCode(msg.Code, msg.Language)

		res := models.SubmissionResult{
			SubmissionID: msg.SubmissionID,
			UserID:       msg.UserID,
			ProblemID:    msg.ProblemID,
			Status:       status,
			Result:       result,
			Runtime:      runtime,
			Memory:       memory,
		}

		// To "submission-results"
		resultMsg, _ := json.Marshal(res)
		if err := kpkg.ProduceMessage(res.SubmissionID, resultMsg); err != nil {
			log.Printf("Failed to produce Kafka message: %v", err)
			continue
		}

		fmt.Printf("Submission %s processed successfully\n", msg.SubmissionID)
	}
}
