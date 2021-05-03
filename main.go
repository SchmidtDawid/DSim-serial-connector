package main

import (
	"fmt"
	"github.com/micmonay/simconnect"
	"time"
)

var eventSC *simconnect.EasySimConnect
var globalSc *simconnect.EasySimConnect
var scGlobal *ScGlobalData

var cComSend = make(chan []string)

func main() {
	var testChannel = make(chan string)

	devices := newDevices()
	devices.getConnectedDevices()
	devices.connectToDevices()

	connected := false
	for !connected {
		go testConnection(testChannel)
		if <-testChannel != "serial fail" {
			connected = true
		}
	}

	eventSC, _ = scConnect("MSFS_events")
	globalSc, _ = scConnect("MSFS_plane")
	scGlobal = newScGlobalData()

	go scGlobal.update()

	devices.listenDevices()
	devices.startLifecycles()
	go devices.monitor()

	varReceiveSC, _ := scConnect("MSFS_vars")
	varReceiveSC.SetDelay(50 * time.Millisecond)

	cSimVar := registerVars(varReceiveSC)
	go startGettingVars(cSimVar, cComSend)

	for {
		if eventSC.IsAlive() {
			fmt.Println("event simconnect is alive")
		}
		time.Sleep(time.Millisecond * 10000)
	}
}
