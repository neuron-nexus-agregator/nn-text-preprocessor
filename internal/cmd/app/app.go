package app

import (
	"os"
	"sync"
	"time"

	"agregator/preprocessor/internal/interfaces"
	"agregator/preprocessor/internal/service/kafka"
	"agregator/preprocessor/internal/service/preprocessor"
)

type App struct {
	preprocessor *preprocessor.Preprocessor
	kafka        *kafka.Kafka
}

func New(updateDBDuration time.Duration, logger interfaces.Logger) (*App, error) {
	preprocessor := preprocessor.New(30, logger)
	brokers := []string{os.Getenv("KAFKA_ADDR")}
	writeTopic := "filter"
	kafka := kafka.New(brokers, writeTopic, logger)
	return &App{
		preprocessor: preprocessor,
		kafka:        kafka,
	}, nil
}

func (a *App) Run() {
	preprocessorInput := a.preprocessor.Input()
	preprocessorOutput := a.preprocessor.Output()

	brokers := []string{os.Getenv("KAFKA_ADDR")}
	readTopic := "preprocessor"

	wg := sync.WaitGroup{}
	wg.Add(3)
	go func() {
		defer wg.Done()
		a.preprocessor.Start()
	}()
	go func() {
		defer wg.Done()
		a.kafka.StartReading(brokers, "news-processor-group", readTopic, preprocessorInput)
	}()
	go func() {
		defer wg.Done()
		a.kafka.StartWriting(preprocessorOutput)
	}()
	wg.Wait()
}
