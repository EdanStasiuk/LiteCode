package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/EdanStasiuk/LiteCode/apps/worker/pkg/sandbox"
	"github.com/EdanStasiuk/LiteCode/pkg/cassandra"
	"github.com/EdanStasiuk/LiteCode/pkg/models"
	rmq "github.com/EdanStasiuk/LiteCode/pkg/rabbitmq"
)

type SubmissionEvent struct {
	SubmissionID string `json:"submission_id"`
	UserID       string `json:"user_id"`
	ProblemID    string `json:"problem_id"`
	Code         string `json:"code"`
	Language     string `json:"language"`
}

func main() {
	// Cassandra
	if err := cassandra.Init(); err != nil {
		log.Fatal("Failed to connect to Cassandra:", err)
	}
	defer cassandra.Close()
	fmt.Println("Cassandra connected successfully (worker)")

	// RabbitMQ
	url := os.Getenv("RABBITMQ_URL")
	if url == "" {
		log.Fatal("RABBITMQ_URL not set")
	}

	// Producer: submission-results
	if err := rmq.InitProducer(url, "submission-results"); err != nil {
		log.Fatal(err)
	}
	defer rmq.CloseProducer()
	log.Println("Worker RabbitMQ producer initialized for submission-results")

	// Consumer: submissions
	consumer, err := rmq.NewConsumer(url, "submissions")
	if err != nil {
		log.Fatal(err)
	}
	defer consumer.Close()
	fmt.Println("Worker listening for submissions...")

	// Consume messages
	for msg := range consumer.Msgs {
		var subEvent SubmissionEvent
		if err := json.Unmarshal(msg.Body, &subEvent); err != nil {
			log.Printf("Invalid submission event: %v", err)
			if err := msg.Nack(false, false); err != nil {
				log.Printf("Failed to NACK message: %v", err)
			}
			continue
		}

		// Optionally fetch from Cassandra if needed
		codeObj, err := cassandra.GetSubmissionCode(subEvent.SubmissionID)
		if err != nil || codeObj == nil {
			log.Printf("Submission code not found: %v", subEvent.SubmissionID)
			if err := msg.Nack(false, false); err != nil {
				log.Printf("Failed to NACK message: %v", err)
			}
			continue
		}

		status, result, runtime, memory := sandbox.RunCode(codeObj.Code, codeObj.Language)

		_ = cassandra.UpdateSubmissionResult(models.SubmissionResult{
			SubmissionID: subEvent.SubmissionID,
			UserID:       subEvent.UserID,
			ProblemID:    subEvent.ProblemID,
			Status:       status,
			Result:       result,
			Runtime:      runtime,
			Memory:       memory,
		})

		if err := msg.Ack(false); err != nil {
			log.Printf("Failed to ACK message: %v", err)
		}
	}

}
