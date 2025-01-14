package main

import (
	"encoding/json"
	"errors"
	"os"
	"strings"

	"LightControl/Dates"
	"LightControl/Hue"
)

const ConfigPath = "config.json"

type RawConfig struct {
	ID              string                     `json:"id"`
	Address         string                     `json:"address"`
	TargetLightName string                     `json:"target_light_name"`
	DefaultColor    string                     `json:"default_color"`
	Programs        map[string]string          `json:"programs"`
	TimedPrograms   map[string]RawTimedProgram `json:"timed_programs"`
}

type RawTimedProgram struct {
	TimeSpan string `json:"span"`
	Color    string `json:"color"`
}

type Config struct {
	ID              string
	Address         string
	TargetLightName string
	DefaultColor    Hue.Color
	Programs        map[string]Hue.Color
	TimedPrograms   map[string]TimedProgram
}

type TimedProgram struct {
	Span  Dates.IDayTimeSpan
	Color Hue.Color
}

func (r *RawConfig) ToConfig() (*Config, error) {
	defaultColor, err := Hue.NewRGBFromHex(r.DefaultColor)
	if err != nil {
		return nil, err
	}

	programs := make(map[string]Hue.Color)
	for name, colorHex := range r.Programs {
		color, err := Hue.NewRGBFromHex(colorHex)
		if err != nil {
			return nil, err
		}

		programs[strings.ToLower(name)] = &color
	}

	timedPrograms := make(map[string]TimedProgram)
	for name, timedProgram := range r.TimedPrograms {
		var color Hue.RGBColor
		var timeSpan Dates.IDayTimeSpan

		color, err = Hue.NewRGBFromHex(timedProgram.Color)
		if err != nil {
			return nil, err
		}

		timeSpan, ok := Dates.ParseStringToSpan(timedProgram.TimeSpan)
		if !ok {
			timeSpan, ok = Dates.ParseStringToSpanSolar(timedProgram.TimeSpan)
		}
		if !ok {
			return nil, errors.New("invalid time span")
		}

		timedPrograms[strings.ToLower(name)] = TimedProgram{
			timeSpan,
			&color,
		}
	}

	return &Config{
		ID:              r.ID,
		Address:         r.Address,
		TargetLightName: strings.ToLower(r.TargetLightName),
		DefaultColor:    &defaultColor,
		Programs:        programs,
		TimedPrograms:   timedPrograms,
	}, nil
}

func ReadConfig() (*Config, error) {
	file, err := os.Open(ConfigPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var rawConfig RawConfig
	if err := json.NewDecoder(file).Decode(&rawConfig); err != nil {
		return nil, err
	}

	config, err := rawConfig.ToConfig()

	if err != nil {
		return nil, err
	}

	return config, nil
}
