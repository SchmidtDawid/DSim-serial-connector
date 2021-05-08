package main

import (
	"fmt"
	"github.com/micmonay/simconnect"
	"strings"
	"time"
)

type ScGlobalData struct {
	planeName string
	sc        *simconnect.EasySimConnect
}

func scConnect(appName string) (*simconnect.EasySimConnect, error) {
	sc, err := simconnect.NewEasySimConnect()
	if err != nil {
		fmt.Println("Cannot connect to ", appName+"...")
		panic(err)
	}
	sc.SetLoggerLevel(simconnect.LogError)

	var connected bool = false
	fmt.Println("connecting to", appName+"...")
	for !connected {
		c, err := sc.Connect(appName)
		if err != nil {
			time.Sleep(time.Second * 2)
			continue
		} else {
			connected = true
		}
		<-c
	}

	return sc, nil
}

func newScGlobalData() *ScGlobalData {
	return &ScGlobalData{
		sc:        globalSc,
		planeName: "",
	}
}

func (d *ScGlobalData) update() {
	cSimVar, err := d.sc.ConnectToSimVar(
		simconnect.SimVarTitle(),
	)
	if err != nil {
		fmt.Println("Can not register Vars")
	}

	var result []simconnect.SimVar

	for range time.Tick(time.Second * 2) {

		result = <-cSimVar
		for _, simVar := range result {
			if strings.Contains(string(simVar.Unit), "String") {
				d.planeName = simVar.GetString()
			}
		}
	}
}
