package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

type Config struct {
	Elements Elements `json:"config"`
}

type Elements struct {
	Buttons  [][]Action `json:"buttons"`
	Switches [][]Action `json:"switches"`
	Encoders [][]Action `json:"encoders"`
}

type Action struct {
	Action   int    `json:"action"`
	SimEvent string `json:"simEvent"`
	Value    int    `json:"value"`
}

func readConfigurationFromFile(file string) Config {
	jsonFile, err := os.Open(file)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println("Successfully Opened", file)
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println(err)
	}

	var config Config
	err = json.Unmarshal(byteValue, &config)
	if err != nil {
		fmt.Println(err)
	}

	return config
}

func (d *Device) updateConfiguration() {
	c := readConfigurationFromFile(d.configFile)
	d.configuration = c
}

func keepUpdateConfig(devices []*Device) {
	for range time.Tick(time.Second * 2) {
		//for _, device := range devices {
		devices[0].updateConfiguration()
		//}
	}
}
