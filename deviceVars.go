package main

import (
	"fmt"
	"github.com/micmonay/simconnect"
	"log"
	"strings"
)

func connectToSimVars(sc *simconnect.EasySimConnect) <-chan []simconnect.SimVar {
	cSimVar, err := sc.ConnectToSimVar(
		simconnect.SimVarPlaneAltitude(),
		simconnect.SimVarPlaneLatitude(simconnect.UnitDegrees),
		simconnect.SimVarPlaneLongitude(),
		simconnect.SimVarIndicatedAltitude(),
		simconnect.SimVarGeneralEngRpm(1),
		simconnect.SimVarAutopilotMaster(),
	)
	if err != nil {
		panic(err)
	}

	for {
		fmt.Println(sc.IsAlive())
		result := <-cSimVar
		fmt.Println(result)
		//for _, simVar := range result {
		//  var f float64
		//  var err error
		//  if strings.Contains(string(simVar.Unit), "String") {
		//    log.Printf("%s : %#v\n", simVar.Name, simVar.GetString())
		//  } else if simVar.Unit == "SIMCONNECT_DATA_LATLONALT" {
		//    data, _ := simVar.GetDataLatLonAlt()
		//    log.Printf("%s : %#v\n", simVar.Name, data)
		//  } else if simVar.Unit == "SIMCONNECT_DATA_XYZ" {
		//    data, _ := simVar.GetDataXYZ()
		//    log.Printf("%s : %#v\n", simVar.Name, data)
		//  } else if simVar.Unit == "SIMCONNECT_DATA_WAYPOINT" {
		//    data, _ := simVar.GetDataWaypoint()
		//    log.Printf("%s : %#v\n", simVar.Name, data)
		//  } else {
		//    f, err = simVar.GetFloat64()
		//    log.Println(simVar.Name, fmt.Sprintf("%f", f))
		//  }
		//  if err != nil {
		//    log.Println("return error :", err)
		//    panic(err)
		//  }
		//}
	}
}

func goSimVars(c <-chan []simconnect.SimVar) {
	for {
		fmt.Println("hello")
		result := <-c
		fmt.Println(result)
		for _, simVar := range result {
			var f float64
			var err error
			if strings.Contains(string(simVar.Unit), "String") {
				log.Printf("%s : %#v\n", simVar.Name, simVar.GetString())
			} else if simVar.Unit == "SIMCONNECT_DATA_LATLONALT" {
				data, _ := simVar.GetDataLatLonAlt()
				log.Printf("%s : %#v\n", simVar.Name, data)
			} else if simVar.Unit == "SIMCONNECT_DATA_XYZ" {
				data, _ := simVar.GetDataXYZ()
				log.Printf("%s : %#v\n", simVar.Name, data)
			} else if simVar.Unit == "SIMCONNECT_DATA_WAYPOINT" {
				data, _ := simVar.GetDataWaypoint()
				log.Printf("%s : %#v\n", simVar.Name, data)
			} else {
				f, err = simVar.GetFloat64()
				log.Println(simVar.Name, fmt.Sprintf("%f", f))
			}
			if err != nil {
				log.Println("return error :", err)
				panic(err)
			}
		}
	}
}

func connect() *simconnect.EasySimConnect {
	sc, err := simconnect.NewEasySimConnect()
	if err != nil {
		panic(err)
	}
	sc.SetLoggerLevel(simconnect.LogInfo)
	c, err := sc.Connect("MyApp")
	if err != nil {
		panic(err)
	}
	<-c // wait connection confirmation
	for {
		if <-sc.ConnectSysEventSim() {
			break // wait simconnect start
		}
	}
	return sc
}

func Example_getSimVar() {
	sc := connect()
	cSimVar, err := sc.ConnectToSimVar(
		simconnect.SimVarPlaneAltitude(),
		simconnect.SimVarPlaneLatitude(simconnect.UnitDegrees), // you can force the units
		simconnect.SimVarPlaneLongitude(),
		simconnect.SimVarIndicatedAltitude(),
		simconnect.SimVarAutopilotAltitudeLockVar(),
		simconnect.SimVarAutopilotMaster(),
	)
	if err != nil {
		log.Fatalln(err)
	}
	for i := 0; i < 1; i++ {
		result := <-cSimVar
		for _, simVar := range result {
			f, err := simVar.GetFloat64()
			if err != nil {
				panic(err)
			}
			log.Printf("%#v\n", f)
		}

	}
	<-sc.Close() // wait close confirmation
	// Output:

}
