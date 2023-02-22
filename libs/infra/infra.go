package infra

import (
	"context"
	"net"

	"github.com/LSDXXX/libs/config"
	"github.com/LSDXXX/libs/constant"
	"github.com/LSDXXX/libs/infra/thirdparty"
	"github.com/LSDXXX/libs/pkg/container"
	"github.com/LSDXXX/libs/pkg/util"
	"github.com/go-redis/redis/v8"
	"github.com/go-zookeeper/zk"
	cron "github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

var initFuncList []func()

func AppendInitFunc(f func()) {
	initFuncList = append(initFuncList, f)
}

// InfraOptions options
//
//	@param *infraOpts
type InfraOptions func(*infraOpts)

type infraOpts struct {
	withRedis         bool
	withDB            bool
	withZK            bool
	withKafkaProducer bool
	withConsul        bool
}

// WithRedis redis opt
//
//	@return InfraOptions
func WithRedis() InfraOptions {
	return func(o *infraOpts) {
		o.withRedis = true
	}
}

// WithDB db opt
//
//	@return InfraOptions
func WithDB() InfraOptions {
	return func(o *infraOpts) {
		o.withDB = true
	}
}

// WithZK zk opt
//
//	@return InfraOptions
func WithZK() InfraOptions {
	return func(o *infraOpts) {
		o.withZK = true
	}
}

// WithKafkaProducer producer
//
//	@return InfraOptions
func WithKafkaProducer() InfraOptions {
	return func(o *infraOpts) {
		o.withKafkaProducer = true
	}
}

// WithConsul consul
//
//	@return InfraOptions
func WithConsul() InfraOptions {
	return func(o *infraOpts) {
		o.withConsul = true
	}
}

// Init init
//
//	@param opts
//	@return error
func Init(opts ...InfraOptions) error {
	var conf *config.Config

	util.PanicWhenError(container.Resolve(&conf))
	thirdparty.SetupDatabase(conf.Mysql)

	c := cron.New()
	c.Start()
	_ = container.Singleton(func() *cron.Cron {
		return c
	})
	var o infraOpts
	for _, opt := range opts {
		opt(&o)
	}
	if o.withDB {
		db, err := thirdparty.NewMysqlDB(conf.Mysql)
		if err != nil {
			return err
		}
		_ = container.Singleton(func() *gorm.DB {
			return db
		})
	}

	if o.withConsul {
		consul, err := thirdparty.NewConsulDiscovery(conf.Consul,
			conf.ServerName, conf.HttpServerPort, conf.GrpcServerPort)
		if err != nil {
			return err
		}
		_ = container.Singleton(func() Discovery {
			return consul
		})
	}

	if o.withRedis {
		db, err := thirdparty.NewRedis(conf.Redis)
		if err != nil {
			return err
		}
		_ = container.Singleton(func() redis.Cmdable {
			return db
		})
	}

	if o.withKafkaProducer {
		producer, err := thirdparty.NewKafkaProducer(conf.Kafka)
		if err != nil {
			return err
		}
		_ = container.Singleton(func() Producer {
			return producer
		})
	}

	if o.withZK {
		c, err := thirdparty.NewZkClient(conf.ZK)
		if err != nil {
			return err
		}
		err = util.CreateZKPathP(context.Background(), c, constant.ZKLockPath)
		if err != nil {
			return err
		}
		_ = container.Singleton(func() *zk.Conn {
			return c
		})
	}

	for _, f := range initFuncList {
		f()
	}

	return nil
}

// Producer .
type Producer interface {
	ProduceMessage(topic string, message []byte) error
	ProduceMessageWithKey(topic string, key, message []byte) error
}

type MockProducer struct {
}

func (p *MockProducer) ProduceMessage(topic string, message []byte) error {
	return nil
}
func (p *MockProducer) ProduceMessageWithKey(topic string, key, message []byte) error {
	return nil
}

// Infrastructure .
type Infrastructure interface {
	DB() *gorm.DB
	ZK() *zk.Conn
	Producer() Producer
	Discovery() Discovery
	Redis() redis.Cmdable
}

// Discovery .
type Discovery interface {
	GetAddress(context.Context, string) (net.Addr, error)
}
