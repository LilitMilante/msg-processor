package broker

import (
	"context"
	"encoding/json"
	"log"

	"msg-processor/internal/entity"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	w *kafka.Writer
}

func NewProducer(addr []string) *Producer {
	w := &kafka.Writer{
		Addr:                   kafka.TCP(addr...),
		Topic:                  "unprocessed-messages",
		Balancer:               &kafka.LeastBytes{},
		RequiredAcks:           1,
		Logger:                 log.Default(),
		ErrorLogger:            log.Default(),
		AllowAutoTopicCreation: true,
	}
	return &Producer{
		w: w,
	}
}

func (p *Producer) SendMsg(ctx context.Context, msgs ...entity.Msg) error {
	msgsToSend := make([]kafka.Message, 0, len(msgs))

	for _, msg := range msgs {
		bJSON, err := json.Marshal(msg)
		if err != nil {
			return err
		}

		msgsToSend = append(msgsToSend, kafka.Message{Key: []byte(msg.ID.String()), Value: bJSON})
	}

	err := p.w.WriteMessages(ctx, msgsToSend...)
	if err != nil {
		return err
	}

	return nil
}
