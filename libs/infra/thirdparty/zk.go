package thirdparty

import (
	"time"

	"github.com/LSDXXX/libs/config"
	"github.com/go-zookeeper/zk"
)

// NewZkClient create zk
//  @param conf
//  @return *zk.Conn
//  @return error
func NewZkClient(conf config.ZKConfig) (*zk.Conn, error) {
	conn, _, err := zk.Connect(conf.Hosts, time.Second*10)
	if err != nil {
		return nil, err
	}
	if len(conf.Auth) != 0 {
		err = conn.AddAuth(conf.Scheme, []byte(conf.Auth))
		return nil, err
	}
	return conn, nil
}
