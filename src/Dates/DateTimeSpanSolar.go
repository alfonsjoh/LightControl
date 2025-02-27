package Dates

import (
	"time"

	"github.com/nathan-osman/go-sunrise"
)

const Latitude = 57.689940
const Longitude = 11.973001

type DayTimeSpanSolar struct {
	SolarState SolarState
}

func (d *DayTimeSpanSolar) ToString() string {
	switch d.SolarState {
	case SolarOn:
		return "solarOn"
	case SolarOff:
		return "solarOff"
	}
	return ""
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

func (d *DayTimeSpanSolar) Contains(t DayTime) bool {
	now := time.Now()
	rise, set := sunrise.SunriseSunset(Latitude, Longitude, now.Year(), now.Month(), now.Day())
	// rise and set are in UTC, convert to local time
	rise = rise.Local()
	set = set.Local()

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
