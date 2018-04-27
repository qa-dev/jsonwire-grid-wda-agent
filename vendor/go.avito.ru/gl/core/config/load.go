package config

import (
	"encoding/json"
	"errors"
	"os"
)

func LoadFromFile(configPath string, config interface{}) error {
	if configPath == "" {
		return errors.New("empty configuration file path")
	}

	configFile, err := os.Open(configPath)
	if err != nil {
		return err
	}

	defer configFile.Close()

	jsonParser := json.NewDecoder(configFile)

	return jsonParser.Decode(&config)
}
