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
	receiveData   bool
	configuration Config
}

func newDevice(device string, port string) (Device, error) {
	fmt.Println(device)
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

	receive := false
	if intParams[2] != 0 {
		receive = true
	}

	configFile := "config_" + strconv.Itoa(intParams[1])
	configuration := readConfigurationFromFile(configFile, "default")

	return Device{
		intParams[1],
		port,
		configFile,
		receive,
		configuration,
	}, nil
}

func (d Device) getConfiguration(device Device) Config {
	return d.configuration
}

func printConfig(device *Device) {
	fmt.Println(device.configuration)
}
