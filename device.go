package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Device struct {
	id            int
	port          string
	configFile    string
	configuration Config
}

func newDevice(device string, port string) (Device, error) {
	device = strings.Split(device, ";")[0]
	params := strings.Split(device, ",")
	if len(params) != 5 {
		return Device{},
			errors.New("wrong event format")
	}

	var intParams []int
	for _, param := range params {
		intParam, err := strconv.Atoi(param)
		if err != nil {
			return Device{}, err
		}
		intParams = append(intParams, intParam)
	}

	configFile := "config_" + strconv.Itoa(intParams[1]) + ".json"
	configuration := readConfiguration(configFile)

	return Device{
		intParams[1],
		port,
		configFile,
		configuration,
	}, nil
}

func (d Device) getConfiguration(device Device) Config {
	return d.configuration
}

func printConfig(device *Device) {
	fmt.Println(device.configuration)
}
