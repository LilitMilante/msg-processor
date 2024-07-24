package broker

import (
	"context"
	"encoding/json"
	"log"

	"msg-processor/internal/entity"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	r *kafka.Reader
}

func NewConsumer(brokers []string) *Consumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     brokers,
		GroupID:     "msg-processor",
		Topic:       "unprocessed-messages",
		MaxBytes:    10e6,
		Logger:      log.Default(),
		ErrorLogger: log.Default(),
	})

	return &Consumer{
		r: r,
	}
}

func (c *Consumer) OnNewMsgEvent(ctx context.Context, fn func(context.Context, entity.Msg) error) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		m, err := c.r.ReadMessage(ctx)
		if err != nil {
			continue
		}

		var msg entity.Msg
		err = json.Unmarshal(m.Value, &msg)
		if err != nil {
			log.Println(err)
			continue
		}

		err = fn(ctx, msg)
		if err != nil {
			log.Println(err)
			continue
		}
	}
}
