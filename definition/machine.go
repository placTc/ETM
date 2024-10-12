package stateDefinition

import (
	"os"

	"gopkg.in/yaml.v3"
)

type MachineDefinition struct {
	Alphabet        AlphabetDefinition `yaml:"alphabet"`
	Tape            TapeDefinition     `yaml:"tape"`
	StateDefinition StateDefinition    `yaml:"state"`
}

func LoadMachineDefinition(filepath string) (machine MachineDefinition, err error) {
	machine = MachineDefinition{}

	data, err := os.ReadFile(filepath)
	if err != nil {
		return machine, err
	}

	err = yaml.Unmarshal(data, &machine)
	return machine, err
}
