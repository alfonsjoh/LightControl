package Colors

import (
	"errors"
	"time"

	Dates2 "LightControl/src/Dates"
)

type TimedRGBColor struct {
	timings map[Dates2.DayTimeSpan]Color
}

func (c *TimedRGBColor) GetColor() (string, error) {
	now := time.Now()
	dayTime := Dates2.NewDayTime(now.Hour(), now.Minute(), now.Second())
	for span, color := range c.timings {
		if span.Contains(dayTime) {
			return color.GetColor()
		}
	}

	return "", errors.New("no color found for current time")
}
