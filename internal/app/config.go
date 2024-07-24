package app

import (
	"encoding/json"
	"os"
)

type Config struct {
	Port        int      `json:"port"`
	PostgresURL string   `json:"postgres_url"`
	KafkaAddrs  []string `json:"kafka_addrs"`
}

func NewConfig(cfgPath string) (Config, error) {
	data, err := os.ReadFile(cfgPath)
	if err != nil {
		return Config{}, err
	}

	var cfg Config

	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return Config{}, err
	}

	if cfg.Port == 0 {
		cfg.Port = 8080
	}

	return cfg, nil
}
