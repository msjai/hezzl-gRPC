package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"time"
)

const (
	topic          = "HezzlLogs"
	broker1Address = "158.160.10.60:9092"
)

type LogJsonMessage struct {
	Message string `json:"message"`
}

func Produce(ctx context.Context, logMessage string) {
	partition := 0

	conn, err := kafka.DialLeader(context.Background(), "tcp", broker1Address, topic, partition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))

	LogMsg := LogJsonMessage{
		Message: logMessage,
	}
	LogMsgJS, err := json.Marshal(LogMsg)

	_, err = conn.WriteMessages(
		/*kafka.Message{Value: []byte(json_data2)},*/
		kafka.Message{Value: LogMsgJS},
	)
	if err != nil {
		log.Fatal("failed to write messages:", err)
	}

	if err := conn.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}
}

func Consume() {
	topic := "HezzlLogs"
	partition := 0

	conn, err := kafka.DialLeader(context.Background(), "tcp", "158.160.10.60:9092", topic, partition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	batch := conn.ReadBatch(0, 1e6) // fetch 10KB min, 1MB max

	b := make([]byte, 10e3) // 10KB max per message
	for {
		n, err := batch.Read(b)
		if err != nil {
			break
		}
		fmt.Println(string(b[:n]))
	}

	if err := batch.Close(); err != nil {
		log.Fatal("failed to close batch:", err)
	}

	if err := conn.Close(); err != nil {
		log.Fatal("failed to close connection:", err)
	}
}
