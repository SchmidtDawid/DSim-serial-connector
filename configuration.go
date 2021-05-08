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
	Elements Elements     `json:"config"`
	Vars     []VarRequest `json:"vars"`
}

type Elements struct {
	Buttons  [][]Action `json:"buttons"`
	Switches [][]Action `json:"switches"`
	Encoders [][]Action `json:"encoders"`
}

type Action struct {
	Action   int        `json:"action"`
	SimEvent []SimEvent `json:"simEvent"`
}

type SimEvent struct {
	Event string `json:"event"`
	Value int    `json:"value"`
}

type VarRequest struct {
	Name     string `json:"name"`
	Unit     string `json:"unit"`
	Settable bool   `json:"settable"`
	Value    string
}

func readConfigurationFromFile(device *Device) Config {

	fileBase := "configs/config_" + strconv.Itoa(device.id)
	//fmt.Println("PLANE", scGlobal.planeName)
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
