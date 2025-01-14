package Hue

type Config struct {
	ID      string
	Address string
}

func NewConfig(id, address string) Config {
	return Config{id, address}
}
