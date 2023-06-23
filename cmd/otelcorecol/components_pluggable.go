//go:build pluggable

package main

import (
	"io/ioutil"
	"log"
	"plugin"
	"strings"

	"go.opentelemetry.io/collector/connector"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/extension"
	"go.opentelemetry.io/collector/otelcol"
	"go.opentelemetry.io/collector/processor"
	"go.opentelemetry.io/collector/receiver"
)

func components() (otelcol.Factories, error) {
	var err error
	factories := otelcol.Factories{}

	r, e, p, c, x := load(".")
	factories.Extensions, err = extension.MakeFactoryMap(x...)
	if err != nil {
		return otelcol.Factories{}, err
	}

	factories.Receivers, err = receiver.MakeFactoryMap(r...)
	if err != nil {
		return otelcol.Factories{}, err
	}

	factories.Exporters, err = exporter.MakeFactoryMap(e...)
	if err != nil {
		return otelcol.Factories{}, err
	}

	factories.Processors, err = processor.MakeFactoryMap(p...)
	if err != nil {
		return otelcol.Factories{}, err
	}

	factories.Connectors, err = connector.MakeFactoryMap(c...)
	if err != nil {
		return otelcol.Factories{}, err
	}

	return factories, nil
}

func load(path string) (r []receiver.Factory, e []exporter.Factory, p []processor.Factory, c []connector.Factory, x []extension.Factory) {
	plugins, err := pluginLookup(path)
	if err != nil {
		panic(err)
	}
	for _, pluginName := range plugins {
		log.Printf("Loading plugin: %s", pluginName)
		pl, err := plugin.Open(pluginName)
		if err != nil {
			panic(err)
		}

		factorySym, err := pl.Lookup("NewFactory")
		if err != nil {
			panic(err)
		}

		if factoryFunc, ok := factorySym.(func() receiver.Factory); ok {
			r = append(r, factoryFunc())
		} else if factoryFunc, ok := factorySym.(func() exporter.Factory); ok {
			e = append(e, factoryFunc())
		} else if factoryFunc, ok := factorySym.(func() processor.Factory); ok {
			p = append(p, factoryFunc())
		} else if factoryFunc, ok := factorySym.(func() connector.Factory); ok {
			c = append(c, factoryFunc())
		} else if factoryFunc, ok := factorySym.(func() extension.Factory); ok {
			x = append(x, factoryFunc())
		} else {
			panic("ups")
		}
	}
	return
}

func pluginLookup(path string) ([]string, error) {
	var soFiles []string

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".so") {
			soFiles = append(soFiles, file.Name())
		}
	}

	return soFiles, nil
}
