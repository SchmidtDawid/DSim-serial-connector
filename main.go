package main

import (
	"fmt"
	"time"
)

func main() {

	portChannel := make(chan string)
	ports, _ := findOpenPorts()
	myDevices := checkDeviceOnPorts(ports)

	readDevices(myDevices, portChannel)
	go keepUpdateConfig(myDevices)

	eventSC, _ := scConnect("MSFS_events")
	varReveiveSC, _ := scConnect("MSFS_vars")
	varReveiveSC.SetDelay(50 * time.Millisecond)

	go startSendEvents(eventSC, myDevices, portChannel)
	go startGettingVars(varReveiveSC)

	for {
		if eventSC.IsAlive() {
			fmt.Println("event simconnect is alive")
		}
		time.Sleep(time.Second * 10)
	}
}
