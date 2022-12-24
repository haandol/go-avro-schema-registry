package config

import (
	"log"

	"github.com/haandol/protobuf/pkg/util"
	"github.com/joho/godotenv"
)

type App struct {
	Stage string `validate:"required"`
}

type Kafka struct {
	GroupID          string   `validate:"required"`
	Topic            string   `validate:"required"`
	Seeds            []string `validate:"required"`
	SchemaRegistry   string   `validate:"required"`
	MessageExpirySec int      `validate:"required,number"`
	BatchSize        int      `validate:"required,number"`
}

type Config struct {
	App   App
	Kafka Kafka
}

// Load config.Config from environment variables for each stage
func Load() Config {
	stage := getEnv("APP_STAGE").String()
	log.Printf("Loading %s config\n", stage)

	if err := godotenv.Load(); err != nil {
		log.Panic("Error loading .env file")
	}

	cfg := Config{
		App: App{
			Stage: getEnv("APP_STAGE").String(),
		},
		Kafka: Kafka{
			Topic:            getEnv("KAFKA_TOPIC").String(),
			GroupID:          getEnv("KAFKA_GROUP_ID").String(),
			Seeds:            getEnv("KAFKA_SEEDS").Split(","),
			MessageExpirySec: getEnv("KAFKA_MESSAGE_EXPIRY_SEC").Int(),
			SchemaRegistry:   getEnv("KAFKA_SCHEMA_REGISTRY").String(),
			BatchSize:        getEnv("KAFKA_BATCH_SIZE").Int(),
		},
	}

	if err := util.ValidateStruct(cfg); err != nil {
		log.Panicf("Error validating config: %s", err.Error())
	}

	return cfg
}
