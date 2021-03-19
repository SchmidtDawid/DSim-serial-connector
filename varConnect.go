package main

import (
	"fmt"
	"github.com/micmonay/simconnect"
	"log"
	"strings"
)

func startGettingVars(sc *simconnect.EasySimConnect) {
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
		result := <-cSimVar
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
