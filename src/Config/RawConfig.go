package Config

import (
	"encoding/json"
	"errors"
	"hash/fnv"
	"io"
	"os"
	"strings"

	"LightControl/src/Dates"
	"LightControl/src/Hue/Colors"
	"LightControl/src/Triggers"
)

const ConfigPath = "config.json"

type RawConfig struct {
	TargetLightName string            `json:"target_light_name"`
	Colors          map[string]string `json:"colors"`
	ProgramGroups   []RawProgramGroup `json:"program_groups"`
	Programs        []RawProgram      `json:"programs"`
	TimedPrograms   []RawTimedProgram `json:"timed_programs"`
}

type RawProgramGroup struct {
	Name     string   `json:"name"`
	Programs []string `json:"programs"`
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

func computeHash(data []byte) []byte {
	h := fnv.New128()
	h.Write(data)
	return h.Sum(nil)
}

func parseColor(color string, colors *map[string]Colors.Color) (Colors.Color, error) {
	if c, ok := (*colors)[strings.ToLower(color)]; ok {
		return c, nil
	}

	rgb, err := Colors.NewRGBFromHex(color)
	if err != nil {
		return nil, err
	}

	return &rgb, nil
}

func setHueCredentials(config *Config) error {
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

func ReadConfig() (*Config, error) {
	file, err := os.Open(ConfigPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var rawConfig RawConfig
	if err = json.Unmarshal(bytes, &rawConfig); err != nil {
		return nil, err
	}

	config, err := rawConfig.ToConfig()
	if err != nil {
		return nil, err
	}

	config.Hash = computeHash(bytes)

	return config, nil
}

func (r *RawConfig) ToConfig() (*Config, error) {
	colors, err := r.getConfigColors()
	if err != nil {
		return nil, err
	}

	// Check if default color is in the map
	var defaultColor Colors.Color
	if color, ok := colors["default"]; ok {
		defaultColor = color
	} else {
		return nil, errors.New("default color not found")
	}

	programGroups := r.getProgramGroups()

	timedProcessTriggers, err := r.getTimedProcessTriggers(&colors, programGroups)
	if err != nil {
		return nil, err
	}

	processTriggers, err := r.getProcessTriggers(&colors, programGroups)
	if err != nil {
		return nil, err
	}

	triggers := make([]ColorTrigger, 0, len(timedProcessTriggers)+len(processTriggers))
	// Timed process triggers is before since they have higher priority
	triggers = append(triggers, timedProcessTriggers...)
	triggers = append(triggers, processTriggers...)

	config := &Config{
		TargetLightName: strings.ToLower(r.TargetLightName),
		DefaultColor:    defaultColor,
		Triggers:        triggers,
	}

	err = setHueCredentials(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func (r *RawConfig) getMatchingNames(name string, programGroups *map[string][]string) []string {
	if strings.HasPrefix(name, "g:") {
		if names, ok := (*programGroups)[strings.ToLower(strings.TrimPrefix(name, "g:"))]; ok {
			return names
		}
	}

	return []string{strings.ToLower(name)}
}

func (r *RawConfig) getTimedProcessTriggers(colors *map[string]Colors.Color, groups map[string][]string) ([]ColorTrigger, error) {
	timedPrograms := make([]ColorTrigger, 0, len(r.TimedPrograms))
	for _, timedProgram := range r.TimedPrograms {
		var color Colors.Color
		color, err := parseColor(timedProgram.Color, colors)
		if err != nil {
			return nil, err
		}

		timeSpan, err := r.parseDayTimeSpan(timedProgram)
		if err != nil {
			return nil, err
		}

		names := r.getMatchingNames(timedProgram.Name, &groups)
		trigger := Triggers.NewTimedProcess(names, timeSpan)
		colorTrigger := NewColorTrigger(trigger, color)
		timedPrograms = append(timedPrograms, colorTrigger)
	}
	return timedPrograms, nil
}

func (r *RawConfig) parseDayTimeSpan(timedProgram RawTimedProgram) (Dates.IDayTimeSpan, error) {
	var timeSpan Dates.IDayTimeSpan
	timeSpan, ok := Dates.ParseStringToSpan(timedProgram.TimeSpan)
	if ok {
		return timeSpan, nil
	}
	timeSpan, ok = Dates.ParseStringToSpanSolar(timedProgram.TimeSpan)
	if ok {
		return timeSpan, nil
	}

	return nil, errors.New("invalid time span")
}

func (r *RawConfig) getProcessTriggers(colors *map[string]Colors.Color, groups map[string][]string) ([]ColorTrigger, error) {
	processTriggers := make([]ColorTrigger, 0, len(r.Programs))
	for _, program := range r.Programs {
		color, err := parseColor(program.Color, colors)
		if err != nil {
			return nil, err
		}
		names := r.getMatchingNames(program.Name, &groups)
		trigger := Triggers.NewProcess(names)
		colorTrigger := NewColorTrigger(trigger, color)
		processTriggers = append(processTriggers, colorTrigger)
	}
	return processTriggers, nil
}

func (r *RawConfig) getProgramGroups() map[string][]string {
	programGroups := make(map[string][]string)
	for _, group := range r.ProgramGroups {
		// Lowercase all programs
		for i := range group.Programs {
			group.Programs[i] = strings.ToLower(group.Programs[i])
		}

		programGroups[strings.ToLower(group.Name)] = group.Programs
	}
	return programGroups
}

func (r *RawConfig) getConfigColors() (map[string]Colors.Color, error) {
	colors := make(map[string]Colors.Color)
	for name, color := range r.Colors {
		rgb, err := Colors.NewRGBFromHex(color)
		if err != nil {
			return nil, err
		}
		colors[name] = &rgb
	}
	return colors, nil
}
