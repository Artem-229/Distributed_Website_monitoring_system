package kafka

import (
	"fmt"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/sirupsen/logrus"
)

const (
	sessionTimeOut = 8000 // ms
	noTimeOut      = -1
)

type Handler interface {
	HandleMessage(msg []byte, offset kafka.Offset) error
}

type Consumer struct {
	consumer *kafka.Consumer
	handler  Handler
	stop     bool
}

func NewConsumer(address []string, topic, groupConsumers string, handler Handler) (*Consumer, error) {

	conf := kafka.ConfigMap{
		"bootstrap.servers":  strings.Join(address, ","),
		"group.id":           groupConsumers,
		"session.timeout.ms": sessionTimeOut,
		"auto.offset.reset":  "earliest",
	}

	c, err := kafka.NewConsumer(&conf)

	if err != nil {
		return nil, fmt.Errorf("error with connecting to host %w", err)
	}

	if err = c.Subscribe(topic, nil); err != nil {
		return nil, err
	}

	return &Consumer{
		consumer: c,
		handler:  handler,
		stop:     false,
	}, err
}

func (c *Consumer) Start() {
	for {
		if c.stop == true {
			break
		}
		kafkaMSg, err := c.consumer.ReadMessage(noTimeOut)
		if err != nil {
			logrus.Error(err)
		}
		if kafkaMSg == nil {
			continue
		}
		if err := c.handler.HandleMessage(kafkaMSg.Value, kafkaMSg.TopicPartition.Offset); err != nil {
			logrus.Error(err)
		}
	}
}

func (c *Consumer) Stop() error {
	c.stop = true
	c.consumer.Commit()
	return c.consumer.Close()
}
