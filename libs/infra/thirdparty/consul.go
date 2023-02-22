package thirdparty

import (
	"context"
	"errors"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/LSDXXX/libs/config"
	"github.com/LSDXXX/libs/pkg/container"
	"github.com/LSDXXX/libs/pkg/log"
	"github.com/LSDXXX/libs/pkg/singleflight"
	"github.com/LSDXXX/libs/pkg/util"
	"github.com/hashicorp/consul/api"
	"github.com/robfig/cron/v3"
	"github.com/spf13/cast"
)

var (
	consulSingleGroup singleflight.Group
)

type address struct {
	ip   string
	port int
	addr string
}

func (a *address) Network() string {
	return ""
}

func (a *address) String() string {
	if len(a.addr) == 0 {
		a.addr = a.ip + ":" + cast.ToString(a.port)
	}
	return a.addr
}

type consulClient struct {
	conn *api.Client
}

func (c *consulClient) GetAddress(ctx context.Context, name string) (net.Addr, error) {
	addr, err, _ := consulSingleGroup.Do(name, func() (interface{}, error) {
		var nilAddr *address
		services, _, err := c.conn.Health().Service(name, "", true, nil)
		if err != nil {
			return nilAddr, err
		}
		if len(services) > 0 {
			service := services[rand.Int()%len(services)]
			ip := strings.Split(service.Service.Address, ":")[0]
			return &address{
				ip:   ip,
				port: service.Service.Port,
				addr: ip + ":" + cast.ToString(service.Service.Port),
			}, nil
		}
		return nilAddr, errors.New("consul " + name + " not found")
	})
	return addr.(*address), err
}

func (c *consulClient) register(serverName string, serverPort int) error {
	address := util.GetLocalIP()
	host := address + ":" + strconv.Itoa(serverPort)
	ttl := time.Second * 4

	check := &api.AgentServiceCheck{
		// HTTP:                           "http://" + host + "/health",
		TTL:     (ttl + time.Second).String(),
		Timeout: (ttl * 2).String(),
		Status:  api.HealthPassing,
	}
	serviceId := serverName + "-" + host

	registration := new(api.AgentServiceRegistration)
	registration.Name = serverName
	registration.ID = serviceId
	registration.Port = serverPort
	registration.Address = address
	registration.Check = check

	err := c.conn.Agent().ServiceRegister(registration)
	if err != nil {
		return err
	}
	var cron *cron.Cron
	container.Resolve(&cron)
	checkId := "service:" + serviceId
	_ = c.conn.Agent().UpdateTTL(checkId, "", check.Status)
	cron.AddFunc("@every "+ttl.String(), func() {
		err = c.conn.Agent().UpdateTTL(checkId, "", check.Status)
		if err != nil {
			log.WithContext(context.Background()).Errorf("update consul ttl error: %+v", err)
		}
	})
	return nil
}

// NewConsulDiscovery new discovery
//  @param conf
//  @param serverName
//  @param serverPort
//  @param grpcPort
//  @return *consulClient
//  @return error
func NewConsulDiscovery(conf config.ConsulConfig, serverName string, serverPort int, grpcPort int) (*consulClient, error) {
	client := &consulClient{}
	consulConf := api.DefaultConfig()
	if len(conf.Address) > 0 {
		consulConf.Address = conf.Address
	}
	if len(conf.Scheme) > 0 {
		consulConf.Scheme = conf.Scheme
	}
	conn, err := api.NewClient(consulConf)
	if err != nil {
		return nil, err
	}
	client.conn = conn
	container.Singleton(func() *api.Client {
		return conn
	})

	err = client.register(serverName, serverPort)
	if err != nil {
		return nil, err
	}

	err = client.register(serverName+"-grpc", grpcPort)
	if err != nil {
		return nil, err
	}
	return client, nil
}
