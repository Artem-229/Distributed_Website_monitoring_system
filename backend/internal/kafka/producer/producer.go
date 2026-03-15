package producer

import (
	"errors"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

const (
	flushtimeout = 8000
)

var unknowntype = errors.New("unkown type of error")

type Producer struct {
	producer *kafka.Producer
}

func NewProducer(address []string) (*Producer, error) {
	conf := kafka.ConfigMap{
		"bootstrap.servers": strings.Join(address, ","),
	}

	p, err := kafka.NewProducer(&conf)

	if err != nil {
		return nil, err
	}

	return &Producer{
		producer: p,
	}, nil

}

func (p *Producer) Produce(message, topic string) error {
	kafkaMsg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value: []byte(message),
	}

	kafkachan := make(chan kafka.Event)
	p.producer.Produce(kafkaMsg, kafkachan)

	e := <-kafkachan

	switch ev := e.(type) {
	case *kafka.Message:
		return nil
	case *kafka.Error:
		return ev
	default:
		return unknowntype
	}
}

func (p *Producer) Close() {
	p.producer.Flush(flushtimeout)
	p.producer.Close()
}
