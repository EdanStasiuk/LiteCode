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

func main() {
	// Cassandra
	if err := cassandra.Init(); err != nil {
		log.Fatal("Failed to connect to Cassandra:", err)
	}
	defer cassandra.Close()
	fmt.Println("Cassandra connected successfully (worker)")

	// RabbitMQ URL
	url := os.Getenv("RABBITMQ_URL")
	if url == "" {
		log.Fatal("RABBITMQ_URL not set")
	}

	// Producer: submission-results
	if err := rmq.InitProducer(url, "submission-results"); err != nil {
		log.Fatal(err)
	}
	defer rmq.CloseProducer()

	// Consumer: submissions
	consumer, err := rmq.NewConsumer(url, "submissions")
	if err != nil {
		log.Fatal(err)
	}
	defer consumer.Close()
	fmt.Println("Worker listening for submissions...")

	// Consume messages
	for msg := range consumer.Msgs {
		var submission models.Submission
		if err := json.Unmarshal(msg.Body, &submission); err != nil {
			log.Printf("Invalid submission: %v", err)
			msg.Nack(false, false)
			continue
		}

		codeObj, err := cassandra.GetSubmissionCode(submission.SubmissionID)
		if err != nil {
			log.Printf("Failed to fetch code: %v", err)
			msg.Nack(false, false)
			continue
		}
		if codeObj == nil {
			log.Printf("Submission code not found: %v", submission.SubmissionID)
			msg.Nack(false, false)
			continue
		}

		// Run the code
		status, result, runtime, memory := sandbox.RunCode(codeObj.Code, codeObj.Language)

		_ = cassandra.UpdateSubmissionResult(models.SubmissionResult{
			SubmissionID: submission.SubmissionID,
			UserID:       submission.UserID,
			ProblemID:    submission.ProblemID,
			Status:       status,
			Result:       result,
			Runtime:      runtime,
			Memory:       memory,
		})
	}

}
