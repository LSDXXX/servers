package config

import (
	"io/ioutil"
	"os"

	"github.com/creasty/defaults"
	"gopkg.in/yaml.v2"
)

// Config config
type Config struct {
	Log            LogConfig          `yaml:"log"`
	Kafka          KafkaConfig        `yaml:"kafka"`
	GinLog         LogConfig          `yaml:"gin_log"`
	Consul         ConsulConfig       `yaml:"consul"`
	ZK             ZKConfig           `yaml:"zk"`
	Mysql          MysqlConfig        `yaml:"mysql"`
	Redis          RedisConfig        `yaml:"redis"`
	Cos            CosConfig          `yaml:"cos"`
	TencentCloud   TencentCloudConfig `yaml:"tencent_cloud"`
	WorkerZKPath   string             `yaml:"worker_zk_path"`
	FlowZKPath     string             `yaml:"flow_zk_path"`
	HttpServerPort int                `yaml:"server_port"`
	GrpcServerPort int                `yaml:"grpc_server_port"`
	ServerName     string             `yaml:"server_name"`

	PrometheusBindURL string `yaml:"prometheus_bind_url"`
	Cron              string `yaml:"cron"`

	EnableWorkerGroup bool `yaml:"enable_worker_group"`
}

/*
func GetConfig() *Config {
	loadOnce.Do(func() {
		file, err := os.OpenFile(FilePath, os.O_RDONLY, 0755)
		if err != nil {
			panic("can't open config file, check the config file path, path: " + FilePath)
		}
		data, err := ioutil.ReadAll(file)
		if err != nil {
			panic("read config file error: " + err.Error())
		}
		defaults.Set(&conf)
		err = yaml.Unmarshal(data, &conf)
		if err != nil {
			panic("unmarshal config file error: " + err.Error())
		}
	})
	return &conf
}
*/

// ReadConfig read config from file
//  @param filePath
//  @param conf
//  @return error
func ReadConfig(filePath string, conf interface{}) error {
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0755)
	if err != nil {
		return err
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	_ = defaults.Set(conf)
	return yaml.Unmarshal(data, conf)
}
