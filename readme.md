
```
 _     _       _     _   _____             _             _ 
| |   (_)     | |   | | /  __ \           | |           | |
| |    _  __ _| |__ | |_| /  \/ ___  _ __ | |_ _ __ ___ | |
| |   | |/ _` | '_ \| __| |    / _ \| '_ \| __| '__/ _ \| |
| |___| | (_| | | | | |_| \__/\ (_) | | | | |_| | | (_) | |
\_____/_|\__, |_| |_|\__|\____/\___/|_| |_|\__|_|  \___/|_|
    __/ |                                            
    |___/
```

This project connects to a philips hue bridge and controls a chosen lamps colors based on open applications on your pc.


## Limitations
* This is only compatible with windows.
* This is only compatible with philips hue lights.
* This can for now only change the color of one light.
* The configuration file only supports setting colors based on open applications and timespans.
  Timespans can be based on sunrise and sunset, but for now coordinates are hardcoded.

## Configuration
#### Configuration file
The configuration file is located in the root of the project and is called `config.json`. The configuration file is a json file with the following structure:

```json
{
  "target_light_name": "{Lamp name}",
  "colors": {
    "default": "{default color}",
    "{optional color}": "color value"
  },
  "program_groups": [
    {
      "name": "{group name}",
      "programs": ["Application name 1", "..."]
    }
  ],
  "programs": [
    { "name": "{application name or group}", "color": "{color}" }
  ],
  "timed_programs": [
    { "name": "{application name or group}", "span": "{solorOn/solarOff}", "color": "{color}" }
  ]
}
```

Colors in "programs" or "timed_programs" are either hex values without the '#' or a color in the colors object. (ex. "default").
A group in timed_programs starts with "g:" and then the group name.
For example if the group name is "dev" then you would write "g:dev".

#### Credentials
Place your credentials to your hue bridge in an environment variable called `HueCredentials`.
It should be in the format `{Hue Id}@{Hue address}`.

## Running
To run the program go needs to be installed on the system. First run `go mod tidy`'
then you can build the program either with `make build` or simply with
`go build -ldflags -H=windowsgui -o ./out/LightControl.exe ./src/`. Then the program will be located at `./out/LightControl.exe`.