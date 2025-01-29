package messaging

import (
	"context"
	"encoding/json"
	"log"
	"nilus-challenge-backend/internal/domain/notification"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafka.Reader
}

func NewConsumer(broker string, topic string, groupID string) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{broker},
			Topic:   topic,
			GroupID: groupID,
		}),
	}
}

func (c *Consumer) StartProcessing(handler func(event notification.NotificationEvent) error) error {
	for {
		msg, err := c.reader.ReadMessage(context.Background())
		if err != nil {
			return err
		}

		var event notification.NotificationEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("Error al deserializar evento: %v", err)
			continue
		}

		log.Printf("Evento recibido: %v", event)

		if err := handler(event); err != nil {
			log.Printf("Error al procesar evento: %v", err)
		}
	}
}
