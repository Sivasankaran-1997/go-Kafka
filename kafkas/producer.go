package kafkas

import (
	"context"
	"encoding/json"
	"kafka/domain"
	"kafka/utils"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/segmentio/kafka-go"
)

func Produce(user domain.User) *utils.Resterr {

	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	topic := os.Getenv("TOPIC")
	broker1Address := os.Getenv("BROKERADDRESS")

	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{broker1Address},
		Topic:   topic,
	})
	reqMarshal, err := json.Marshal(user)
	if err != nil {
		return utils.BadRequest("Invalid Request")
	}
	strTime := strconv.Itoa(int(time.Now().Unix()))
	reqString := string(reqMarshal)
	if err := w.WriteMessages(context.TODO(), kafka.Message{
		Key:   []byte(strTime),
		Value: []byte(reqString),
	}); err != nil {
		return utils.BadRequest("Kafka Server Not Connected")
	}
	return nil
}
