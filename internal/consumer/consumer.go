package consumer

import (
	"fmt"
	"log"

	"github.com/IBM/sarama"
)

func StartCepConsumer() {
	consumer, err := sarama.NewConsumer([]string{"localhost:9092"}, nil)
	if err != nil {
		log.Fatalln("Failed to start consumer:", err)
	}
	defer func() {
		if err := consumer.Close(); err != nil {
			log.Fatalln("Failed to close consumer:", err)
		}
	}()

	partitionConsumer, err := consumer.ConsumePartition("test-topic", 0, sarama.OffsetOldest)
	if err != nil {
		log.Fatalln("Failed to start partition consumer:", err)
	}
	defer func() {
		if err := partitionConsumer.Close(); err != nil {
			log.Fatalln("Failed to close partition consumer:", err)
		}
	}()

	for message := range partitionConsumer.Messages() {
		fmt.Printf("Received message: Key = %s, Value = %s\n", string(message.Key), string(message.Value))
	}
}
