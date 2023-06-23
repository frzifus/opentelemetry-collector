package main

import (
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/exporter/loggingexporter"
)

//export NewFactory
func NewFactory() exporter.Factory {
	return loggingexporter.NewFactory()
}
