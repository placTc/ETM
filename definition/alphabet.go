package stateDefinition

import "etm/common/utils"

type AlphabetDefinition struct {
	Symbols []string                       `yaml:"symbols"`
	Blank   string                         `yaml:"blank"`
	Input   utils.Either[string, []string] `yaml:"input"`
}
