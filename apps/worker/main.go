package main

import (
	"fmt"
	"log"

	"github.com/EdanStasiuk/LiteCode/apps/worker/handlers"
	"github.com/EdanStasiuk/LiteCode/pkg/cassandra"
	"github.com/EdanStasiuk/LiteCode/pkg/kafka"
)

func main() {
	// Cassandra
	if err := cassandra.Init(); err != nil {
		log.Fatal("Failed to connect to Cassandra:", err)
	}
	defer cassandra.Close()
	fmt.Println("Cassandra connected successfully (worker)")

	// Kafka consumer
	reader := kafka.NewConsumer([]string{"kafka:9092"}, "submissions", "worker-group-1")
	defer reader.Close()

	fmt.Println("Worker listening for submissions...")

	handlers.ConsumeSubmissions(reader)
}
