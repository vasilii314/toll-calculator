package main

import (
	"context"
	"encoding/json"
	"time"
	"tolling/client"
	"tolling/types"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/sirupsen/logrus"
)

// This can also be called KafkaTransport
type KafkaConsumer struct {
	consumer  *kafka.Consumer
	isRunning bool
	service   CalculatorServicer
	aggClient client.Client
}

func NewKafkaConsumer(topic string, svc CalculatorServicer, aggClient client.Client) (*KafkaConsumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		return nil, err
	}

	c.SubscribeTopics([]string{topic, "^aRegex.*[Tt]opic"}, nil)

	return &KafkaConsumer{
		consumer: c,
		service:  svc,
		aggClient: aggClient,
	}, nil
}

func (c *KafkaConsumer) readMessageLoop() {
	for c.isRunning {
		msg, err := c.consumer.ReadMessage(-1)
		if err != nil {
			logrus.Errorf("kafka consume error %s", err)
			continue
		}
		var data types.OBUData
		if err := json.Unmarshal(msg.Value, &data); err != nil {
			logrus.Errorf("JSON serialization error: %s", err)
			continue
		}
		distance, err := c.service.CalculateDistance(data)
		if err != nil {
			logrus.Errorf("calculation error: %s", err)
			continue
		}
		req := types.AggregateRequest {
			Value: distance,
			Unix: time.Now().Unix(),
			ObuId: int32(data.OBUID),
		}
		if err := c.aggClient.Aggregate(context.Background(), &req); err != nil {
			logrus.Errorf("aggregation error: %s", err.Error())
			continue
		}
	}
}

func (c *KafkaConsumer) Start() {
	logrus.Info("kafka transport started")
	c.isRunning = true
	c.readMessageLoop()
}
