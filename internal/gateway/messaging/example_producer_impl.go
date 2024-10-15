package messaging

import (
	"user-simple-crud/internal/model"
	kafkaserver "user-simple-crud/pkg/broker/kafkaservice"
)

type ExampleProducerImpl struct {
	ProducerKafka[*model.ExampleMessage]
}

func NewExampleKafkaProducerImpl(producer *kafkaserver.KafkaService, topic string) ExampleProducer {
	return &ExampleProducerImpl{
		ProducerKafka: ProducerKafka[*model.ExampleMessage]{
			Topic:         topic,
			KafkaProducer: producer,
		},
	}
}
