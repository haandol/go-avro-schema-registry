package main

import (
	"context"
	"fmt"
	"os"

	"github.com/haandol/protobuf/pkg/config"
	"github.com/haandol/protobuf/pkg/util"
	"github.com/twmb/franz-go/pkg/sr"
)

func main() {
	cfg := config.Load()
	logger := util.InitLogger(cfg.App.Stage).With(
		"module", "main",
	)
	logger.Infow("\n==== Config ====\n\n", "config", fmt.Sprintf("%v", cfg))

	rcl, err := sr.NewClient(sr.URLs(cfg.Kafka.SchemaRegistry))
	if err != nil {
		logger.Errorw("failed to connect schema registry", "error", err.Error())
		return
	}

	msgProto, err := os.ReadFile("idl/base.proto")
	if err != nil {
		logger.Errorw("failed to read proto file", "error", err.Error())
		return
	}
	msgSchema, err := rcl.CreateSchema(context.Background(), "message.Message", sr.Schema{
		Schema: string(msgProto),
		Type:   sr.TypeProtobuf,
	})
	if err != nil {
		logger.Errorw("failed to create schema", "error", err.Error())
		return
	}
	logger.Infow("cmd schema created", "version", msgSchema.Version)

	cmdProto, err := os.ReadFile("idl/commandpb/car.proto")
	if err != nil {
		logger.Errorw("failed to read proto file", "error", err.Error())
		return
	}
	cmdSchema, err := rcl.CreateSchema(context.Background(), "command.BookCar", sr.Schema{
		Schema: string(cmdProto),
		Type:   sr.TypeProtobuf,
		References: []sr.SchemaReference{
			sr.SchemaReference{
				Name:    "base.proto",
				Subject: msgSchema.Subject,
				Version: 1,
			},
		},
	})
	if err != nil {
		logger.Errorw("failed to create schema", "error", err.Error())
		return
	}
	logger.Infow("cmd schema created", "version", cmdSchema.Version)

	rcl.SetCompatibilityLevel(context.Background(), sr.CompatBackward, msgSchema.Subject, cmdSchema.Subject)
}
