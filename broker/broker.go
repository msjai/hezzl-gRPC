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

type LogMessageType struct {
	TimeStamp string `json:"timestamp"`
	Message   string `json:"message"`
}

func Produce(ctx context.Context, logMessage string) {
	partition := 0

	conn, err := kafka.DialLeader(ctx, "tcp", broker1Address, topic, partition)
	if err != nil {
		log.Print("failed to dial leader:", err)
	}

	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))

	t := time.Now()
	formatted := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
	LogMsg := LogMessageType{
		Message:   logMessage,
		TimeStamp: formatted,
	}
	LogMsgJson, err := json.Marshal(LogMsg)

	_, err = conn.WriteMessages(
		kafka.Message{Value: LogMsgJson},
	)
	if err != nil {
		log.Print("failed to write messages:", err)
	}

	if err := conn.Close(); err != nil {
		log.Print("failed to close writer:", err)
	}
}
