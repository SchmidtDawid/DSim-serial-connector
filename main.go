package main

import (
	"fmt"
	"time"
)

func main() {

	testChannel := make(chan string)
	cComReceive := make(chan string)
	cComSend := make(chan []string)
	ports, _ := findOpenPorts()
	myDevices := checkDeviceOnPorts(ports)

	connected := false
	for !connected {
		go testConnection(testChannel)
		if <-testChannel != "connection fail" {
			connected = true
		}
	}

	connectDevices(myDevices, cComReceive, cComSend)

	eventSC, _ := scConnect("MSFS_events")
	varReceiveSC, _ := scConnect("MSFS_vars")
	varReceiveSC.SetDelay(50 * time.Millisecond)
	planeSC, _ := scConnect("MSFS_plane")
	planeSC.SetDelay(2 * time.Second)

	go keepUpdateConfig(myDevices, planeSC)

	go startSendEvents(eventSC, myDevices, cComReceive)

	cSimVar := registerVars(varReceiveSC)
	go startGettingVars(cSimVar, cComSend)

	for {
		if eventSC.IsAlive() {
			fmt.Println("event simconnect is alive")
		}
		time.Sleep(time.Millisecond * 10000)
	}
}
