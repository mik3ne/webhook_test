package config

import "go.uber.org/config"

const CONFIG_PATH = "config/config.yml"

type Configuration struct {
	URL      string
	Requests struct {
		Amount    int
		PerSecond int `yaml:"per_second"`
	}
}

func LoadConfig() (Configuration, error) {
	var c Configuration

	cfg, err := config.NewYAML(config.File(CONFIG_PATH))
	if err != nil {
		return c, err
	}

	if err := cfg.Get("").Populate(&c); err != nil {
		return c, err
	}

	return c, nil
}
