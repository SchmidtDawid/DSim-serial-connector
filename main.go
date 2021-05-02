package main

import (
	"fmt"
	"time"
)

func main() {

	//testChannel := make(chan string)
	//cComSend := make(chan []string)
	//ports, _ := findActivePorts()
	//myDevices := checkDeviceOnPorts(ports)

	devices := newDevices()
	devices.getConnectedDevices()
	devices.connectToDevices()
	devices.listenDevices()
	for _, device := range devices {
		fmt.Println(*device)
	}

	for msg := range cDevicesReceive {
		fmt.Println(msg.device.port)
		//incomingEvents, err := collectEvents(msg.msg)
		//if err != nil {
		//}
		//deviceEvents = append(deviceEvents, incomingEvents...)
		//for len(deviceEvents) > 0 {
		//	event := deviceEvents[0]
		//	deviceEvents = deviceEvents[1:]
		//	fmt.Println(event)
		//	executeEvents(sc, event, myDevices)
		//}
	}

	time.Now()
	time.Sleep(time.Millisecond * 30000)

	//connected := false
	//for !connected {
	//	go testConnection(testChannel)
	//	if <-testChannel != "serial fail" {
	//		connected = true
	//	}
	//}
	//
	//connectDevices(myDevices, cDevicesReceive, cComSend)
	//
	//eventSC, _ := scConnect("MSFS_events")
	//varReceiveSC, _ := scConnect("MSFS_vars")
	//varReceiveSC.SetDelay(50 * time.Millisecond)
	//planeSC, _ := scConnect("MSFS_plane")
	//planeSC.SetDelay(2 * time.Second)
	//
	//go keepUpdateConfig(myDevices, planeSC)
	//
	//go startSendEvents(eventSC, myDevices, cDevicesReceive)
	//
	//cSimVar := registerVars(varReceiveSC)
	//go startGettingVars(cSimVar, cComSend)
	//
	//for {
	//	if eventSC.IsAlive() {
	//		fmt.Println("event simconnect is alive")
	//	}
	//	time.Sleep(time.Millisecond * 10000)
	//}
}
