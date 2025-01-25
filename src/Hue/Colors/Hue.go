package Colors

import (
	"encoding/json"
)

type Hue struct {
	Hue        int  `json:"hue"`
	Saturation int  `json:"sat"`
	Brightness int  `json:"bri"`
	On         bool `json:"on"`
}

func (c *Hue) GetColor() (string, error) {
	c.setOn()
	bytes, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (c *Hue) setOn() {
	c.On = c.Brightness > 0
}
