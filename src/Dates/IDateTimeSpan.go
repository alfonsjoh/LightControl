package Dates

import "strings"

type IDayTimeSpan interface {
	Contains(t DayTime) bool
	ToString() string
}

func TimeSpanEquals(t1 IDayTimeSpan, t2 IDayTimeSpan) bool {
	return t1.ToString() == t2.ToString()
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
