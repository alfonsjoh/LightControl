package Colors

import (
	"errors"

	"LightControl/src/Dates"
)

type TimedRGBColor struct {
	timings map[Dates.DayTimeSpan]Color
}

func (c *TimedRGBColor) GetColor() (string, error) {
	dayTime := Dates.DayTimeNow()
	for span, color := range c.timings {
		if span.Contains(dayTime) {
			return color.GetColor()
		}
	}

	return "", errors.New("no color found for current time")
}
