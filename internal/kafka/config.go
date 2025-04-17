package kafka

import "github.com/IBM/sarama"

func NewKafkaConfig() *sarama.Config {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Consumer.Return.Errors = true
	return config
}

func GetBrokers() []string {
	return []string{"kafka:9092"} // ajuste se quiser pegar de env
}
