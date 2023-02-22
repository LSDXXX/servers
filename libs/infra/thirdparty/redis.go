package thirdparty

import (
	"context"
	"strings"

	"github.com/LSDXXX/libs/config"
	"github.com/go-redis/redis/v8"
)

// NewRedis create redis
//  @param conf
//  @return redis.Cmdable
//  @return error
func NewRedis(conf config.RedisConfig) (redis.Cmdable, error) {
	var client redis.Cmdable
	if conf.IsCluster {
		client = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    strings.Split(conf.Addr, ","),
			Password: conf.Password,
			Username: conf.Username,
		})
	} else {
		client = redis.NewClient(&redis.Options{
			Addr:     conf.Addr,
			Password: conf.Password,
			Username: conf.Username,
			DB:       conf.DB,
		})
	}
	cmd := client.Ping(context.Background())
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}

	return client, cmd.Err()
}
