#!/bin/bash

go build -v -o loggingexporter.so -buildmode=plugin ../../exporter/loggingexporter/plugin/main.go
go build -v -o otlpreceiver.so -buildmode=plugin ../../receiver/otlpreceiver/plugin/main.go
go build -tags pluggable
./otelcorecol --config=config.yaml
