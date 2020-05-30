package db

import (
	"FBI-Analyzer/logger"
	"FBI-Analyzer/rule"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Kafka struct {
	Broker  string
	GroupID string
	Topic   []string
	Offset  string
}

// 消费kafka并解析成access结构体
func (k *Kafka) Consumer(kchan chan rule.AccessLog, i int) {

	logger.Print(logger.INFO, "consumer thread no.%d started.", i)
	var access rule.AccessLog
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": k.Broker,
		"group.id":          k.GroupID,
		"auto.offset.reset": k.Offset,
	})
	if err != nil {
		logger.Print(logger.ERROR, "kafka new FBI-Analyzer error : %v", err)
	}
	c.SubscribeTopics(k.Topic, nil)

	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			json.Unmarshal(msg.Value, &access)
			kchan <- access
		} else {
			// The client will automatically try to recover from all errors.
			logger.Print(logger.ERROR, "Consumer error: %v (%v)\n", err, msg)
		}
	}

	c.Close()
}
