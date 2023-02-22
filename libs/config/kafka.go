package config

// KafkaConfig kafka config
type KafkaConfig struct {
	Brokers []string            `yaml:"brokers"`
	Version string              `yaml:"version" default:"2.4.0"`
	Group   ConsumerGroupConfig `yaml:"consumer"`
}

// ConsumerGroupConfig config
type ConsumerGroupConfig struct {
	Topics  []string `yaml:"topics"`
	GroupID string   `yaml:"group_id"`
}
