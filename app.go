package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"reflect"
	"strings"

	"cloud.google.com/go/pubsub"
	log "github.com/Sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type message struct {
	fields log.Fields
	body   []byte
}

func main() {
	topics := flag.String("topics", "playSonny1", "topics to subscribe in comma separated form")
	projectID := flag.String("project", "tokopedia-970", "projectID on GCP")
	flag.Parse()

	ctx := context.Background()

	cli, err := pubsub.NewClient(ctx, *projectID)
	if err != nil {
		log.Fatalf("failed to create new pubsub client: %v", err)
	}

	topicSlice := strings.Split(*topics, ",")
	msgChan := make(chan message)

	for _, topic := range topicSlice {
		go listenToTopic(ctx, cli, msgChan, topic)
	}

	for msg := range msgChan {
		log.WithFields(msg.fields).Infof("%s", msg.body)
	}

	log.Println("pubtail finished")
	os.Exit(0)
}

func listenToTopic(ctx context.Context, cli *pubsub.Client, ch chan<- message, topicName string) error {
	subname := fmt.Sprintf("pubtail___%s___", topicName)

	topic := cli.Topic(topicName)

	if ok, err := topic.Exists(ctx); !ok || err != nil {
		log.Fatalf("topic %s doesn't exist", topicName)
	}

	sub, err := cli.CreateSubscription(ctx, subname, pubsub.SubscriptionConfig{
		Topic: topic,
	})
	if status.Convert(err).Code() == codes.AlreadyExists {
		log.Infoln("subscription already exists")
		sub = cli.Subscription(subname)
	} else if err != nil {
		log.Fatalf("failed to create subscription %v, %v", err, reflect.TypeOf(err))
	}

	log.Infof("listening to topic %s...", topicName)

	// This will block
	err = sub.Receive(ctx, func(ctx context.Context, m *pubsub.Message) {
		f := log.Fields{
			"topic": topicName,
		}
		for k, v := range m.Attributes {
			f[k] = v
		}

		ch <- message{
			fields: f,
			body:   m.Data,
		}

		m.Ack()
	})
	if err != nil {
		log.Printf("error on receive(): %v", err)
	}
	return err
}
