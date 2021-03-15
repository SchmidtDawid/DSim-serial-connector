package main

import (
	"fmt"
	"github.com/tarm/serial"
	"log"
	"strconv"
	"time"
)

func findOpenPorts() ([]string, []string) {
	var openPorts []string
	var closedPorts []string

	for i := 1; i < 50; i++ {
		portName := "COM" + strconv.Itoa(i)

		config := &serial.Config{Name: portName, Baud: 57600}
		s, err := serial.OpenPort(config)
		if err != nil {
			closedPorts = append(closedPorts, portName)
		} else {
			openPorts = append(openPorts, portName)
			_ = s.Close()
		}

	}

	return openPorts, closedPorts
}

func checkDeviceOnPorts(ports []string) []*Device {
	var myDevices []*Device

	for _, port := range ports {
		//fmt.Println("test", port)
		config := &serial.Config{Name: port, Baud: 57600, ReadTimeout: time.Millisecond * 100}
		s, err := serial.OpenPort(config)
		if err != nil {
			fmt.Println(err)
		} else {
			go s.Write([]byte("?"))
			time.Sleep(time.Millisecond * 1000)
			buf := make([]byte, 128)
			n, err := s.Read(buf)
			if err != nil {
				log.Fatal(err)
			} else {
				fmt.Println(string(buf[:n]))
				if string(buf[0]) == "3" {
					d, err := newDevice(string(buf[:n]), port)
					if err != nil {
						fmt.Println(err)
					}
					myDevices = append(myDevices, &d)
				}
			}
			_ = s.Close()
		}
	}

	return myDevices
}

func readPorts(ports []string, c chan string) {
	for _, port := range ports {
		config := &serial.Config{Name: port, Baud: 57600}
		s, err := serial.OpenPort(config)
		if err != nil {
			fmt.Println(err)
		} else {
			go readCom(s, c)
		}
	}
}

func readDevices(devices []*Device, c chan string) {
	for _, device := range devices {
		config := &serial.Config{Name: device.port, Baud: 57600}
		s, err := serial.OpenPort(config)
		if err != nil {
			fmt.Println(err)
		} else {
			go readCom(s, c)
		}
	}
}

func readCom(s *serial.Port, c chan string) {
	for {
		buf := make([]byte, 30)
		n, err := s.Read(buf)
		if err != nil {
			log.Fatal(err)
		}
		c <- string(buf[:n])
	}
}