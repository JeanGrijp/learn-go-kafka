package producer

import (
	"log"

	"github.com/IBM/sarama"
)

func SendCepToKafka(cep string) error {
	producer, err := sarama.NewSyncProducer([]string{"localhost:9092"}, nil)
	if err != nil {
		log.Fatalln("Failed to start producer:", err)
	}
	defer func() {
		if err := producer.Close(); err != nil {
			log.Fatalln("Failed to close producer:", err)
		}
	}()

	message := &sarama.ProducerMessage{
		Topic: "test-topic",
		Value: sarama.StringEncoder("Hello, Kafka!"),
	}

	partition, offset, err := producer.SendMessage(message)
	if err != nil {
		log.Fatalln("Failed to send message:", err)
	}

	log.Printf("Message sent! Partition = %d, Offset = %d\n", partition, offset)

	return nil
}
