package Hue

import (
	"encoding/json"
	"errors"
	"math"
	"strconv"
	"time"

	"LightControl/Dates"
)

type Color interface {
	GetColor() (string, error)
}

type HueColor struct {
	Hue        int `json:"hue"`
	Saturation int `json:"sat"`
	Brightness int `json:"bri"`
}

type RGBColor struct {
	Red   int
	Green int
	Blue  int
}

type TimedRGBColor struct {
	timings map[Dates.DayTimeSpan]Color
}

func (c *HueColor) GetColor() (string, error) {
	bytes, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (c *RGBColor) GetColor() (string, error) {
	hueColor := c.ToHueColor()
	return hueColor.GetColor()
}

func (c *TimedRGBColor) GetColor() (string, error) {
	now := time.Now()
	dayTime := Dates.NewDayTime(now.Hour(), now.Minute(), now.Second())
	for span, color := range c.timings {
		if span.Contains(dayTime) {
			return color.GetColor()
		}
	}

	return "", errors.New("no color found for current time")
}

func NewRGB(r int, g int, b int) RGBColor {
	return RGBColor{r, g, b}
}

func NewRGBFromHex(hex string) (RGBColor, error) {
	if len(hex) == 6 {
		r, err1 := strconv.ParseInt(hex[0:2], 16, 16)
		g, err2 := strconv.ParseInt(hex[2:4], 16, 16)
		b, err3 := strconv.ParseInt(hex[4:6], 16, 16)
		err := errors.Join(err1, err2, err3)
		if err != nil {
			return RGBColor{}, err
		}

		return NewRGB(int(r), int(g), int(b)), nil
	}

	return RGBColor{}, errors.New("invalid hex string length")
}

func (c *RGBColor) ToHueColor() HueColor {
	r, g, b := float64(c.Red)/255.0, float64(c.Green)/255.0, float64(c.Blue)/255.0
	max, min := math.Max(r, math.Max(g, b)), math.Min(r, math.Min(g, b))
	delta := max - min

	// Calculate Hue
	hue := 0.0
	switch {
	case delta == 0:
		hue = 0
	case max == r:
		hue = math.Mod((g-b)/delta, 6)
	case max == g:
		hue = (b-r)/delta + 2
	case max == b:
		hue = (r-g)/delta + 4
	}
	hue = math.Mod(hue*60+360, 360) // Ensure hue is non-negative
	hue = hue * (65535 / 360.0)     // Scale hue to range 0–65535

	// Calculate Saturation and Brightness
	saturation := 0.0
	if max != 0 {
		saturation = (delta / max) * 254
	}
	brightness := max * 254

	return HueColor{
		Hue:        int(hue),
		Saturation: int(saturation),
		Brightness: int(brightness),
	}
}
