package config

// ZKConfig zk conf
type ZKConfig struct {
	Hosts  []string `yaml:"hosts"`
	Scheme string   `yaml:"scheme"`
	Auth   string   `yaml:"auth"`
}
