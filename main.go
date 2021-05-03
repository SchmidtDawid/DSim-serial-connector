package main

import (
	"time"
)

var eventSC, _ = scConnect("MSFS_events")
var scGlobal = newScGlobalData()

func main() {

	go scGlobal.update()

	//testChannel := make(chan string)
	//cComSend := make(chan []string)
	//ports, _ := findActivePorts()
	//myDevices := checkDeviceOnPorts(ports)

	devices := newDevices()
	devices.getConnectedDevices()
	devices.connectToDevices()
	devices.listenDevices()
	devices.startLifecycles()
	go devices.monitor()

	for _, device := range devices {
		go func(d *Device) {
			for msg := range d.cRec {
				//fmt.Println(msg.msg)

				incomingEvents, err := collectEvents(msg.msg)
				if err != nil {
				}
				deviceEvents = append(deviceEvents, incomingEvents...)
				for len(deviceEvents) > 0 {
					event := deviceEvents[0]
					deviceEvents = deviceEvents[1:]
					executeEvents(event, d)
				}

			}
		}(device)
	}

	for {
		time.Sleep(time.Millisecond * 30000)
	}

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
	//varReceiveSC, _ := scConnect("MSFS_vars")
	//varReceiveSC.SetDelay(50 * time.Millisecond)

	//
	//go keepUpdateConfig(myDevices, globalSc)
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
