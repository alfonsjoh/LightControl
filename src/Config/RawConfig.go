package Config

import (
	"encoding/json"
	"errors"
	"os"
	"strings"

	"LightControl/src/Dates"
	"LightControl/src/Hue/Colors"
)

const ConfigPath = "config.json"

type RawConfig struct {
	TargetLightName string            `json:"target_light_name"`
	DefaultColor    string            `json:"default_color"`
	Programs        []RawProgram      `json:"programs"`
	TimedPrograms   []RawTimedProgram `json:"timed_programs"`
}

type RawProgram struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

type RawTimedProgram struct {
	Name     string `json:"name"`
	TimeSpan string `json:"span"`
	Color    string `json:"color"`
}

func parseColor(color string, defaultColor Colors.Color) (Colors.Color, error) {
	if strings.ToLower(color) == "default" {
		return defaultColor, nil
	}

	rgb, err := Colors.NewRGBFromHex(color)
	if err != nil {
		return nil, err
	}

	return &rgb, nil
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

func SetHueCredentials(config *Config) error {
	// Read from environment variable: "HueCredentials"
	credentials := os.Getenv("HueCredentials")
	if credentials == "" {
		return errors.New("HueCredentials environment variable not set")
	}

	// Parse credentials formatted like "id@ip"
	split := strings.Split(credentials, "@")
	if len(split) != 2 {
		return errors.New("invalid HueCredentials format")
	}

	config.Address = split[1]
	config.ID = split[0]

	return nil
}

func (r *RawConfig) ToConfig() (*Config, error) {
	defaultColor, err := Colors.NewRGBFromHex(r.DefaultColor)
	if err != nil {
		return nil, err
	}

	programs := make([]Program, 0, len(r.Programs))
	for _, program := range r.Programs {
		color, err := parseColor(program.Color, &defaultColor)
		if err != nil {
			return nil, err
		}

		programs = append(programs, Program{
			Name:  strings.ToLower(program.Name),
			Color: color})
	}

	timedPrograms := make([]TimedProgram, 0, len(r.TimedPrograms))
	for _, timedProgram := range r.TimedPrograms {
		var color Colors.Color
		color, err = parseColor(timedProgram.Color, &defaultColor)
		if err != nil {
			return nil, err
		}

		var timeSpan Dates.IDayTimeSpan
		timeSpan, ok := Dates.ParseStringToSpan(timedProgram.TimeSpan)
		if !ok {
			timeSpan, ok = Dates.ParseStringToSpanSolar(timedProgram.TimeSpan)
		}
		if !ok {
			return nil, errors.New("invalid time span")
		}

		timedPrograms = append(timedPrograms, TimedProgram{
			strings.ToLower(timedProgram.Name),
			timeSpan,
			color,
		})
	}

	config := &Config{
		TargetLightName: strings.ToLower(r.TargetLightName),
		DefaultColor:    &defaultColor,
		Programs:        programs,
		TimedPrograms:   timedPrograms,
	}

	err = SetHueCredentials(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
