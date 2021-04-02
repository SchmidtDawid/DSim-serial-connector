package main

import (
	"fmt"
	"time"
)

func main() {

	testChannel := make(chan string)
	portChannel := make(chan string)
	ports, _ := findOpenPorts()
	myDevices := checkDeviceOnPorts(ports)

	connected := false
	for !connected {
		go testConnection(testChannel)
		if <-testChannel != "connection fail" {
			connected = true
		}
	}

	readDevices(myDevices, portChannel)

	eventSC, _ := scConnect("MSFS_events")
	varReceiveSC, _ := scConnect("MSFS_vars")
	varReceiveSC.SetDelay(200 * time.Millisecond)
	planeSC, _ := scConnect("MSFS_plane")
	planeSC.SetDelay(2 * time.Second)

	go keepUpdateConfig(myDevices, planeSC)

	go startSendEvents(eventSC, myDevices, portChannel)

	cSimVar := registerVars(varReceiveSC)
	go startGettingVars(cSimVar)

	for {
		if eventSC.IsAlive() {
			fmt.Println("event simconnect is alive")
		}
		time.Sleep(time.Millisecond * 10000)
	}
}
