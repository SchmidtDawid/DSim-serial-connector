package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
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

func readConfigurationFromFile(device *Device) Config {

	fileBase := "config_" + strconv.Itoa(device.id)
	fmt.Println("PLANE", scGlobal.planeName)
	optionalFileName := fileBase + "_" + strings.ToLower(strings.ReplaceAll(scGlobal.planeName, " ", "_"))

	var jsonFile *os.File
	var err error

	jsonFile, err = os.Open(optionalFileName + ".json")
	if err != nil {
		jsonFile, err = os.Open(fileBase + "_default" + ".json")
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
		fmt.Println(fileBase, err)
	}

	return config
}

//func keepUpdateConfig(devices []*Device, sc *simconnect.EasySimConnect) {
//	cSimVar, err := sc.ConnectToSimVar(
//		simconnect.SimVarTitle(),
//	)
//	if err != nil {
//		fmt.Println("Can not register Vars")
//	}
//
//	var result []simconnect.SimVar
//	var planeName string
//
//	for range time.Tick(time.Second * 2) {
//
//		result = <-cSimVar
//		for _, simVar := range result {
//			if strings.Contains(string(simVar.Unit), "String") {
//				planeName = simVar.GetString()
//			}
//		}
//
//		for _, device := range devices {
//			device.updateConfiguration(planeName)
//		}
//	}
//}
