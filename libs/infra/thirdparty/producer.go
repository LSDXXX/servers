package thirdparty

import (
	"time"

	"github.com/LSDXXX/libs/config"
	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
)

type kafkaProducer struct {
	producer sarama.AsyncProducer
}

func (p *kafkaProducer) ProduceMessage(topic string, message []byte) error {
	input := &sarama.ProducerMessage{
		Value: sarama.ByteEncoder(message),
		Topic: topic,
	}

	p.producer.Input() <- input
	return nil
}

func (p *kafkaProducer) ProduceMessageWithKey(topic string, key []byte, message []byte) error {
	input := &sarama.ProducerMessage{
		Value: sarama.ByteEncoder(message),
		Topic: topic,
	}

	if len(key) > 0 {
		input.Key = sarama.ByteEncoder(key)
	}

	p.producer.Input() <- input
	return nil
}

// NewKafkaProducer create producer
//  @param kconf
//  @return *kafkaProducer
//  @return error
func NewKafkaProducer(kconf config.KafkaConfig) (*kafkaProducer, error) {

	conf := sarama.NewConfig()

	conf.Producer.RequiredAcks = sarama.WaitForLocal
	conf.Producer.Flush.Frequency = 500 * time.Millisecond

	producer, err := sarama.NewAsyncProducer(kconf.Brokers, conf)

	if err != nil {
		return nil, err
	}

	go func() {
		for err := range producer.Errors() {
			logrus.Errorf("%v: %+v", err.Err, err.Msg)
		}
	}()

	return &kafkaProducer{producer: producer}, nil
}
