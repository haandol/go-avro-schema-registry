package main

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/google/uuid"
	"github.com/haandol/protobuf/pkg/config"
	"github.com/haandol/protobuf/pkg/idlpb"
	"github.com/haandol/protobuf/pkg/idlpb/commandpb"
	"github.com/haandol/protobuf/pkg/producer"
	"github.com/haandol/protobuf/pkg/util"
	"github.com/twmb/franz-go/pkg/sr"
	"google.golang.org/protobuf/proto"
)

func main() {
	cfg := config.Load()
	logger := util.InitLogger(cfg.App.Stage).With(
		"module", "main",
	)
	logger.Infow("\n==== Config ====\n\n", "config", fmt.Sprintf("%v", cfg))

	producer, err := producer.Connect(&cfg.Kafka)
	if err != nil {
		logger.Errorw("failed to connect kafka", "error", err.Error())
		return
	}

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

	var serde sr.Serde
	serde.Register(
		cmdSchema.ID,
		&commandpb.BookCar{},
		sr.EncodeFn(func(v any) ([]byte, error) {
			return proto.Marshal(v.(*commandpb.BookCar))
		}),
	)

	for {
		cmd := &commandpb.BookCar{
			Message: &idlpb.Message{
				Name:      reflect.ValueOf(commandpb.BookCar{}).Type().Name(),
				Version:   "1.0.0",
				Id:        uuid.NewString(),
				CreatedAt: time.Now().Format(time.RFC3339),
			},
			Body: &commandpb.BookCarBody{
				TripId: 1,
				CarId:  1,
			},
		}

		if err := producer.Produce(context.Background(),
			cfg.Kafka.Topic, "", serde.MustEncode(cmd),
		); err != nil {
			logger.Errorw("failed to produce kafka", "error", err.Error())
			return
		}

		time.Sleep(1 * time.Second)
	}
}
