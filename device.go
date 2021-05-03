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
	connected       bool
	isFamiliar      bool
	configFile      string
	isReceivingData bool
	serial          *serial.Port
	configuration   Config
	cRec            chan DeviceMsg
	cSnd            chan DeviceMsg
	connectionTime  time.Time
	lastSeen        time.Time
}

type DeviceMsg struct {
	device *Device
	msg    string
}

var cDevicesReceive = make(chan DeviceMsg)

func newEmptyDevice(port string) *Device {
	return &Device{
		port: port,
		cRec: make(chan DeviceMsg),
		cSnd: make(chan DeviceMsg),
	}
}

//func newDevice(deviceID string, port string) (Device, error) {
//  fmt.Println(deviceID)
//  deviceID = strings.Split(deviceID, ";")[0]
//  params := strings.Split(deviceID, ",")
//  if len(params) != 5 {
//    return Device{},
//      errors.New("wrong event format")
//  }
//
//  var intParams []int
//  for _, param := range params {
//    intParam, err := strconv.Atoi(param)
//    if err != nil {
//      return Device{}, err
//    }
//    intParams = append(intParams, intParam)
//  }
//
//  receive := false
//  if intParams[2] != 0 {
//    receive = true
//  }
//
//  configFile := "config_" + strconv.Itoa(intParams[1])
//  configuration := readConfigurationFromFile(configFile, "default")
//
//  return Device{
//    intParams[1],
//    port,
//    true,
//    configFile,
//    receive,
//    configuration,
//  }, nil
//}

func (d *Device) connect() {
	config := &serial.Config{Name: d.port, Baud: 57600}
	s, err := serial.OpenPort(config)
	if err != nil {
		d.connected = false
	} else {
		d.serial = s
		d.connected = true
		d.connectionTime = time.Now()
	}
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

func (d *Device) poke() {
	if !d.connected && time.Now().Before(d.connectionTime.Add(time.Second*4)) {
		return
	}
	d.serial.Write([]byte("?;"))
	//go d.serial.Write([]byte("?;"))
	//time.Sleep(time.Millisecond * 20)
	//buf := make([]byte, 128)
	//n, err := d.serial.Read(buf)
	//if err != nil {
	//  log.Fatal(err)
	//} else {
	//  fmt.Println(string(buf[:n]))
	//  if string(buf[0]) == "3" {
	//    //d, err := newDevice(string(buf[:n]), port)
	//    if err != nil {
	//      fmt.Println(err)
	//    }
	//    //myDevices = append(myDevices, &d)
	//  }
	//}
}

func (d Device) sanitizeCheck(presentation devicePresentationEvent) {
	if d.id == presentation.deviceID {
		fmt.Println("sanitaze OK!")
	}
}

func (d *Device) updateConfiguration(planeName string) {
	c := readConfigurationFromFile(d, planeName)
	d.configuration = c
}

func (d *Device) getConfiguration() Config {
	return d.configuration
}

func (d *Device) lifecycle() {

	pokeTimer := time.NewTicker(time.Second * 1)

	for {
		select {
		case _ = <-pokeTimer.C:
			d.poke()
		}
	}
	//for {
	//  time.Sleep(time.Millisecond * 2000)
	//  d.poke()
	//  time.Sleep(time.Millisecond * 2000)
	//  d.poke()
	//  d.somethingElse()
	//}
}

func printConfig(device *Device) {
	fmt.Println(device.configuration)
}
