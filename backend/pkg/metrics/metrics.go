package metrics

import (
	"time"
)

func (c *Client) DBQueryTime(queryName string, startTime time.Time) {
	c.dbQueryTimeHistogram.WithLabelValues(queryName).Observe(time.Since(startTime).Seconds())
}

func (c *Client) KafkaConsumeTime(topic string, startTime time.Time) {
	c.kafkaConsumeTimeHistogram.WithLabelValues(topic).Observe(time.Since(startTime).Seconds())
}
