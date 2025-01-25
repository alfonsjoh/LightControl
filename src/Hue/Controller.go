package Hue

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"LightControl/src/Extensions"
	"LightControl/src/Hue/Colors"
	"LightControl/src/Hue/Models"
)

var previousLightStates = make(map[string]string)

type Controller struct {
	configLock *Extensions.StructLock[Config]
}

func NewController(configLock *Extensions.StructLock[Config]) *Controller {
	return &Controller{configLock}
}

func (c *Controller) GetAllLights() (Models.AllLights, error) {
	config := c.configLock.Get()
	resp, err := http.Get(fmt.Sprintf("http://%s/api/%s/lights", config.Address, config.ID))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var lights Models.AllLights
	if err := json.NewDecoder(resp.Body).Decode(&lights); err != nil {
		return nil, err
	}

	return lights, nil
}

func (c *Controller) SetColor(lightID string, color Colors.Color) error {
	config := c.configLock.Get()
	colorString, err := color.GetColor()
	if prev, ok := previousLightStates[lightID]; ok {
		if prev == colorString {
			return nil
		}
	}
	previousLightStates[lightID] = colorString

	fmt.Println("Setting light state:", colorString)
	if err != nil {
		return err
	}
	requestBody := bytes.NewBufferString(colorString)

	url := fmt.Sprintf("http://%s/api/%s/lights/%s/state", config.Address, config.ID, lightID)
	request, err := http.NewRequest("PUT", url, requestBody)
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (c *Controller) SetOnOff(lightID string, on bool) error {
	config := c.configLock.Get()
	requestBody := bytes.NewBufferString(fmt.Sprintf("{\"on\": %t}", on))

	url := fmt.Sprintf("http://%s/api/%s/lights/%s/state", config.Address, config.ID, lightID)
	request, err := http.NewRequest("PUT", url, requestBody)
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
