package grpc

import (
	"errors"
	"fmt"
	"net"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/LSDXXX/libs/pkg/container"
	"github.com/LSDXXX/libs/pkg/util"
	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc/resolver"
)

const (
	defaultPort = "8500"
)

var (
	errMissingAddr = errors.New("consul resolver: missing address")

	errAddrMisMatch = errors.New("consul resolver: invalied uri")

	errEndsWithColon = errors.New("consul resolver: missing port after port-separator colon")

	regexConsul, _ = regexp.Compile("^([A-z0-9.]+)(:[0-9]{1,5})?/([A-z_]+)$")
)

// RegisterConsulResolver register consul builder
func RegisterConsulResolver() {
	resolver.Register(NewBuilder())
}

type consulBuilder struct {
	consul *api.Client `container:"type"`
}

type consulResolver struct {
	wg                   sync.WaitGroup
	cc                   resolver.ClientConn
	consul               *api.Client
	name                 string
	disableServiceConfig bool
	ch                   chan int
	closech              chan struct{}
}

// NewBuilder create consul builder
//  @return resolver.Builder
func NewBuilder() resolver.Builder {
	out := &consulBuilder{}
	util.PanicWhenError(container.Fill(out))
	return out
}

func (cb *consulBuilder) Build(target resolver.Target, cc resolver.ClientConn,
	opts resolver.BuildOptions) (resolver.Resolver, error) {

	cr := &consulResolver{
		name:                 target.URL.Path,
		cc:                   cc,
		disableServiceConfig: opts.DisableServiceConfig,
		ch:                   make(chan int),
		closech:              make(chan struct{}),
		consul:               cb.consul,
	}
	go cr.watcher()
	cr.ResolveNow(resolver.ResolveNowOptions{})
	return cr, nil

}

func (cr *consulResolver) watcher() {
	client := cr.consul
	t := time.NewTicker(2000 * time.Millisecond)
	defer func() {
		fmt.Println("defer done")
	}()
	for {
		select {
		case <-t.C:
			//fmt.Println("定时")
		case <-cr.ch:
			//fmt.Println("ch call")
		case <-cr.closech:
			return
		}
		//api添加了 lastIndex   consul api中并不兼容附带lastIndex的查询
		services, _, err := client.Health().Service(cr.name, "", true, &api.QueryOptions{})
		if err != nil {
			fmt.Printf("error retrieving instances from Consul: %v", err)
		}

		newAddrs := make([]resolver.Address, 0)
		for _, service := range services {
			addr := net.JoinHostPort(service.Service.Address, strconv.Itoa(service.Service.Port))
			newAddrs = append(newAddrs, resolver.Address{
				Addr: addr,
				//type：不能是grpclb，grpclb在处理链接时会删除最后一个链接地址，不用设置即可 详见=> balancer_conn_wrappers => updateClientConnState
				ServerName: service.Service.Service,
			})
		}
		//cr.cc.NewAddress(newAddrs)
		//cr.cc.NewServiceConfig(cr.name)
		_ = cr.cc.UpdateState(resolver.State{Addresses: newAddrs})
	}

}

func (cb *consulBuilder) Scheme() string {
	return "consul"
}

func (cr *consulResolver) ResolveNow(opt resolver.ResolveNowOptions) {
	cr.ch <- 1
}

func (cr *consulResolver) Close() {
	close(cr.closech)
}
