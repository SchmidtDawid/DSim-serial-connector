package main

import (
	"fmt"
	"github.com/tarm/serial"
	"log"
	"time"
)

type Device struct {
	id              int
	port            string
	isFamiliar      bool
	configFile      string
	isReceivingData bool
	serial          *serial.Port
	connection      DeviceConnection
	configuration   Config
	cRec            chan DeviceMsg
	cSnd            chan DeviceMsg
}

type DeviceConnection struct {
	connected          bool
	begin              time.Time
	pokeInterval       time.Duration
	connectionFailures int
	lastSeen           time.Time
}

type DeviceMsg struct {
	device *Device
	msg    string
}

func newEmptyDevice(port string) *Device {
	return &Device{
		port: port,
		cRec: make(chan DeviceMsg),
		cSnd: make(chan DeviceMsg),
	}
}

func (d *Device) connect() {
	config := &serial.Config{Name: d.port, Baud: 57600}
	s, err := serial.OpenPort(config)
	if err != nil {
		d.connection.connected = false
		return
	}
	d.serial = s
	d.connection.connected = true
	d.connection.pokeInterval = time.Second * 5
	d.connection.begin = time.Now()
}

func (d *Device) listen() {
	for {
		buf := make([]byte, 100)
		n, err := d.serial.Read(buf)
		if err != nil {
			log.Fatal(err)
		}
		d.cRec <- DeviceMsg{
			d,
			string(buf[:n]),
		}
	}
}

func (d *Device) sanitize() {
	if !d.connection.connected && time.Now().Before(d.connection.begin.Add(time.Second*4)) {
		return
	}
	if !d.connection.lastSeen.Add(d.connection.pokeInterval * 3).Before(time.Now()) {
		fmt.Println(d.id, " device lost")
	}
	d.serial.Write([]byte("?;"))
}

func (d Device) sanitizeCheck(presentation DevicePresentationEvent) {
	if d.id == presentation.deviceID {
		fmt.Println("sanitaze OK!")
	} else {
		//TODO reset device
	}
}

func (d *Device) updateConfiguration() {
	if !d.isFamiliar {
		return
	}
	c := readConfigurationFromFile(d)
	d.configuration = c
}

func (d *Device) getConfiguration() Config {
	return d.configuration
}

func (d *Device) lifecycle() {

	pokeTimer := time.NewTicker(d.connection.pokeInterval)
	configTimer := time.NewTicker(time.Second * 2)

	for {
		select {
		case _ = <-pokeTimer.C:
			d.sanitize()
		case _ = <-configTimer.C:
			d.updateConfiguration()
		}
	}
}

func printConfig(device *Device) {
	fmt.Println(device.configuration)
}
