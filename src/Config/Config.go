package Config

import (
	"bytes"
	"fmt"
	"os"
	"time"

	"LightControl/src/Extensions"
	"LightControl/src/Hue"
	"LightControl/src/Hue/Colors"
)

type Config struct {
	ID              string
	Address         string
	TargetLightName string
	DefaultColor    Colors.Color
	Triggers        []ColorTrigger
	Hash            []byte // 128-bit hash
}

func (c *Config) Equals(other *Config) bool {
	return bytes.Equal(c.Hash, other.Hash)
}

func Watcher() (*Extensions.StructLock[Config], error) {
	prevFileStat, err := os.Stat(ConfigPath)
	if err != nil {
		return nil, err
	}

	config, err := ReadConfig()
	if err != nil {
		return nil, err
	}

	configLock := Extensions.NewStructLock[Config](*config)

	go func(prevFileStat os.FileInfo, configLock *Extensions.StructLock[Config], beginConfig *Config) {
		for {
			// Sleep 1 second
			time.Sleep(1 * time.Second)

			fileStat, err := os.Stat(ConfigPath)
			if err != nil {
				continue
			}

			// Continue if the file hasn't changed
			if fileStat.ModTime() == prevFileStat.ModTime() && fileStat.Size() == prevFileStat.Size() {
				continue
			}
			prevFileStat = fileStat

			config, err := ReadConfig()
			if err != nil {
				fmt.Println(err)
				continue
			}

			// Clear the channel
			fmt.Println("Config updated")

			// Send the new config to all subscribers
			configLock.Set(*config)
		}
	}(prevFileStat, configLock, config)

	return configLock, nil
}

func HueWatcher(configLock *Extensions.StructLock[Config]) (*Extensions.StructLock[Hue.Config], error) {
	config := configLock.Get()
	hueConfig := Hue.NewConfig(config.ID, config.Address)
	hueConfigLock := Extensions.NewStructLock[Hue.Config](hueConfig)

	go func(configLock *Extensions.StructLock[Config], hueConfigLock *Extensions.StructLock[Hue.Config]) {
		prevConfig := configLock.Get()
		for {
			// Sleep 1 second
			time.Sleep(1 * time.Second)
			config := configLock.Get()
			if config.Equals(&prevConfig) {
				continue
			}
			prevConfig = config

			fmt.Println("Hue config updated")

			hueConfig := Hue.NewConfig(config.ID, config.Address)
			hueConfigLock.Set(hueConfig)

		}
	}(configLock, hueConfigLock)

	return hueConfigLock, nil
}
