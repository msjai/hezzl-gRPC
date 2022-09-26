package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"hezzl/auth"
	"log"
	"time"
)

type LogMessageType struct {
	TimeStamp string `json:"timestamp"`
	Message   string `json:"message"`
}

func Produce(ctx context.Context, logMessage string) {
	topic, err := auth.GetToken("#KafkaTopic")
	if err != nil {
		log.Fatal(err)
	}

	broker1Address, err := auth.GetToken("#KafkaBroker1Address")
	if err != nil {
		log.Fatal(err)
	}

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
