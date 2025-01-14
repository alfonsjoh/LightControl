package Dates

import (
	"strings"
	"time"

	"github.com/nathan-osman/go-sunrise"
)

const Latitude = 57.689940
const Longitude = 11.973001

type SolarState int

const (
	SolarOn SolarState = iota
	SolarOff
)

type IDayTimeSpan interface {
	Contains(t DayTime) bool
}

type DayTimeSpan struct {
	Start DayTime
	End   DayTime
}

type DayTimeSpanSolar struct {
	SolarState SolarState
}

func NewDayTimeSpan(start, end DayTime) DayTimeSpan {
	return DayTimeSpan{start, end}
}

func ParseStringToSpan(s string) (IDayTimeSpan, bool) {
	// Parse string with format "HH:MM:SS-HH:MM:SS"
	split := strings.Split(s, "-")
	if len(split) != 2 {
		return nil, false
	}

	start, ok1 := ParseString(split[0])
	end, ok2 := ParseString(split[1])

	if !ok1 || !ok2 {
		return nil, false
	}

	return &DayTimeSpan{start, end}, true
}

func ParseStringToSpanSolar(s string) (IDayTimeSpan, bool) {
	if s == "solarOn" {
		return &DayTimeSpanSolar{SolarOn}, true
	}
	if s == "solarOff" {
		return &DayTimeSpanSolar{SolarOff}, true
	}
	return nil, false
}

func (d *DayTimeSpan) Contains(t DayTime) bool {
	return d.Start.Diff(t) <= 0 && d.End.Diff(t) >= 0
}

func (d *DayTimeSpanSolar) Contains(t DayTime) bool {
	now := time.Now()
	rise, set := sunrise.SunriseSunset(Latitude, Longitude, now.Year(), now.Month(), now.Day())
	timeSpan := NewDayTimeSpan(
		NewDayTime(rise.Hour(), rise.Minute(), rise.Second()),
		NewDayTime(set.Hour(), set.Minute(), set.Second()),
	)
	switch d.SolarState {
	case SolarOn:
		return timeSpan.Contains(t)
	case SolarOff:
		return !timeSpan.Contains(t)
	}
	return false
}
