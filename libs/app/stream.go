package app

import (
	"context"

	"github.com/LSDXXX/libs/api"
	"github.com/LSDXXX/libs/config"
	"github.com/LSDXXX/libs/pkg/util"
	"github.com/Shopify/sarama"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

// Streaming stream interface
type Streaming interface {
	SetHandler(api.StreamMessageHandler)

	Start(context.Context)
}

type kafkaConsumer struct {
	handlers map[string]api.StreamMessageHandler
	ready    chan bool
	client   sarama.ConsumerGroup
	topics   []string
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *kafkaConsumer) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(consumer.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *kafkaConsumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (consumer *kafkaConsumer) ConsumeClaim(session sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim) error {
	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/Shopify/sarama/blob/main/consumer_group.go#L27-L29
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorf("process consumed msg panic: %v", r)
		}
	}()
	for message := range claim.Messages() {
		if h, ok := consumer.handlers[message.Topic]; ok {
			logrus.Infof("consume message, data: %s", util.UnsafeString(message.Value))
			err := h.Process(message.Value)
			if err != nil {
				log.Errorf("process message error: %v, msg: %s, topic: %s",
					err, string(message.Value), message.Topic)
			}
		}
		session.MarkMessage(message, "")
	}

	return nil
}

func (consumer *kafkaConsumer) Start(ctx context.Context) {
	for {
		if err := consumer.client.Consume(ctx, consumer.topics, consumer); err != nil {
			log.Errorf("consume msg error: %s", err.Error())
		}
		if ctx.Err() != nil {
			return
		}
		consumer.ready = make(chan bool)
	}
}

func (consumer *kafkaConsumer) SetHandler(h api.StreamMessageHandler) {
	if len(h.Topic()) == 0 {
		return
	}
	if util.SliceIndex(consumer.topics, h.Topic()) == -1 {
		consumer.topics = append(consumer.topics, h.Topic())
	}
	consumer.handlers[h.Topic()] = h
}

// NewKafkaConsumer new consumer
//  @param kconf
//  @return Streaming
//  @return error
func NewKafkaConsumer(kconf config.KafkaConfig) (Streaming, error) {
	conf := sarama.NewConfig()
	// version, err := sarama.ParseKafkaVersion(kconf.Version)
	// version, err := sarama.ParseKafkaVersion("3.1.0")
	// if err != nil {
	// 	// panic("parse kafka version err: " + err.Error())
	// 	return nil, err
	// }
	// conf.Version = version
	if len(kconf.Group.GroupID) == 0 {
		kconf.Group.GroupID = uuid.New().String()
	}
	// conf.Consumer.Offsets.Initial = sarama.OffsetOldest
	conf.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	client, err := sarama.NewConsumerGroup(kconf.Brokers, kconf.Group.GroupID, conf)
	if err != nil {
		// panic("create consumer group error: " + err.Error())
		return nil, errors.WithMessage(err, "create consumer group")
	}

	return &kafkaConsumer{
		ready:    make(chan bool),
		handlers: make(map[string]api.StreamMessageHandler),
		topics:   kconf.Group.Topics,
		client:   client,
	}, nil
}
