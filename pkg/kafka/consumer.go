package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/EdanStasiuk/LiteCode/apps/backend/server/models"
	"github.com/EdanStasiuk/LiteCode/pkg/cassandra"
	"github.com/segmentio/kafka-go"
)

func NewConsumer(brokers []string, topic, groupID string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		GroupID: groupID,
		Topic:   topic,
	})
}

// ConsumeSubmissionResults listens for submission results and updates Cassandra
func ConsumeSubmissionResults(brokers []string, topic, groupID string) {
	reader := NewConsumer(brokers, topic, groupID)
	defer reader.Close()

	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Error reading result message: %v", err)
			continue
		}

		var res models.SubmissionResult
		if err := json.Unmarshal(m.Value, &res); err != nil {
			log.Printf("Invalid result message format: %v", err)
			continue
		}

		if err := cassandra.UpdateSubmissionResult(res); err != nil {
			log.Printf("Failed to update Cassandra for submission %s: %v", res.SubmissionID, err)
			continue
		}

		log.Printf("Submission %s updated in Cassandra (status: %s, result: %s)", res.SubmissionID, res.Status, res.Result)
	}
}
