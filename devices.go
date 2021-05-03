package main

import (
	"fmt"
	"time"
)

type Devices []*Device

func newDevices() Devices {
	d := Devices{}
	return d
}

func (d *Devices) getConnectedDevices() {
	openedPorts, busyPorts := findActivePorts()
	allPorts := append(openedPorts, busyPorts...)

	for _, port := range allPorts {
		*d = append(*d, newEmptyDevice(port))
	}
}

func (d *Devices) connectToDevices() {
	for _, device := range *d {
		device.connect()
	}
}

func (d *Devices) askDevices() {
	for _, device := range *d {
		device.sanitize()
	}
}

func (d *Devices) listenDevices() {
	for _, device := range *d {
		go device.listen()
		time.Sleep(time.Millisecond * 100)
	}
}

func (d *Devices) startLifecycles() {
	for _, device := range *d {
		go device.lifecycle()
		time.Sleep(time.Millisecond * 100)
	}
}

func (d *Devices) monitor() {
	for {
		time.Sleep(time.Second)
		for _, device := range *d {
			fmt.Printf("ID: %v, FAMILIAR: %v\n", device.id, device.isFamiliar)
			//fmt.Printf("%+v\n", device)
		}
		fmt.Println("----------")
	}
}
