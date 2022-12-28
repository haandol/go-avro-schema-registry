# Go Protobuf Schema Registry

this repository is for the Go example of the Confluent Schema Registry using Protobuf.

# Prerequisites

- Docker
- Go 1.19
- [Taskfile](https://taskfile.dev/#/installation)
- [Protobuf](https://developers.google.com/protocol-buffers/docs/gotutorial)

# Getting Started

## Install depencencies

```bash
$ go mod tidy && go mod vendor
```

## Run infrastructure

```bash
$ docker-compose --profile backend up -d
```

## Copy .ENV

```bash
$ cp env/local.env .env
```

## Build IDL

build Protobuf IDL interface 

```bash
$ task idl
task: [idl] mkdir -p pkg/idlpb
task: [idl] protoc -I ./idl --go_opt=paths=source_relative --go_out=pkg/idlpb ./idl/**/*.proto
```

## Create schema on registry

register schema

```bash
$ task schema
```

find schema on registry
  
```bash
$ open http://localhost:8081/schemas
```

## Run the application

```bash
$ docker-compose up dev
```

find messages on kafka-ui

```bash
$ open http://localhost:8080/ui/clusters/local/topics
```