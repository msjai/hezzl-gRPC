package broker

import (
	"context"
	"github.com/segmentio/kafka-go"
	"log"
	"time"
)

const (
	topic          = "HezzlLogs"
	broker1Address = "158.160.10.60:9092"
)

func Produce(ctx context.Context, logMessage string) {
	partition := 0

	conn, err := kafka.DialLeader(context.Background(), "tcp", broker1Address, topic, partition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	_, err = conn.WriteMessages(
		kafka.Message{Value: []byte(logMessage)},
	)
	if err != nil {
		log.Fatal("failed to write messages:", err)
	}

	if err := conn.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}
}
