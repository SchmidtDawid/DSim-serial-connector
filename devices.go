package main

import "fmt"

type Devices []Device

func newDevices() Devices {
	d := Devices{}
	return d
}

func getConnectedDevices() {
	openedPorts, closedPorts := findActivePorts()
	allPorts := append(openedPorts, closedPorts...)

	fmt.Println(openedPorts)
	fmt.Println(closedPorts)
	fmt.Println(allPorts)
}
