package main

import (
	"fmt"
	"github.com/micmonay/simconnect"
	"log"
	"strings"
	"time"
)

func registerVars(sc *simconnect.EasySimConnect, vars []simconnect.SimVar) <-chan []simconnect.SimVar {
	cSimVar, err := sc.ConnectToSimVar(vars...)
	if err != nil {
		fmt.Println("Can not register Vars", err)
	}

	return cSimVar
}

func startGettingVars(c <-chan []simconnect.SimVar, cSnd chan []string) error {
	for {
		var vars []string
		select {
		case result := <-c:
			for _, simVar := range result {
				var f float64
				var err error
				if strings.Contains(string(simVar.Unit), "String") {
					//log.Printf("%s : %#v\n", simVar.Name, simVar.GetString())
					vars = append(vars, simVar.GetString())
				} else if simVar.Unit == "SIMCONNECT_DATA_LATLONALT" {
					//data, _ := simVar.GetDataLatLonAlt()
					//log.Printf("%s : %#v\n", simVar.Name, data)
				} else if simVar.Unit == "SIMCONNECT_DATA_XYZ" {
					//data, _ := simVar.GetDataXYZ()
					//log.Printf("%s : %#v\n", simVar.Name, data)
				} else if simVar.Unit == "SIMCONNECT_DATA_WAYPOINT" {
					//data, _ := simVar.GetDataWaypoint()
					//log.Printf("%s : %#v\n", simVar.Name, data)
				} else {
					f, err = simVar.GetFloat64()
					//log.Println(simVar.Name, fmt.Sprintf("%f", f))
					vars = append(vars, fmt.Sprintf("%f", f))
				}
				if err != nil {
					log.Println("return error :", err)
					panic(err)
				}
			}
		case <-time.After(5 * time.Second):
			return fmt.Errorf("Can't Connect to MSFS")
		}

		cSnd <- vars
	}
}

func testConnection(c chan string) {
	sc, _ := scConnect("MSFS_test")

	cSimVar, err := sc.ConnectToSimVar(
		simconnect.SimVarPlaneAltitude(),
	)
	if err != nil {
		fmt.Println("Can not register Vars")
	}

	select {
	case _ = <-cSimVar:
		c <- "serial success"
	case <-time.After(10 * time.Second):
		{
			sc.Close()
			fmt.Println("serial failed")
			c <- "serial fail"
		}
	}
}
