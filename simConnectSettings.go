package main

import (
	"fmt"
	"github.com/micmonay/simconnect"
	"time"
)

func scConnect(appName string) (*simconnect.EasySimConnect, error) {
	sc, err := simconnect.NewEasySimConnect()
	if err != nil {
		fmt.Println("Cannot connect to ", appName+"...")
		panic(err)
	}
	sc.SetLoggerLevel(simconnect.LogError)

	var connected bool = false
	fmt.Println("connecting to", appName+"...")
	for !connected {
		c, err := sc.Connect("Com_listener")
		if err != nil {
			time.Sleep(time.Second * 2)
			continue
		} else {
			connected = true
		}
		<-c // Wait serial confirmation
	}

	return sc, nil
}
