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
	go func(d *Devices) {
		for {
			time.Sleep(time.Second)
			for _, device := range *d {

				if device.connection.connected == false {
					device.connect()
				}

				fmt.Printf("port : %v, ID: %v, familiar: %v, isReceiving:  %v, connected: %v\n",
					device.port, device.id, device.isFamiliar, device.isReceivingData, device.connection.connected)
			}
			fmt.Println("----------")

		}
	}(d)

	go func(d *Devices) {
		checkNewTimer := time.NewTicker(time.Second * 3)

		for {
			select {
			case _ = <-checkNewTimer.C:
				openedPorts, _ := findActivePorts()
				for _, port := range openedPorts {
					for _, device := range *d {
						if device.port == port {
							return
						}
					}
					*d = append(*d, newEmptyDevice(port))
				}
			}
		}
	}(d)
}
