package main

import (
	"errors"
	"fmt"
	"github.com/micmonay/simconnect"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
)

type deviceEvent struct {
	eventType int
	deviceID  int
	value1    int
	value2    int
	value3    int
}

type DeviceActionEvent struct {
	device          *Device
	componentType   int
	componentNumber int
	action          int
}

type DevicePresentationEvent struct {
	device          *Device
	deviceID        int
	isReceivingData bool
	data            int
	data2           int
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
			logrus.Errorf("skiped event")
			continue
			//return nil, err
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

func executeEvents(event deviceEvent, device *Device) {
	device.connection.lastSeen = time.Now()
	if event.eventType == 1 {
		executeActionEvent(
			DeviceActionEvent{
				device,
				event.value1,
				event.value2,
				event.value3,
			},
		)
	}
	if event.eventType == 3 {
		executePresentationEvent(
			DevicePresentationEvent{
				device,
				event.deviceID,
				event.value1 != 0,
				event.value2,
				event.value3,
			},
		)
	}
}

func executeActionEvent(event DeviceActionEvent) {
	if eventSC == nil || !eventSC.IsAlive() {
		return
	}

	if event.componentType == 1 {
		if len(event.device.configuration.Elements.Buttons) >= event.componentNumber &&
			len(event.device.configuration.Elements.Buttons[event.componentNumber-1]) >= event.action {
			ev := event.device.configuration.Elements.Buttons[event.componentNumber-1][event.action-1]
			for _, singleEvent := range ev.SimEvent {
				fmt.Println("Button", event.componentNumber, event.action, singleEvent.Event)
				event := eventSC.NewSimEvent(simconnect.KeySimEvent(singleEvent.Event))
				event.RunWithValue(singleEvent.Value)
			}
		}
	}
	if event.componentType == 2 {
		if len(event.device.configuration.Elements.Switches) >= event.componentNumber &&
			len(event.device.configuration.Elements.Switches[event.componentNumber-1]) >= event.action {
			ev := event.device.configuration.Elements.Switches[event.componentNumber-1][event.action-1]
			for _, singleEvent := range ev.SimEvent {
				fmt.Println("Switches", event.componentNumber, event.action, singleEvent.Event)
				event := eventSC.NewSimEvent(simconnect.KeySimEvent(singleEvent.Event))
				event.RunWithValue(singleEvent.Value)
			}
		}
	}
	if event.componentType == 3 {
		if len(event.device.configuration.Elements.Encoders) >= event.componentNumber &&
			len(event.device.configuration.Elements.Encoders[event.componentNumber-1]) >= event.action {
			ev := event.device.configuration.Elements.Encoders[event.componentNumber-1][event.action-1]
			for _, singleEvent := range ev.SimEvent {
				fmt.Println("Encoder", event.componentNumber, event.action, singleEvent.Event)
				event := eventSC.NewSimEvent(simconnect.KeySimEvent(singleEvent.Event))
				event.RunWithValue(singleEvent.Value)
			}
		}
	}
}

func executePresentationEvent(event DevicePresentationEvent) {
	if event.device.id == 0 {
		event.device.id = event.deviceID
		event.device.isFamiliar = true
		event.device.isReceivingData = event.isReceivingData
		event.device.updateConfiguration()
		event.device.writeTo()
		return
	}
	event.device.sanitizeCheck(event)
}
