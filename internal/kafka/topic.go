package kafka

import (
	"log"

	"github.com/IBM/sarama"
)

func CreateTopic(brokers []string, topic string, partitions int32, replicationFactor int16) error {
	config := sarama.NewConfig()
	admin, err := sarama.NewClusterAdmin(brokers, config)
	if err != nil {
		return err
	}
	defer admin.Close()

	err = admin.CreateTopic(topic, &sarama.TopicDetail{
		NumPartitions:     partitions,
		ReplicationFactor: replicationFactor,
	}, false)

	if err != nil {
		log.Printf("Erro ao criar tópico: %v", err)
	} else {
		log.Printf("Tópico '%s' criado com sucesso", topic)
	}

	return err
}

func DeleteTopic(brokers []string, topic string) error {
	config := sarama.NewConfig()
	admin, err := sarama.NewClusterAdmin(brokers, config)
	if err != nil {
		return err
	}
	defer admin.Close()

	err = admin.DeleteTopic(topic)
	if err != nil {
		log.Printf("Erro ao deletar tópico: %v", err)
	} else {
		log.Printf("Tópico '%s' deletado com sucesso", topic)
	}

	return err
}
