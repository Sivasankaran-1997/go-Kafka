package kafkas

import (
	"context"
	"encoding/json"
	"kafka/domain"
	"kafka/utils"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/segmentio/kafka-go"
)

func Consume() (*domain.User, *utils.Resterr) {

	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	topic := os.Getenv("TOPIC")
	brokerAddress := os.Getenv("BROKERADDRESS")

	l := log.New(os.Stdout, "kafka reader: ", 0)
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{brokerAddress},
		Topic:   topic,
		GroupID: "my-group",
		Logger:  l,
	})

	msg, err := r.ReadMessage(context.TODO())
	if err != nil {
		return nil, utils.NotFound("Consumer Message Not found")
	}
	text := string(msg.Value)
	bytes := []byte(text)

	var user domain.User
	json.Unmarshal(bytes, &user)

	// after receiving the message, log its value
	return &user, nil
}
