package Config

import (
	"LightControl/src/Hue/Colors"
	"LightControl/src/Triggers"
)

type ColorTrigger struct {
	Triggers.Trigger
	Color Colors.Color
}

func NewColorTrigger(trigger Triggers.Trigger, color Colors.Color) ColorTrigger {
	return ColorTrigger{
		Trigger: trigger,
		Color:   color,
	}
}

func (trigger *ColorTrigger) GetIfEnabled(state *Triggers.TriggerState) (Colors.Color, bool) {
	if !trigger.Trigger.IsEnabled(state) {
		return nil, false
	}
	return trigger.Color, true
}
