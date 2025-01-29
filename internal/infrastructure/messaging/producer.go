package messaging

import (
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(broker string, topic string) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(broker),
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func (p *Producer) Publish(event interface{}) error {
	message, err := json.Marshal(event)
	if err != nil {
		return err
	}

	err = p.writer.WriteMessages(nil, kafka.Message{Value: message})
	if err != nil {
		log.Printf("Error al publicar evento: %v", err)
		return err
	}

	log.Printf("Evento publicado: %s", string(message))
	return nil
}
