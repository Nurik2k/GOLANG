package main

import (
	"fmt"
	"github.com/Shopify/sarama"
	"log"
	"os"
	"os/signal"
	"sync"
)

const (
	// Укажите явно IP-адрес и порт брокера Kafka
	kafkaConn = "localhost:9092"
	topic     = "test_topic"
)

func main() {
	// Настройка Sarama логирования
	sarama.Logger = log.New(os.Stdout, "", log.Ltime)

	// Создание конфигурации потребителя
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	// Создание клиента Kafka
	client, err := sarama.NewClient([]string{kafkaConn}, config)
	if err != nil {
		log.Fatalf("Error creating Kafka client: %v", err)
	}
	defer client.Close()

	// Создание потребителя
	consumer, err := sarama.NewConsumerFromClient(client)
	if err != nil {
		log.Fatalf("Error creating Kafka consumer: %v", err)
	}
	defer consumer.Close()

	// Указание темы для потребления
	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetOldest)
	if err != nil {
		log.Fatalf("Error creating partition consumer: %v", err)
	}
	defer partitionConsumer.Close()

	// Создание канала для обработки сигнала завершения
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, os.Interrupt)

	// Создание WaitGroup для ожидания завершения работы потребителя
	var wg sync.WaitGroup
	wg.Add(1)

	// Основной цикл потребителя
	go func() {
		defer wg.Done()
		for {
			select {
			case msg := <-partitionConsumer.Messages():
				fmt.Printf("Received message: %s\n", string(msg.Value))
			case err := <-partitionConsumer.Errors():
				log.Printf("Error: %s\n", err.Error())
			case <-sigterm:
				return
			}
		}
	}()

	// Ожидание завершения работы потребителя
	wg.Wait()
}
