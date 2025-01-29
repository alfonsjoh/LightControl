package Triggers

import "LightControl/src/Dates"

type TriggerState struct {
	ProcessNames []string
	CurrentTime  Dates.DayTime
}

func NewTriggerState(processNames []string, currentTime Dates.DayTime) *TriggerState {
	return &TriggerState{
		ProcessNames: processNames,
		CurrentTime:  currentTime,
	}
}
