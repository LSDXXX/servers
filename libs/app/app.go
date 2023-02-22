package app

import (
	"context"
	"errors"
	"sync"

	"github.com/LSDXXX/libs/api"
	"github.com/LSDXXX/libs/config"
	"github.com/LSDXXX/libs/pkg/container"
	"github.com/LSDXXX/libs/pkg/log"
	"github.com/google/uuid"
)

// Start app start
//  @param ctx
//  @return error
func Start(ctx context.Context) error {
	var wg sync.WaitGroup
	var conf *config.Config
	container.Resolve(&conf)
	routers := api.GetHttpRouters()
	if len(routers) > 0 {
		server := NewHttpServer(&conf.GinLog, api.GetHttpRouters()...)
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := server.Start(conf.HttpServerPort)
			if err != nil {
				log.WithContext(ctx).
					Panicf("start gin http server error: %s", err.Error())
				panic(err)
			}
		}()
	}
	log.WithContext(ctx).
		Infof("start gin http server success, router len: %d", len(routers))

	streamHandlers := api.GetStreamMessageHandler()
	if len(streamHandlers) > 0 {
		serverGroup := conf.Kafka.Group.GroupID
		if len(serverGroup) == 0 {
			serverGroup = uuid.NewString()
		}
		groups := make(map[string][]api.StreamMessageHandler)
		for _, handler := range streamHandlers {
			if len(handler.GroupID()) == 0 {
				groups[serverGroup] = append(groups[serverGroup], handler)
			} else {
				groups[handler.GroupID()] = append(groups[handler.GroupID()], handler)
			}
		}
		kafkaConf := conf.Kafka
		for groupID, handlers := range groups {
			kafkaConf.Group.GroupID = groupID
			var topics []string
			for _, handler := range handlers {
				topics = append(topics, handler.Topic())
			}
			kafkaConf.Group.Topics = topics
			s, err := NewKafkaConsumer(kafkaConf)
			if err != nil {
				return err
			}
			for _, handler := range handlers {
				s.SetHandler(handler)
				log.WithContext(ctx).
					Infof("add stream handler, group: %s, topic: %s", groupID, handler.Topic())
			}
			wg.Add(1)
			go func() {
				defer wg.Done()
				s.Start(ctx)
			}()
		}
	}

	grpcServices := api.GetGrpcServices()
	if len(grpcServices) > 0 {
		server := NewGrpcServer(grpcServices)
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := server.Start(conf.GrpcServerPort)
			if err != nil {
				log.WithContext(ctx).
					Panicf("start grpc server error: %s", err.Error())
				panic(err)
			}
		}()
	}
	log.WithContext(ctx).
		Infof("start grpc server success, services len: %d", len(grpcServices))

	wg.Wait()
	return errors.New("all app stopped")
}
