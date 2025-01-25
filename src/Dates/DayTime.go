package Dates

import (
	"errors"
	"strconv"
	"strings"
)

type DayTime struct {
	Hour   int
	Minute int
	Second int
}

func NewDayTime(hour int, minute int, second int) DayTime {
	return DayTime{hour, minute, second}
}

func (d *DayTime) Diff(other DayTime) int {
	return (d.Hour-other.Hour)*3600 + (d.Minute-other.Minute)*60 + d.Second - other.Second
}

func ParseString(s string) (DayTime, bool) {
	// Parse string with format "HH:MM:SS"
	split := strings.Split(s, ":")
	if !(len(split) == 3 || len(split) == 2) {
		return DayTime{}, false
	}

	hour, err1 := strconv.Atoi(split[0])
	minute, err2 := strconv.Atoi(split[1])
	second := 0
	var err3 error
	if len(split) == 3 {
		second, err3 = strconv.Atoi(split[2])
	}

	err := errors.Join(err1, err2, err3)

	if err != nil {
		return DayTime{}, false
	}

	if hour < 0 || hour > 23 || minute < 0 || minute > 59 || second < 0 || second > 59 {
		return DayTime{}, false
	}

	return NewDayTime(hour, minute, second), true
}

func (d *DayTime) ToString() string {
	return strconv.Itoa(d.Hour) + ":" + strconv.Itoa(d.Minute) + ":" + strconv.Itoa(d.Second)
}
