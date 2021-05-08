package main

import (
	"github.com/micmonay/simconnect"
	"time"
)

var eventSC *simconnect.EasySimConnect
var globalSc *simconnect.EasySimConnect
var scGlobal *ScGlobalData

func main() {
	var testChannel = make(chan string)

	devices := newDevices()
	devices.getConnectedDevices()

	connected := false
	for !connected {
		go testConnection(testChannel)
		if <-testChannel != "serial fail" {
			connected = true
		}
	}

	eventSC, _ = scConnect("MSFS_events")
	globalSc, _ = scConnect("MSFS_global")
	globalSc.SetDelay(2 * time.Second)
	scGlobal = newScGlobalData()

	go scGlobal.update()

	devices.monitor()

	for {
		//if eventSC.IsAlive() {
		//	fmt.Println("event simconnect is alive")
		//}
		time.Sleep(time.Millisecond * 10000)
	}
}
