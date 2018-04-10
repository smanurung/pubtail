package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"reflect"
	"strings"
	"syscall"

	"cloud.google.com/go/pubsub"
	log "github.com/Sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type message struct {
	fields    log.Fields
	body      []byte
	formatted []byte
}

type formatter interface {
	format([]byte) ([]byte, error)
}

var (
	fmtmap = map[string]formatter{
		"avro": new(avroFormatter),
		"text": new(textFormatter),
	}
)

func main() {
	topics := flag.String("topics", "playSonny1", "topics to subscribe in comma separated form")
	projectID := flag.String("project", "projectid", "projectID on GCP")
	format := flag.String("format", "text", "global format for printing message body")
	flag.Parse()

	ctx := context.Background()

	log.Infof("connecting to projectid %s, using format %s", *projectID, *format)

	cli, err := pubsub.NewClient(ctx, *projectID)
	if err != nil {
		log.Fatalf("failed to create new pubsub client: %v", err)
	}

	topicSlice := strings.Split(*topics, ",")
	msgChan := make(chan message)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		s := <-c
		log.Infof("receiving signal: %s - will exit gracefully", s)

		for _, topicName := range topicSlice {
			subName := fmt.Sprintf("pubtail___%s___", topicName)
			sub := cli.Subscription(subName)

			if err := sub.Delete(ctx); err != nil {
				log.Errorln(err)
				return
			}
			log.Infof("Subscription %s has been deleted.", subName)
		}
		log.Info("cleaned")
	}()

	for _, topic := range topicSlice {
		go listenToTopic(ctx, cli, msgChan, topic, *format)
	}

	for msg := range msgChan {
		log.WithFields(msg.fields).Infof("%s", msg.formatted)
	}

	log.Println("pubtail finished")
	os.Exit(0)
}

func listenToTopic(ctx context.Context, cli *pubsub.Client, ch chan<- message, topicName, format string) error {
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

		formatted, err := fmtmap[format].format(m.Data)
		if err != nil {
			log.Errorf("failed to use format %s for message %s", format, m.Data)
			return
		}

		ch <- message{
			fields:    f,
			body:      m.Data,
			formatted: formatted,
		}

		m.Ack()
	})
	if err != nil {
		log.Errorf("error on receive(): %v", err)
	}
	return err
}
