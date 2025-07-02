package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"order_service/infra/log"
	"order_service/models"
)

type KafkaProducer struct {
	writer *kafka.Writer
}

func NewKafkaProducer(brokers []string, topic string) *KafkaProducer {
	writer := &kafka.Writer{
		Topic:    topic,
		Addr:     kafka.TCP(brokers...),
		Balancer: &kafka.LeastBytes{},
	}

	return &KafkaProducer{
		writer: writer,
	}
}

func (k *KafkaProducer) Close() error {
	return k.writer.Close()
}

func (k *KafkaProducer) PublishOrderCreated(ctx context.Context, event interface{}) error {
	value, err := json.Marshal(event)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"err":   err.Error(),
			"event": event,
		})
		return err
	}

	msg := kafka.Message{
		Key:   []byte(fmt.Sprintf("order-%d", event.(models.OrderCreatedEvent).OrderID)),
		Value: value,
	}
	return k.writer.WriteMessages(ctx, msg)
}
