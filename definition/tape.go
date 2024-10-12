package stateDefinition

type TapeDefinition struct {
	InitialTape  []string `yaml:"initial_tape"`
	InitialIndex int64    `yaml:"initial_index"`
}
