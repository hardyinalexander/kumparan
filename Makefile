#!/bin/bash

start:
	@make build-producer && make build-consumer && make -j 2 run-producer run-consumer

build-producer:
	@echo "Building producer..."
	@go build -o bin/producer ./cmd/producer/
	@echo "Done."

build-consumer:
	@echo "Building consumer..."
	@go build -o bin/consumer ./cmd/consumer/
	@echo "Done."

run-producer:
	@echo "Running producer binary..."
	@bin/producer

run-consumer:
	@echo "Running consumer binary..."
	@bin/consumer