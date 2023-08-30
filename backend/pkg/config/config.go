package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	WebSocketAPIKey string

	GraphQLAddress string `default:":8085"`

	KafkaAddress       string `default:"localhost:9092"`
	KafkaTopic         string `default:"ship-data-topic"`
	KafkaConsumerGroup string `default:"consumer-group-1"`

	PostgresUsername string `default:"postgres"`
	PostgresPassword string `default:"postgres"`
	PostgresAddress  string `default:"localhost:5432"`
	PostgresDBName   string `default:"ship_db"`

	PrometheusServerAddress string `default:":2112"`

	TracingCollectorAddress string `default:"localhost:4318"`
}

func Load() (*Config, error) {
	var cfg Config
	err := envconfig.Process("SHIPLOC", &cfg)
	if err != nil {
		return nil, fmt.Errorf("error on reading config env vars :%w", err)
	}
	return &cfg, nil
}
