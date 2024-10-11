package config

import (
	u "etm/common/utils"
	"os"

	"gopkg.in/yaml.v3"
)

const configFilename = "config.yaml"

type Config struct {
	General  General  `yaml:"general"`
	Executor Executor `yaml:"executor"`
}

type General struct {
	TuringMachineConfigurationFile string `yaml:"turingMachineConfiguration"`
}

type Executor struct {
	ExecutionDelayMs int      `yaml:"executionDelayMs"`
	InitialTape      []string `yaml:"initialTape"`
	InitialIndex     int64    `yaml:"initialIndex"`
}

func LoadConfig() Config {
	config := Config{}

	fileData, err := os.ReadFile(configFilename)
	u.PanicOnErrorWithCustomMessage(err, "Couldn't read config file!")

	err = yaml.Unmarshal(fileData, &config)
	u.PanicOnErrorWithCustomMessage(err, "Couldn't parse config.yaml file!")

	return config
}
