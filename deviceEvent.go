package main

import (
	"errors"
	"fmt"
	"github.com/micmonay/simconnect"
	"strconv"
	"strings"
)

type deviceEvent struct {
	eventType       int
	device          int
	componentType   int
	componentNumber int
	action          int
}

var portChannelMessageBuffer string

func collectEvents(msg string) ([]deviceEvent, error) {
	var newActions []string
	var deviceEvents []deviceEvent

	portChannelMessageBuffer += msg
	newActions = strings.Split(portChannelMessageBuffer, ";")
	if len(newActions) == 0 || newActions[len(newActions)-1] != "" {
		return nil, nil
	}

	for _, newAction := range newActions {
		if newAction == "" {
			continue
		}
		de, err := decodeEvent(newAction)
		if err != nil {
			return nil, err
		} else {
			deviceEvents = append(deviceEvents, de)
		}
	}
	portChannelMessageBuffer = ""
	return deviceEvents, nil
}

func decodeEvent(event string) (deviceEvent, error) {
	params := strings.Split(event, ",")
	if len(params) != 5 {
		return deviceEvent{},
			errors.New("wrong event format")
	}

	var intParams []int
	for _, param := range params {
		intParam, err := strconv.Atoi(param)
		if err != nil {
			return deviceEvent{}, err
		}
		intParams = append(intParams, intParam)
	}

	return deviceEvent{
		intParams[0],
		intParams[1],
		intParams[2],
		intParams[3],
		intParams[4],
	}, nil
}

func executeEvents(sc *simconnect.EasySimConnect, event deviceEvent, devices []*Device) {
	if event.eventType != 1 {
		return
	}

	for _, device := range devices {
		if device.id == event.device {
			if event.componentType == 1 {
				if len(device.configuration.Elements.Buttons) >= event.componentNumber {
					ev := device.configuration.Elements.Buttons[event.componentNumber-1][event.action-1]
					event := sc.NewSimEvent(simconnect.KeySimEvent(ev.SimEvent))
					event.RunWithValue(ev.Value)
					fmt.Println(ev.SimEvent)
				}
			}
			if event.componentType == 2 {
				if len(device.configuration.Elements.Switches) >= event.componentNumber {
					ev := device.configuration.Elements.Switches[event.componentNumber-1][event.action-1]
					event := sc.NewSimEvent(simconnect.KeySimEvent(ev.SimEvent))
					event.RunWithValue(ev.Value)
					fmt.Println(ev.SimEvent)
				}
			}
			if event.componentType == 3 {
				if len(device.configuration.Elements.Encoders) >= event.componentNumber {
					ev := device.configuration.Elements.Encoders[event.componentNumber-1][event.action-1]
					event := sc.NewSimEvent(simconnect.KeySimEvent(ev.SimEvent))
					event.RunWithValue(ev.Value)
					fmt.Println(ev.SimEvent)
				}
			}
		}
	}
}

//if event.eventType == 1{
//  if event.componentType == 1 {
//    if event.componentNumber == 1 {
//      event := sc.NewSimEvent(sim.KeyG1000PfdMenuButton)
//      event.RunWithValue(100)
//      fmt.Println("KeyG1000PfdMenuButton")
//    }
//    if event.componentNumber == 2 {
//      event := sc.NewSimEvent(sim.KeyG1000MfdMenuButton)
//      event.RunWithValue(100)
//      fmt.Println("KeyG1000MfdMenuButton")
//    }
//  }
//  if event.componentType == 3 {
//    if event.componentNumber == 1 {
//      if event.action == 1 {
//        event := sc.NewSimEvent(sim.KeyApAltVarDec)
//        event.RunWithValue(100)
//        fmt.Println("KeyApAltVarDec")
//      }
//      if event.action == 2 {
//        event := sc.NewSimEvent(sim.KeyApAltVarInc)
//        event.Run()
//        fmt.Println("KeyApAltVarInc")
//      }
//    }
//    if event.componentNumber == 2 {
//      if event.action == 1 {
//        event := sc.NewSimEvent(sim.KeyComRadioFractInc)
//        event.Run()
//        fmt.Println("KeyComRadioFractInc")
//      }
//      if event.action == 2 {
//        event := sc.NewSimEvent(sim.KeyComRadioFractDec)
//        event.Run()
//        fmt.Println("KeyComRadioFractDec")
//      }
//    }
//  }
//}
//}
