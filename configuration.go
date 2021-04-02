package main

import (
	"encoding/json"
	"fmt"
	"github.com/micmonay/simconnect"
	"io/ioutil"
	"os"
	"strings"
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

func readConfigurationFromFile(file string, planeName string) Config {

	optionalFileName := file + "_" + strings.ToLower(strings.ReplaceAll(planeName, " ", "_"))

	var jsonFile *os.File
	var err error

	jsonFile, err = os.Open(optionalFileName + ".json")
	if err != nil {
		jsonFile, err = os.Open(file + "_default" + ".json")
		if err != nil {
			fmt.Println(err)
		}
	}

	//fmt.Println("Successfully Opened", file)
	defer jsonFile.Close()

	byteValue, bErr := ioutil.ReadAll(jsonFile)
	if bErr != nil {
		fmt.Println(err)
	}

	var config Config
	err = json.Unmarshal(byteValue, &config)
	if err != nil {
		fmt.Println(err)
	}

	return config
}

func (d *Device) updateConfiguration(planeName string) {
	c := readConfigurationFromFile(d.configFile, planeName)
	d.configuration = c
}

func keepUpdateConfig(devices []*Device, sc *simconnect.EasySimConnect) {
	cSimVar, err := sc.ConnectToSimVar(
		simconnect.SimVarTitle(),
	)
	if err != nil {
		fmt.Println("Can not register Vars")
	}

	var result []simconnect.SimVar
	var planeName string

	for range time.Tick(time.Second * 2) {

		result = <-cSimVar
		for _, simVar := range result {
			if strings.Contains(string(simVar.Unit), "String") {
				planeName = simVar.GetString()
			}
		}

		for _, device := range devices {
			device.updateConfiguration(planeName)
		}
	}
}
