package main

type Devices []*Device

func newDevices() Devices {
	d := Devices{}
	return d
}

func (d *Devices) getConnectedDevices() {
	openedPorts, closedPorts := findActivePorts()
	allPorts := append(openedPorts, closedPorts...)

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
		device.ask()
	}
}

func (d *Devices) listenDevices() {
	for _, device := range *d {
		device.listen()
	}
}
