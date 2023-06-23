package main

import (
	"go.opentelemetry.io/collector/receiver"
	"go.opentelemetry.io/collector/receiver/otlpreceiver"
)

//export NewFactory
func NewFactory() receiver.Factory {
	return otlpreceiver.NewFactory()
}
