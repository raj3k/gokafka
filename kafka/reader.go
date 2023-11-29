package kafka

import (
	"context"
	"fmt"
	"log"

	"github.com/pkg/errors"
	kafkago "github.com/segmentio/kafka-go"
)

type KafkaReader struct {
	Reader *kafkago.Reader
}

func NewKafkaReader() *KafkaReader {
	reader := kafkago.NewReader(kafkago.ReaderConfig{
		Brokers: []string{"localhost:9093"},
		Topic:   "random_numbers",
		GroupID: "group",
	})

	return &KafkaReader{
		Reader: reader,
	}
}

// func (k *KafkaReader) FetchMessages(ctx context.Context, messages chan<- kafkago.Message) error {
// 	for {
// 		message, err := k.Reader.FetchMessage(ctx)
// 		if err != nil {
// 			return err
// 		}
// 		fmt.Println(string(message.Value))

// 		select {
// 		case <-ctx.Done():
// 			return ctx.Err()
// 		case messages <- message:
// 			log.Printf("message fetched and send to a channel: %v \n", string(message.Value))
// 		}
// 	}
// }

func (k *KafkaReader) FetchMessages(ctx context.Context) error {
	for {
		message, err := k.Reader.FetchMessage(ctx)
		if err != nil {
			return err
		}
		fmt.Println(string(message.Value))
	}
}

func (k *KafkaReader) CommitMessage(ctx context.Context, messageCommitChan <-chan kafkago.Message) error {
	for {
		select {
		case <-ctx.Done():
		case msg := <-messageCommitChan:
			err := k.Reader.CommitMessages(ctx, msg)
			if err != nil {
				return errors.Wrap(err, "Reader.CommitMessages")
			}
			log.Printf("committed an message: %v \n", string(msg.Value))
		}
	}
}
