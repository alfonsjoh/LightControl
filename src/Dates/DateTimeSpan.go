package Dates

import (
	"fmt"
)

type SolarState int

const (
	SolarOn SolarState = iota
	SolarOff
)

type DayTimeSpan struct {
	Start DayTime
	End   DayTime
}

func (d *DayTimeSpan) ToString() string {
	return fmt.Sprintf("%s-%s", d.Start.ToString(), d.End.ToString())
}

func NewDayTimeSpan(start, end DayTime) DayTimeSpan {
	return DayTimeSpan{start, end}
}

func (d *DayTimeSpan) Contains(t DayTime) bool {
	return d.Start.Diff(t) <= 0 && d.End.Diff(t) >= 0
}
