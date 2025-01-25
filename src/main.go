package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"LightControl/src/Config"
	"LightControl/src/Extensions"
	"LightControl/src/Hue"

	_ "github.com/shirou/gopsutil/v4/process"
)

const AsciiLogo string = " _     _       _     _   _____             _             _   \n" +
	"| |   (_)     | |   | | /  __ \\           | |           | |\n" +
	"| |    _  __ _| |__ | |_| /  \\/ ___  _ __ | |_ _ __ ___ | |\n" +
	"| |   | |/ _` | '_ \\| __| |    / _ \\| '_ \\| __| '__/ _ \\| |\n" +
	"| |___| | (_| | | | | |_| \\__/\\ (_) | | | | |_| | | (_) | |\n" +
	"\\_____/_|\\__, |_| |_|\\__|\\____/\\___/|_| |_|\\__|_|  \\___/|_|\n" +
	"            __/ |\n" +
	"           |___/"

func ActivityWatcher(controller *Hue.Controller, lightId string, configLock *Extensions.StructLock[Config.Config]) {
	for {
		config := configLock.Get()
		color, err := GetProcessColor(&config)
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
	fmt.Println(AsciiLogo)

	configLock, err := Config.Watcher()
	if err != nil {
		panic(err)
	}

	hueConfigLock, err := Config.HueWatcher(configLock)
	if err != nil {
		panic(err)
	}

	config := configLock.Get()

	controller := Hue.NewController(hueConfigLock)
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

	go ActivityWatcher(controller, targetLightId, configLock)

	waitForShutDown(controller, targetLightId)
}

func waitForShutDown(controller *Hue.Controller, targetLightId string) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGABRT, os.Interrupt)

	select {
	case <-sigs:
		// Turn off the light when the program exits
		err := controller.SetOnOff(targetLightId, false)
		if err != nil {
			fmt.Println(err)
		}

	}

}
