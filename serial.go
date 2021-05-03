package main

import (
	"errors"
	"fmt"
	"github.com/tarm/serial"
	"log"
	"os"
	"strconv"
	"strings"
)

func findActivePorts() ([]string, []string) {
	var openPorts []string
	var busyPorts []string

	for i := 1; i < 100; i++ {
		portName := "COM" + strconv.Itoa(i)

		config := &serial.Config{Name: portName, Baud: 57600}
		s, err := serial.OpenPort(config)
		if err != nil {
			if errors.Is(err, os.ErrPermission) {
				busyPorts = append(busyPorts, portName)
			}
			//if err.Error() == "Access is denied." {
			//	busyPorts = append(busyPorts, portName)
			//}
		} else {
			openPorts = append(openPorts, portName)
			_ = s.Close()
		}
	}

	return openPorts, busyPorts
}

func readCom(s *serial.Port, c chan string) {
	for {
		buf := make([]byte, 100)
		n, err := s.Read(buf)
		if err != nil {
			log.Fatal(err)
		}
		c <- string(buf[:n])
	}
}

func writeCom(s *serial.Port, c chan []string) {
	currentTransmision := ""
	lastTransmision := ""
	for vars := range c {
		currentTransmision = strings.Join(vars, ",")
		if currentTransmision != lastTransmision {
			transmit := currentTransmision + ",epmtyData" + ";"
			fmt.Println(transmit)
			s.Write([]byte(transmit))
			lastTransmision = currentTransmision
		}
	}
}
