package Triggers

import (
	"strings"
)

type Process struct {
	Names []string // All names that match this process
}

func NewProcess(names []string) *Process {
	return &Process{
		Names: names,
	}
}

func (program *Process) IsEnabled(p *TriggerState) bool {
	for _, processName := range p.ProcessNames {
		for _, name := range program.Names {
			if strings.Contains(processName, name) {
				return true
			}
		}
	}
	return false
}
