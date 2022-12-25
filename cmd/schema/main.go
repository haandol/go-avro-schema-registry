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
	msgSchema := sr.Schema{
		Schema: string(msgProto),
		Type:   sr.TypeProtobuf,
	}
	isCompatible, err := rcl.CheckCompatibility(context.Background(), "message.Message", -1, msgSchema)
	if err != nil {
		logger.Errorw("failed to check compatibility", "error", err.Error())
		return
	}
	if !isCompatible {
		logger.Errorw("schema is not compatible")
		return
	}
	msgSubSchema, err := rcl.CreateSchema(context.Background(), "message.Message", msgSchema)
	if err != nil {
		logger.Errorw("failed to create schema", "error", err.Error())
		return
	}
	logger.Infow("msg schema created", "version", msgSubSchema.Version)

	cmdProto, err := os.ReadFile("idl/commandpb/car.proto")
	if err != nil {
		logger.Errorw("failed to read proto file", "error", err.Error())
		return
	}
	cmdSchema := sr.Schema{
		Schema: string(cmdProto),
		Type:   sr.TypeProtobuf,
		References: []sr.SchemaReference{
			sr.SchemaReference{
				Name:    "base.proto",
				Subject: msgSubSchema.Subject,
				Version: 1,
			},
		},
	}
	isCompatible, err = rcl.CheckCompatibility(context.Background(), "command.BookCar", -1, cmdSchema)
	if err != nil {
		logger.Errorw("failed to check compatibility", "error", err.Error())
		return
	}
	if !isCompatible {
		logger.Errorw("schema is not compatible")
		return
	}
	cmdSubSchema, err := rcl.CreateSchema(context.Background(), "command.BookCar", cmdSchema)
	if err != nil {
		logger.Errorw("failed to create schema", "error", err.Error())
		return
	}
	logger.Infow("cmd schema created", "version", cmdSubSchema.Version)

	rcl.SetCompatibilityLevel(context.Background(),
		sr.CompatBackward,
		msgSubSchema.Subject, cmdSubSchema.Subject,
	)
}
