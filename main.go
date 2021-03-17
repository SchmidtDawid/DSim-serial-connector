package main

import (
	"fmt"
	"github.com/micmonay/simconnect"
	"time"
)

func main() {

	portChannel := make(chan string)
	ports, _ := findOpenPorts()

	fmt.Println(ports)
	myDevices := checkDeviceOnPorts(ports)
	fmt.Println(myDevices)
	readDevices(myDevices, portChannel)

	go keepUpdateConfig(myDevices)

	//----------------------------------------
	sc, err := simconnect.NewEasySimConnect()
	if err != nil {
		panic(err)
	}
	sc.SetLoggerLevel(simconnect.LogError)

	var connected bool = false
	fmt.Println("connecting to MSFS...")
	for !connected {
		c, err := sc.Connect("Com_listener")
		if err != nil {
			time.Sleep(time.Second * 2)
			continue
		} else {
			connected = true
		}
		<-c // Wait connection confirmation
	}

	//go goEvents(sc, myDevices, portChannel)

	//go (func(){
	//	for {
	// 	Example_getSimVar()
	//	}
	//})()

	go connectToSimVars(sc)

	time.Sleep(3 * time.Second)
	event := sc.NewSimEvent(simconnect.KeyAutopilotOff)
	event.Run()

	for {
		time.Sleep(time.Second * 1000)
	}
}
