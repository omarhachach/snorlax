package snorlax

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

var cachedConfigPath string

// Config is the configuration struct for the Snorlax bot.
type Config struct {
	// AutoDelete determines whether or not to automatically delete command
	// messages.
	AutoDelete    bool     `json:"autoDelete"`
	DBPath        string   `json:"dbPath"`
	Debug         bool     `json:"debug"`
	DisplayAuthor bool     `json:"displayAuthor"`
	Token         string   `json:"token"`
	Owners        []string `json:"owners"`
}

// UpdateFile will update the config.json file.
func (c *Config) UpdateFile() error {
	newFile, err := json.MarshalIndent(*c, "", "    ")
	if err != nil {
		return err
	}

	file, err := os.OpenFile(cachedConfigPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	file.Write(newFile)
	return nil
}

// ParseConfig parses a config file, and returns a new Config.
func ParseConfig(path string) (*Config, error) {
	cachedConfigPath = path
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
