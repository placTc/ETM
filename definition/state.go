package stateDefinition

import "etm/common/utils"

type StateDefinition struct {
	Initial  string                           `yaml:"initial"`
	Final    utils.Either[string, []string]   `yaml:"final"`
	NullMove bool                             `yaml:"null_move"`
	States   map[string]SingleStateDefinition `yaml:"states"`
}

type SingleStateDefinition map[string]StateActionDefinition

type StateActionDefinition struct {
	Write      string `yaml:"write,omitempty"`
	Move       string `yaml:"move,omitempty"`
	Transition string `yaml:"transition,omitempty"`
}
