package configuration

import (
	"gopkg.in/yaml.v3"
	"io"
	"os"
)

type Config struct {
	DeviceConfig DeviceConfig `yaml:"device_config"`
	Auth         Auth         `yaml:"auth"`
}

type Auth struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type DeviceConfig struct {
	Local  bool   `yaml:"local"`
	Host   string `yaml:"host"`
	Port   int    `yaml:"port"`
	CPU0ID uint8  `yaml:"cpu0_id"`
	CPU1ID uint8  `yaml:"cpu1_id"`
}

// ReadConfigFile read config from file
func ReadConfigFile(configPath string) (Config, error) {
	config := Config{}

	data, err := os.Open(configPath)
	if err != nil {
		return Config{}, err
	}

	configFile, err := io.ReadAll(data)
	if err != nil {
		return Config{}, err
	}
	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}
