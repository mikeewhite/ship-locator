package kafka

import "time"

type Metrics interface {
	KafkaConsumeTime(topic string, startTime time.Time)
}
