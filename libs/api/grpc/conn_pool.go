package grpc

import (
	"context"
	"sync"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ClientPool grpc client pool
type ClientPool struct {
	pool sync.Map
}

// Get get grpc client
//  @receiver c
//  @param dns
//  @return *grpc.ClientConn
//  @return error
func (c *ClientPool) Get(dns string) (*grpc.ClientConn, error) {
	if conn, ok := c.pool.Load(dns); ok {
		return conn.(*grpc.ClientConn), nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, dns,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(UnaryClientInterceptor),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`))
	if err != nil {
		return nil, errors.Wrap(err, "grpc conn")
	}
	c.pool.Store(dns, conn)
	return conn, nil
}

// Delete delete client
//  @receiver c
//  @param dns
func (c *ClientPool) Delete(dns string) {
	c.pool.Delete(dns)
}
