package kafka

import (
	"context"
	"encoding/json"
	"time"

	"github.com/segmentio/kafka-go"

	"agregator/preprocessor/internal/interfaces"
	kafka_model "agregator/preprocessor/internal/model/kafka"
)

type Kafka struct {
	writer *kafka.Writer
	reader *kafka.Reader
	logger interfaces.Logger
}

func New(brokers []string, topic string, logger interfaces.Logger) *Kafka {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
	return &Kafka{
		writer: writer,
		logger: logger,
	}
}

// Функция для чтения сообщений из Kafka
func (k *Kafka) StartReading(brokers []string, groupID, topic string, output chan<- kafka_model.Item) {
	k.reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		GroupID: groupID,
		Topic:   topic,
	})
	defer k.reader.Close()
	for {
		msg, err := k.reader.ReadMessage(context.Background())
		if err != nil {
			k.logger.Error("Error reading message from Kafka", "error", err)
			continue
		}
		item := kafka_model.Item{}
		err = json.Unmarshal(msg.Value, &item)
		if err != nil {
			k.logger.Error("Error decoding Kafka message", "error", err)
			continue
		}
		output <- item
	}
}

// Функция для записи сообщений в Kafka
func (k *Kafka) StartWriting(input <-chan kafka_model.Item) {
	for item := range input {
		data, err := json.Marshal(item)
		if err != nil {
			k.logger.Error("Error encoding item to JSON", "error", err)
			continue
		}
		message := kafka.Message{
			Key:   []byte(item.MD5), // Используем MD5 как ключ
			Value: data,
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		err = k.writer.WriteMessages(ctx, message)
		cancel()
		if err != nil {
			k.logger.Error("Error writing message to Kafka", "error", err)
		}
	}
}
