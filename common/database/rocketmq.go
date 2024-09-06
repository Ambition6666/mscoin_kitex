package database

import (
	"context"
	"errors"
	"fmt"
	rmq "github.com/apache/rocketmq-clients/golang/v5"
	"github.com/cloudwego/kitex/pkg/klog"
	"log"
	"sync"
	"time"
)

var (
	// maximum waiting time for receive func
	awaitDuration = time.Second * 5
	// maximum number of messages received at one time
	maxMessageNum int32 = 16
	// invisibleDuration should > 20s
	invisibleDuration = time.Second * 20
	// repeated number of send the message
	retries = 3
)

type RocketMQProducer struct {
	w      rmq.Producer
	data   chan RocketMQData
	closed bool
	mutex  sync.Mutex
}

type RocketMQData struct {
	Topic string
	Key   []string
	Data  []byte
}

func NewRocketMQProducer(config *rmq.Config, cap int) (*RocketMQProducer, error) {
	producer, err := rmq.NewProducer(config)
	if err != nil {
		return nil, err
	}
	err = producer.Start()
	if err != nil {
		return nil, err
	}
	return &RocketMQProducer{
		w:    producer,
		data: make(chan RocketMQData, cap),
	}, nil
}

func (p *RocketMQProducer) StartWrite() {
	go p.sendRocketMQ()
}

func (p *RocketMQProducer) Send(message RocketMQData) {
	defer func() {
		if err := recover(); err != nil {
			p.closed = true
		}
	}()
	p.data <- message
	p.closed = false
}

func (p *RocketMQProducer) Close() {
	if p.w != nil {
		p.w.GracefulStop()
		p.mutex.Lock()
		defer p.mutex.Unlock()
		if !p.closed {
			close(p.data)
			p.closed = true
		}
	}
}

func (p *RocketMQProducer) sendRocketMQ() {
	for {
		if p.closed {
			break
		}
		select {
		case data := <-p.data:
			message := &rmq.Message{
				Topic: data.Topic,
				Body:  data.Data,
			}
			message.SetKeys(data.Key...)

			var err error

			success := false
			for i := 0; i < retries; i++ {
				// attempt to create topic prior to publishing the message
				_, err = p.w.Send(context.TODO(), message)
				if err == nil {
					success = true
					break
				}
				if errors.Is(err, rmq.ErrNoAvailableBrokers) || errors.Is(err, context.DeadlineExceeded) {
					time.Sleep(time.Millisecond * 250)
					success = false
					continue
				}
				if err != nil {
					success = false
					log.Printf("kafka send writemessage err %s \n", err.Error())
				}
			}
			if !success {
				//重新放进去等待消费
				p.Send(data)
			}
		}
	}
}

type topicInfo struct {
	consumer rmq.SimpleConsumer
	data     chan RocketMQData
	closed   bool
	mutex    sync.Mutex
}

type RocketMQConsumer struct {
	r     map[string]*topicInfo
	mutex sync.RWMutex
}

func NewRocketMQConsumer() *RocketMQConsumer {
	return &RocketMQConsumer{
		r: make(map[string]*topicInfo),
	}
}

func (c *RocketMQConsumer) AddConsumer(config *rmq.Config, cap int, topic string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	consumer, err := rmq.NewSimpleConsumer(config, rmq.WithAwaitDuration(awaitDuration))
	if err != nil {
		return err
	}

	err = consumer.Start()
	if err != nil {
		return err
	}

	err = consumer.Subscribe(topic, rmq.SUB_ALL)
	if err != nil {
		return err
	}

	c.r[topic] = &topicInfo{
		consumer: consumer,
		closed:   false,
		data:     make(chan RocketMQData, cap),
	}
	return nil
}

func (c *RocketMQConsumer) StartRead(topic string) error {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if c.r[topic] == nil {
		return errors.New("rocketmq consumer not exist")
	}

	go c.readMsg(topic)
	return nil
}

func (c *RocketMQConsumer) readMsg(topic string) {
	for {
		c.mutex.RLock()
		topicInfo := c.r[topic]
		c.mutex.RUnlock()

		if topicInfo.closed {
			close(topicInfo.data)
			return
		}

		mvs, err := topicInfo.consumer.Receive(context.TODO(), maxMessageNum, invisibleDuration)
		if err != nil {
			klog.Error(err)
			continue
		}

		for _, mv := range mvs {
			data := RocketMQData{
				Topic: mv.GetTopic(),
				Data:  mv.GetBody(),
				Key:   mv.GetKeys(),
			}

			topicInfo.data <- data
			if err := topicInfo.consumer.Ack(context.TODO(), mv); err != nil {
				fmt.Println("ack message error: " + err.Error())
			}
		}
	}
}

func (c *RocketMQConsumer) Read(topic string) (RocketMQData, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if c.r[topic] == nil {
		return RocketMQData{}, errors.New("rocketmq consumer not exist")
	}

	msg := <-c.r[topic].data
	return msg, nil
}

func (c *RocketMQConsumer) Rput(topic string, data RocketMQData) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if c.r[topic] != nil && !c.r[topic].closed {
		c.r[topic].data <- data
	}
}

func (c *RocketMQConsumer) Close() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for _, info := range c.r {
		info.mutex.Lock()
		if !info.closed {
			info.closed = true
			info.consumer.GracefulStop()
			close(info.data)
		}
		info.mutex.Unlock()
	}
}
