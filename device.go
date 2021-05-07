package main

import (
	"context"
	"fmt"
	"github.com/tarm/serial"
	"time"
)

type Device struct {
	id              int
	port            string
	isFamiliar      bool
	configFile      string
	isReceivingData bool
	hasLifecycle    bool
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
	ctx                context.Context
	cancel             context.CancelFunc
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
	if d.connection.connectionFailures > 5 {
		return
	}
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
	d.connection.ctx, d.connection.cancel = context.WithCancel(context.Background())
	go d.listen()
	go d.lifecycle()
	go func() {
		time.Sleep(time.Second * 3)
		d.serial.Write([]byte("?;"))
	}()
}

func (d *Device) disconnect() {
	d.connection.cancel()
	d.serial.Close()
	d.connection.connected = false
}

func (d *Device) listen() {
	if !d.connection.connected {
		return
	}
	go func(d *Device) {
		if !d.connection.connected {
			return
		}
		var deviceEvents []deviceEvent
		for msg := range d.cRec {
			select {
			case <-d.connection.ctx.Done():
				return
			default:
				//fmt.Println(d.id)

				incomingEvents, err := collectEvents(msg.msg)
				if err != nil {
					fmt.Println(err)
				}
				deviceEvents = append(deviceEvents, incomingEvents...)
				for len(deviceEvents) > 0 {
					event := deviceEvents[0]
					deviceEvents = deviceEvents[1:]
					executeEvents(event, d)
				}
			}
		}
	}(d)

	for {
		select {
		case <-d.connection.ctx.Done():
			return
		default:
			buf := make([]byte, 100)
			n, err := d.serial.Read(buf)
			if err != nil {
				fmt.Println(err)
				d.disconnect()
			}
			d.cRec <- DeviceMsg{
				d,
				string(buf[:n]),
			}
		}
	}
}

func (d *Device) writeTo() {
	if d.isReceivingData {
		go writeCom(d.serial, cComSend)
	}
}

func (d *Device) sanitize() {
	if !d.connection.connected || time.Now().Before(d.connection.begin.Add(time.Second*3)) {
		return
	}
	d.serial.Write([]byte("?;"))
	if time.Now().After(d.connection.lastSeen.Add(time.Second * 11)) {
		fmt.Println("error? ", d.id)
		d.connection.connectionFailures++
		if d.connection.connectionFailures > 5 {
			d.disconnect()
		}
	}
}

func (d Device) sanitizeCheck(presentation DevicePresentationEvent) {
	if d.id == presentation.deviceID {
		//fmt.Println("sanitaze OK!")
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
	if d.hasLifecycle {
		return
	}

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
