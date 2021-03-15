package main

import (
	"fmt"
	"github.com/micmonay/simconnect"
	"time"
)

func main() {

	portChannel := make(chan string)
	ports, _ := findOpenPorts()

	fmt.Println(ports)
	myDevices := checkDeviceOnPorts(ports)
	fmt.Println(myDevices)
	readDevices(myDevices, portChannel)

	go keepUpdateConfig(myDevices)

	//----------------------------------------
	sc, err := simconnect.NewEasySimConnect()
	if err != nil {
		panic(err)
	}
	sc.SetLoggerLevel(simconnect.LogInfo)

	var connected bool = false
	fmt.Println("connecting to MSFS...")
	for !connected {
		c, err := sc.Connect("Com_listener")
		if err != nil {
			time.Sleep(time.Second * 2)
			continue
		} else {
			connected = true
		}
		<-c // Wait connection confirmation
	}

	var deviceEvents []deviceEvent
	for msg := range portChannel {
		incomingEvents, err := collectEvents(msg)
		if err != nil {
		}
		deviceEvents = append(deviceEvents, incomingEvents...)
		for len(deviceEvents) > 0 {
			event := deviceEvents[0]
			deviceEvents = deviceEvents[1:]
			executeEvents(sc, event, myDevices)
		}
	}

	//cSimVar, err := sc.ConnectToSimVar(
	// sim.SimVarPlaneAltitude(),
	// sim.SimVarPlaneLatitude(sim.UnitDegrees), // You can force the units
	// sim.SimVarPlaneLongitude(),
	// sim.SimVarIndicatedAltitude(),
	// sim.SimVarGeneralEngRpm(1),
	// sim.SimVarAutopilotMaster(),
	//)
	//if err != nil {
	// panic(err)
	//}
	//cSimStatus := sc.ConnectSysEventSim()
	////wait sim start
	//for {
	// if <-cSimStatus {
	//   break
	// }
	//}
	//crashed := sc.ConnectSysEventCrashed()
	//for {
	// select {
	// case result := <-cSimVar:
	//   for _, simVar := range result {
	//     var f float64
	//     var err error
	//     if strings.Contains(string(simVar.Unit), "String") {
	//       log.Printf("%s : %#v\n", simVar.Name, simVar.GetString())
	//     } else if simVar.Unit == "SIMCONNECT_DATA_LATLONALT" {
	//       data, _ := simVar.GetDataLatLonAlt()
	//       log.Printf("%s : %#v\n", simVar.Name, data)
	//     } else if simVar.Unit == "SIMCONNECT_DATA_XYZ" {
	//       data, _ := simVar.GetDataXYZ()
	//       log.Printf("%s : %#v\n", simVar.Name, data)
	//     } else if simVar.Unit == "SIMCONNECT_DATA_WAYPOINT" {
	//       data, _ := simVar.GetDataWaypoint()
	//       log.Printf("%s : %#v\n", simVar.Name, data)
	//     } else {
	//       f, err = simVar.GetFloat64()
	//       log.Println(simVar.Name, fmt.Sprintf("%f", f))
	//     }
	//     if err != nil {
	//       log.Println("return error :", err)
	//     }
	//   }
	//
	// case <-crashed:
	//   log.Println("Your are crashed !!")
	//   <-sc.Close() // Wait close confirmation
	//   return       // This example close after crash in the sim
	// }
	//
	// fmt.Println(<-portChannel)
	//
	//}
}

//func upd(devices) {
//  myDevices[0].updateConfiguration()
//  time.Sleep(time.Second * 2)
//}
