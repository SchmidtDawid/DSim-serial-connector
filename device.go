package main

import (
	"context"
	"fmt"
	"github.com/micmonay/simconnect"
	"github.com/sirupsen/logrus"
	"github.com/tarm/serial"
	"strconv"
	"strings"
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
	cSnd            chan []string
	simData         simData
	eventsBuffer    []deviceEvent
	msgBuffer       string
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

type simData struct {
	data    []VarRequest
	sc      *simconnect.EasySimConnect
	cSimVar <-chan []simconnect.SimVar
}

type DeviceMsg struct {
	device *Device
	msg    string
}

func newEmptyDevice(port string) *Device {
	return &Device{
		port: port,
		cRec: make(chan DeviceMsg),
		cSnd: make(chan []string),
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
	d.connection.pokeInterval = time.Second * 30
	d.connection.begin = time.Now()
	d.connection.ctx, d.connection.cancel = context.WithCancel(context.Background())
	go d.listen()
	go d.processReceivedData()
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
	d.simData.sc.Close()
	logrus.Errorf(strconv.Itoa(d.id), " disconnected")
}

func (d *Device) listen() {
	if !d.connection.connected {
		return
	}
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

func (d *Device) processReceivedData() {
	if !d.connection.connected {
		return
	}
	for msg := range d.cRec {
		select {
		case <-d.connection.ctx.Done():
			return
		default:
			//fmt.Println(d.id)
			d.collectEvents(msg.msg)

			for len(d.eventsBuffer) > 0 {
				event := d.eventsBuffer[0]
				d.eventsBuffer = d.eventsBuffer[1:]
				executeEvents(event, d)
			}
		}
	}
}

func (d *Device) collectEvents(msg string) {
	var newActions []string
	var deviceEvents []deviceEvent

	d.msgBuffer += msg

	newActions = strings.Split(d.msgBuffer, ";")
	if len(newActions) == 0 || newActions[len(newActions)-1] != "" {
		return
	}

	for _, newAction := range newActions {
		if newAction == "" {
			continue
		}
		de, err := decodeEvent(newAction)
		fmt.Println(newAction)
		if err != nil {
			logrus.Errorf("skiped event")
			continue
		} else {
			deviceEvents = append(deviceEvents, de)
		}
	}
	d.msgBuffer = ""
	d.eventsBuffer = append(d.eventsBuffer, deviceEvents...)
}

func (d *Device) writeTo() {
	time.Sleep(time.Millisecond * 100)
	if d.isReceivingData {
		d.simData.sc, _ = scConnect(strconv.Itoa(d.id) + "_Vars")
		d.simData.sc.SetDelay(50 * time.Millisecond)

		d.simData.cSimVar = registerVars(d.simData.sc, d.getSimVars())
		go startGettingVars(d.simData.cSimVar, d.cSnd)

		go writeCom(d.serial, d.cSnd)
	}
}

func (d *Device) getSimVars() []simconnect.SimVar {
	var Vars []simconnect.SimVar
	for _, vReq := range d.configuration.Vars {
		Vars = append(Vars, simconnect.SimVar{
			Name:     vReq.Name,
			Unit:     simconnect.SimVarUnit(vReq.Unit),
			Settable: vReq.Settable,
		})
	}
	return Vars
}

func (d *Device) sanitize() {
	if !d.connection.connected || time.Now().Before(d.connection.begin.Add(time.Second*5)) {
		return
	}
	//fmt.Println("poke -> ", d.id)
	d.serial.Write([]byte("?;"))
	if time.Now().After(d.connection.lastSeen.Add(d.connection.pokeInterval * 3)) {
		fmt.Println("error? ", d.id)
		d.connection.connectionFailures++
		if d.connection.connectionFailures > 5 && d.id == 0 {
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
	configTimer := time.NewTicker(time.Second * 5)

	for {
		select {
		case _ = <-pokeTimer.C:
			d.sanitize()
		case _ = <-configTimer.C:
			d.updateConfiguration()
		}
	}
}
