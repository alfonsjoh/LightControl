package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"LightControl/Hue"
	_ "github.com/shirou/gopsutil/v4/process"
)

func ActivityWatcher(controller *Hue.Controller, lightId string, config *Config, doneCh chan struct{}) {
	defer func() {
		doneCh <- struct{}{}
	}()

	for {
		color, err := GetProcessColor(config)
		if err != nil {
			color = config.DefaultColor
		}
		err = controller.SetColor(lightId, color)
		if err != nil {
			fmt.Println(err)
		}

		time.Sleep(500 * time.Millisecond)
	}
}

func main() {
	config, err := ReadConfig()
	if err != nil {
		panic(err)
	}

	hueConfig := Hue.NewConfig(config.ID, config.Address)
	controller := Hue.NewController(hueConfig)
	lights, err := controller.GetAllLights()
	if err != nil {
		panic(err)
	}

	targetLightId := ""
	for id, light := range lights {
		if strings.ToLower(light.Name) == strings.ToLower(config.TargetLightName) {
			targetLightId = id
			break
		}
	}

	if targetLightId == "" {
		panic("Target light not found")
	}

	err = controller.SetColor(targetLightId, config.DefaultColor)
	if err != nil {
		panic(err)
	}

	err = controller.SetOnOff(targetLightId, true)
	if err != nil {
		panic(err)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	doneCh := make(chan struct{})

	go func() {
		<-sigs
		// Turn off the light when the program exits
		err := controller.SetOnOff(targetLightId, false)
		if err != nil {
			fmt.Println(err)
		}
		close(doneCh)
	}()

	go ActivityWatcher(controller, targetLightId, config, doneCh)

	<-doneCh

}
