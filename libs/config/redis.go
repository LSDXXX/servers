package config

// RedisConfig redis conf
type RedisConfig struct {
	Network   string `yaml:"network"`
	Addr      string `yaml:"addr"`
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
	DB        int    `yaml:"db"`
	IsCluster bool   `yaml:"is_cluster"`
}
