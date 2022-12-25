package main

import (
	"context"
	"fmt"

	"github.com/haandol/protobuf/pkg/config"
	"github.com/haandol/protobuf/pkg/consumer"
	"github.com/haandol/protobuf/pkg/idlpb/commandpb"
	"github.com/haandol/protobuf/pkg/util"
	"github.com/twmb/franz-go/pkg/sr"
	"google.golang.org/protobuf/proto"
)

var (
	serde sr.Serde
)

func handler(ctx context.Context, msg *consumer.Message) error {
	logger := util.GetLogger().With(
		"module", "handler",
	)

	val := &commandpb.BookCar{}
	if err := serde.Decode(msg.Value, val); err != nil {
		logger.Errorw("failed to unmarshal json", "error", err.Error())
		return err
	}

	logger.Infow("handler called", "message", fmt.Sprintf("%v", val))

	return nil
}

func main() {
	cfg := config.Load()
	logger := util.InitLogger(cfg.App.Stage).With(
		"module", "main",
	)
	logger.Infow("\n==== Config ====\n\n", "config", fmt.Sprintf("%v", cfg))

	c := consumer.NewKafkaConsumer(&cfg.Kafka)
	defer c.Close(context.Background())

	rcl, err := sr.NewClient(sr.URLs(cfg.Kafka.SchemaRegistry))
	if err != nil {
		logger.Errorw("failed to connect schema registry", "error", err.Error())
		return
	}

	cmdSchema, err := rcl.SchemaByVersion(context.Background(), "command.BookCar", 1, sr.HideDeleted)
	if err != nil {
		logger.Errorw("failed to get schema", "error", err.Error())
		return
	}

	serde.Register(
		cmdSchema.ID,
		&commandpb.BookCar{},
		sr.DecodeFn(func(b []byte, v any) error {
			return proto.Unmarshal(b, v.(*commandpb.BookCar))
		}),
	)

	if err := c.Consume(context.Background(), handler); err != nil {
		logger.Fatalw("failed to consume kafka", "error", err.Error())
	}
}
