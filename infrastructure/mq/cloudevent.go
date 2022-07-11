package mq

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Shopify/sarama"
	kafka_sarama "github.com/cloudevents/sdk-go/protocol/kafka_sarama/v2"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/client"
	"github.com/opensourceways/xihe-server/config"
	log "github.com/sirupsen/logrus"
)

const (
	SourceUrl = "github.com/opensourceways/xihe-server"
)

type MQEventName string

var (
	ProjectLikeCountIncreaseEvent MQEventName = "XiheServerProjectLikeCountIncrease"
)

var (
	saramaConfig *sarama.Config
	brokers      []string
	MQClientMap  map[MQEventName]cloudevents.Client
)

type (
	Closeable interface {
		Close()
	}

	Notifier interface {
		Closeable
		PushEvent(eventType, subject string, data map[string]interface{})
	}
)

func InitMQ(cfg *config.Config) {
	saramaConfig = sarama.NewConfig()
	saramaConfig.Version = sarama.V2_8_1_0
	saramaConfig.Consumer.Offsets.Initial = sarama.OffsetNewest
	brokers = strings.Split(cfg.Kafka.KafkaBrokers, ",")
	MQClientMap = make(map[MQEventName]cloudevents.Client)
	initEventNotifier(cfg, ProjectLikeCountIncreaseEvent)
}

func initEventNotifier(cfg *config.Config, eventType MQEventName) error {
	sender, err := kafka_sarama.NewSender(brokers, saramaConfig, string(eventType))
	if err != nil {
		return errors.New(fmt.Sprintf("failed to create protocol: %v", err.Error()))
	}
	cloudEventClient, err := cloudevents.NewClient(sender, cloudevents.WithTimeNow(), cloudevents.WithUUIDs())
	if err != nil {
		return errors.New(fmt.Sprintf("failed to create cloudevents client, %v", err))
	}
	MQClientMap[ProjectLikeCountIncreaseEvent] = cloudEventClient
	//-------------do other event ...
	return nil
}
func PushEvent(eventType MQEventName, data map[string]interface{}) error {
	clientItem := MQClientMap[eventType]
	if clientItem == nil {
		return fmt.Errorf("[PushEvent] has no [%v] event ", eventType)
	}
	e := cloudevents.NewEvent()
	e.SetSpecVersion(cloudevents.VersionV1)
	e.SetType(string(eventType))
	e.SetSource(SourceUrl)
	e.SetData(cloudevents.ApplicationJSON, data)
	go func() {
		err := clientItem.Send(kafka_sarama.WithMessageKey(context.Background(), sarama.StringEncoder(e.ID())), e)
		if err != nil {
			log.Error(fmt.Sprintf("[PushEvent] failed to send message ,error: %v", err))
		}
		log.Info(fmt.Sprintf("[PushEvent] message send with event type %s,  data %v", eventType, data))
	}()
	return nil
}

// go StartEventLinsten(ProjectLikeCountIncreaseEvent, handleDownloadStatusEvent)
func StartEventLinsten(topicName MQEventName, GroupID string, handleFunc func(ctx context.Context, event cloudevents.Event)) {
	receiver, err := kafka_sarama.NewConsumer(brokers, saramaConfig, GroupID, string(topicName))
	if err != nil {
		log.Printf("failed to create protocol: %s", err.Error())

	}
	defer receiver.Close(context.Background())
	clientItem, err := cloudevents.NewClient(receiver, client.WithPollGoroutines(1))
	if err != nil {
		log.Printf("failed to create client, %v", err.Error())
	}
	err = clientItem.StartReceiver(context.Background(), handleFunc)
	if err != nil {
		log.Printf("failed to start receiver(%s) error: %s", topicName, err.Error())

	}
	log.Printf(" TOPIC :(%s) 停止监听。\n", topicName)
}
