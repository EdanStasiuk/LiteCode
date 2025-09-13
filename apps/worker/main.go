package main

import (
	"fmt"
	"log"

	"github.com/EdanStasiuk/LiteCode/apps/worker/handlers"
	"github.com/EdanStasiuk/LiteCode/pkg/cassandra"
	kpkg "github.com/EdanStasiuk/LiteCode/pkg/kafka"
)

func main() {
	// Cassandra
	if err := cassandra.Init(); err != nil {
		log.Fatal("Failed to connect to Cassandra:", err)
	}
	defer cassandra.Close()
	fmt.Println("Cassandra connected successfully (worker)")

	// Kafka consumer for submissions
	subReader := kpkg.NewConsumer(
		[]string{"kafka:9092"},
		"submissions",
		"worker-group-1",
	)
	defer subReader.Close()
	fmt.Println("Worker listening for submissions...")

	// Kafka producer for submission results
	kpkg.InitProducer("kafka:9092", "submission-results")
	defer func() {
		if err := kpkg.CloseProducer(); err != nil {
			log.Printf("Failed to close Kafka producer: %v", err)
		}
	}()
	fmt.Println("Kafka producer ready for submission-results")

	// Start consuming submissions (blocking)
	handlers.ConsumeSubmissions(subReader)
}
