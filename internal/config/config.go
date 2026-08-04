package config

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type (
	Config struct {
		KafkaBrokers       []string
		KafkaTopic         string
		KafkaConsumerGroup string
	}
	topicConfig struct {
		TopicName string `json:"topic_name"`
		Brokers   string `json:"brokers"`
		Group     string `json:"group"`
	}
)

func Load(path string) (Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return Config{}, fmt.Errorf("open properties file: %w", err)
	}
	defer file.Close()

	values := make(map[string]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
			continue
		}
		key, value, found := strings.Cut(line, "=")
		if found {
			values[strings.TrimSpace(key)] = strings.TrimSpace(value)
		}
	}
	if err := scanner.Err(); err != nil {
		return Config{}, fmt.Errorf("read properties file: %w", err)
	}

	var topic topicConfig
	if err := json.Unmarshal([]byte(values["auth.topic.config"]), &topic); err != nil {
		return Config{}, fmt.Errorf("parse auth.topic.config: %w", err)
	}

	cfg := Config{
		KafkaBrokers:       splitCSV(topic.Brokers),
		KafkaTopic:         topic.TopicName,
		KafkaConsumerGroup: topic.Group,
	}
	if len(cfg.KafkaBrokers) == 0 || cfg.KafkaTopic == "" || cfg.KafkaConsumerGroup == "" {
		return Config{}, fmt.Errorf("incomplete Kafka configuration in %s", path)
	}

	return cfg, nil
}

func splitCSV(value string) []string {
	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		if part = strings.TrimSpace(part); part != "" {
			result = append(result, part)
		}
	}
	return result
}
