package main

import (
	"errors"
	"math"
	"strings"
	"time"

	"LightControl/Dates"
	"LightControl/Hue"
	"github.com/shirou/gopsutil/v4/process"
)

func GetProcessColor(config *Config) (Hue.Color, error) {
	processes, err := process.Processes()
	if err != nil {
		panic(err)
	}
	priority := math.MaxInt
	var resultColor Hue.Color

	for _, process := range processes {
		var name string
		name, err = process.Name()
		if err != nil {
			continue
		}
		name = strings.ToLower(name)

		i := 0

		now := time.Now()
		dayTime := Dates.NewDayTime(now.Hour(), now.Minute(), now.Second())
		for programName, timedProgram := range config.TimedPrograms {
			if strings.Contains(name, programName) && timedProgram.Span.Contains(dayTime) {
				_, err = timedProgram.Color.GetColor()
				if err != nil {
					continue
				}

				priority = i
				resultColor = timedProgram.Color
			}
		}
		for programName, color := range config.Programs {
			if strings.Contains(name, programName) && i < priority {
				// Check if color is valid
				_, err = color.GetColor()
				if err != nil {
					continue
				}

				priority = i
				resultColor = color
			}
			i++
		}
	}

	if priority == math.MaxInt {
		return nil, errors.New("no matching process found")
	}

	return resultColor, nil
}
