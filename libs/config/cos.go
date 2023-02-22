package config

type CosConfig struct {
	Domain    string `yaml:"domain"`
	SecretId  string `yaml:"secret_id"`
	SecretKey string `yaml:"secret_key"`
}
