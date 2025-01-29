package Triggers

import (
	"strings"

	"LightControl/src/Dates"
)

type TimedProcess struct {
	Names []string // All names that match this process
	Span  Dates.IDayTimeSpan
}

func NewTimedProcess(names []string, span Dates.IDayTimeSpan) *TimedProcess {
	return &TimedProcess{
		Names: names,
		Span:  span,
	}
}

func (program *TimedProcess) IsEnabled(p *TriggerState) bool {
	if !program.Span.Contains(p.CurrentTime) {
		return false
	}

	for _, processName := range p.ProcessNames {
		for _, name := range program.Names {
			if strings.Contains(name, processName) {
				return true
			}
		}
	}
	return false
}
