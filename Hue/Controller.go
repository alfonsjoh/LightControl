package Hue

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"LightControl/Hue/Models"
)

type Controller struct {
	config Config
}

func NewController(config Config) *Controller {
	return &Controller{config}
}

func (c *Controller) GetAllLights() (Models.AllLights, error) {
	resp, err := http.Get(fmt.Sprintf("http://%s/api/%s/lights", c.config.Address, c.config.ID))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	//data, err := io.ReadAll(resp.Body)
	//if err != nil {
	//        return nil, err
	//}
	//
	//fmt.Println(string(data))

	var lights Models.AllLights
	if err := json.NewDecoder(resp.Body).Decode(&lights); err != nil {
		return nil, err
	}

	return lights, nil
}

func (c *Controller) SetColor(lightID string, color Color) error {
	colorString, err := color.GetColor()
	if err != nil {
		return err
	}
	requestBody := bytes.NewBufferString(colorString)

	url := fmt.Sprintf("http://%s/api/%s/lights/%s/state", c.config.Address, c.config.ID, lightID)
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
