package snorlax

import (
	"encoding/json"
	"io/ioutil"
)

// Config is the configuration struct for the Snorlax bot.
type Config struct {
	// AutoDelete determines whether or not to automatically delete command
	// messages.
	AutoDelete bool     `json:"autoDelete"`
	Debug      bool     `json:"debug"`
	Token      string   `json:"token"`
	Owners     []string `json:"owners"`
}

// ParseConfig parses a config file, and returns a new Config.
func ParseConfig(path string) (*Config, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	config := &Config{}
	err = json.Unmarshal(bytes, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
