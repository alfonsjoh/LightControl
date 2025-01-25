package Models

type AllLights map[string]Light

type Light struct {
	State struct {
		On        bool   `json:"on"`
		Bri       int    `json:"bri"`
		Ct        int    `json:"ct"`
		Alert     string `json:"alert"`
		Colormode string `json:"colormode"`
		Mode      string `json:"mode"`
		Reachable bool   `json:"reachable"`
	} `json:"state"`
	Swupdate struct {
		State       string `json:"state"`
		Lastinstall string `json:"lastinstall"`
	} `json:"swupdate"`
	Type             string `json:"type"`
	Name             string `json:"name"`
	Modelid          string `json:"modelid"`
	Manufacturername string `json:"manufacturername"`
	Productname      string `json:"productname"`
	Capabilities     struct {
		Certified bool `json:"certified"`
		Control   struct {
			Mindimlevel int `json:"mindimlevel"`
			Maxlumen    int `json:"maxlumen"`
			Ct          struct {
				Min int `json:"min"`
				Max int `json:"max"`
			} `json:"ct"`
		} `json:"control"`
		Streaming struct {
			Renderer bool `json:"renderer"`
			Proxy    bool `json:"proxy"`
		} `json:"streaming"`
	} `json:"capabilities"`
	Config struct {
		Archetype string `json:"archetype"`
		Function  string `json:"function"`
		Direction string `json:"direction"`
		Startup   struct {
			Mode       string `json:"mode"`
			Configured bool   `json:"configured"`
		} `json:"startup"`
	} `json:"config"`
	Uniqueid   string `json:"uniqueid"`
	Swversion  string `json:"swversion"`
	Swconfigid string `json:"swconfigid"`
	Productid  string `json:"productid"`
}
