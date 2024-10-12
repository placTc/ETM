package structure

import (
	. "etm/common"
	. "etm/common/settings"
	. "etm/common/utils"
	. "etm/definition"
	"fmt"
	"slices"
)

const BlankSymbolReference = "blank"
const DefaultReference = "default"

type MachineConfiguration struct {
	ExecutionDelay int
	Symbols        []string
	BlankSymbol    string
	PermittedInput Either[string, []string]
	InitialState   string
	HaltingStates  Either[string, []string]
	StateMap       map[string]State
	Tape           TapeDefinition
}

type State map[string]StateAction

type StateAction struct {
	Write      string
	Move       Move
	Transition string
}

func New(definition MachineDefinition, settings Settings) (MachineConfiguration, error) {
	machineConfiguration := MachineConfiguration{}
	machineConfiguration.ExecutionDelay = settings.Executor.ExecutionDelayMs
	if machineConfiguration.ExecutionDelay < 0 {
		return MachineConfiguration{}, fmt.Errorf(
			"Could not create Turing Machine configuration, "+
				"the execution delay was negative: %v",
			settings.Executor.ExecutionDelayMs,
		)
	}

	if definition.Alphabet.Symbols == nil || len(definition.Alphabet.Symbols) == 0 {
		return MachineConfiguration{}, fmt.Errorf(
			"Could not create Turing Machine configuration, symbol alphabet was not specified",
		)
	}
	machineConfiguration.Symbols = definition.Alphabet.Symbols

	if slices.Contains(machineConfiguration.Symbols, DefaultReference) {
		return MachineConfiguration{}, fmt.Errorf(
			"Could not create Turing Machine configuration, '%v' is a protected keyword and cannot be set as a symbol",
			DefaultReference,
		)
	}

	if slices.Contains(machineConfiguration.Symbols, BlankSymbolReference) {
		return MachineConfiguration{}, fmt.Errorf(
			"Could not create Turing Machine configuration, '%v' is a protected keyword and cannot be set as a symbol",
			BlankSymbolReference,
		)
	}

	if definition.Alphabet.Blank == "" {
		return MachineConfiguration{}, fmt.Errorf(
			"Could not create Turing Machine configuration, blank symbol was not specified",
		)
	}
	machineConfiguration.BlankSymbol = definition.Alphabet.Blank

	if !slices.Contains(machineConfiguration.Symbols, machineConfiguration.BlankSymbol) {
		return MachineConfiguration{}, fmt.Errorf(
			"Could not create Turing Machine configuration, "+
				"the blank symbol '%v' was not found in the alphabet '%v'",
			definition.Alphabet.Blank,
			definition.Alphabet.Symbols,
		)
	}

	inputs := ConvertSingleOrArrayEitherToArray(definition.Alphabet.Input)
	if inputs[0] == "" && len(inputs) == 1 || len(inputs) == 0 {
		return MachineConfiguration{}, fmt.Errorf(
			"Could not create Turing Machine configuration, input symbols were not defined",
		)
	}

	for i := 0; i < len(inputs); i++ {
		if !slices.Contains(machineConfiguration.Symbols, inputs[i]) {
			return MachineConfiguration{}, fmt.Errorf(
				"Could not create Turing Machine configuration, "+
					"the input symbol '%v' was not found in the alphabet '%v'",
				inputs[i],
				inputs,
				definition.Alphabet.Symbols,
			)
		}
	}
	machineConfiguration.PermittedInput = definition.Alphabet.Input

	for i := range definition.Tape.InitialTape {
		if definition.Tape.InitialTape[i] == BlankSymbolReference {
			definition.Tape.InitialTape[i] = machineConfiguration.BlankSymbol
		}
		if !slices.Contains(machineConfiguration.Symbols, definition.Tape.InitialTape[i]) {
			return MachineConfiguration{}, fmt.Errorf(
				"Could not create Turing Machine configuration, "+
					"initial tape confgiuration contained symbol '%v' not present in alphabet '%v'",
				definition.Tape.InitialTape[i],
				machineConfiguration.Symbols,
			)
		}
	}
	machineConfiguration.Tape.InitialTape = definition.Tape.InitialTape

	if definition.Tape.InitialIndex < 0 || definition.Tape.InitialIndex > int64(len(machineConfiguration.Tape.InitialTape)-1) {
		return MachineConfiguration{}, fmt.Errorf(
			"Could not create Turing Machine configuration, "+
				"initial tape index %v was outside the range of the initial tape %v",
			definition.Tape.InitialIndex,
			machineConfiguration.Tape.InitialTape,
		)
	}
	machineConfiguration.Tape.InitialIndex = definition.Tape.InitialIndex

	var stateMap map[string]State = map[string]State{}
	for stateName, stateDefinition := range definition.StateDefinition.States {
		var state State = State{}
		for stateActionSymbol, stateActionDefinition := range stateDefinition {
			if stateActionSymbol == BlankSymbolReference {
				stateActionSymbol = definition.Alphabet.Blank
			}

			if stateActionDefinition.Transition == "" {
				stateActionDefinition.Transition = stateName
			}

			if stateActionDefinition.Write == BlankSymbolReference {
				stateActionDefinition.Write = machineConfiguration.BlankSymbol
			}

			if definition.StateDefinition.NullMove && stateActionDefinition.Move == "" {
				stateActionDefinition.Move = string(NullMove)
			}

			if !slices.Contains(machineConfiguration.Symbols, stateActionSymbol) && stateActionSymbol != DefaultReference {
				return MachineConfiguration{}, fmt.Errorf(
					"Could not create Turing Machine configuration, "+
						"the action symbol '%v' for state '%v' was not found in the alphabet '%v'",
					stateActionSymbol,
					stateName,
					machineConfiguration.Symbols,
				)
			}
			_, exists := definition.StateDefinition.States[stateActionDefinition.Transition]

			if !slices.Contains(machineConfiguration.Symbols, stateActionDefinition.Write) &&
				stateActionSymbol != DefaultReference && stateActionDefinition.Write != "" {
				return MachineConfiguration{}, fmt.Errorf(
					"Could not create Turing Machine configuration, "+
						"the write symbol '%v' for action symbol '%v' in state '%v' was not found in the alphabet '%v'",
					stateActionDefinition.Write,
					stateActionSymbol,
					stateName,
					machineConfiguration.Symbols,
				)
			}

			if !slices.Contains(
				append(
					ConvertSingleOrArrayEitherToArray(machineConfiguration.PermittedInput),
					machineConfiguration.BlankSymbol,
				),
				stateActionDefinition.Write,
			) && stateActionSymbol != DefaultReference && stateActionDefinition.Write != "" {
				return MachineConfiguration{}, fmt.Errorf(
					"Could not create Turing Machine configuration, "+
						"the write symbol '%v' for action symbol '%v' in state '%v' was not in the permitted inputs '%v'",
					stateActionDefinition.Write,
					stateActionSymbol,
					stateName,
					ConvertSingleOrArrayEitherToArray(machineConfiguration.PermittedInput),
				)
			}

			if !slices.Contains([]string{string(MoveRight), string(MoveLeft), string(NullMove)}, stateActionDefinition.Move) {
				return MachineConfiguration{}, fmt.Errorf(
					"Could not create Turing Machine configuration, "+
						"the specified move '%v' for action symbol '%v' in state '%v' was not L, R or N",
					stateActionDefinition.Move,
					stateActionSymbol,
					stateName,
				)
			}

			if !exists {
				return MachineConfiguration{}, fmt.Errorf(
					"Could not create Turing Machine configuration, "+
						"the specified transition state '%v' for state '%v' was not found in the list of state definitions",
					stateActionDefinition.Transition,
					stateName,
				)
			}

			if !definition.StateDefinition.NullMove && stateActionDefinition.Move == "" {
				return MachineConfiguration{}, fmt.Errorf(
					"Could not create Turing Machine configuration, "+
						"null move specified in state '%v' for symbol '%v', while null moves are disallowed",
					stateName,
					stateActionSymbol,
				)
			}

			state[stateActionSymbol] = StateAction{
				Move:       Move(stateActionDefinition.Move),
				Transition: stateActionDefinition.Transition,
				Write:      stateActionDefinition.Write,
			}
		}
		for i := range machineConfiguration.Symbols {
			_, exists := state[machineConfiguration.Symbols[i]]
			if !exists {
				state[machineConfiguration.Symbols[i]] = StateAction{
					Move:       state[DefaultReference].Move,
					Transition: state[DefaultReference].Transition,
					Write:      machineConfiguration.Symbols[i],
				}
			}
		}
		delete(state, DefaultReference)
		stateMap[stateName] = state
	}
	machineConfiguration.StateMap = stateMap

	if definition.StateDefinition.Initial == "" {
		return MachineConfiguration{}, fmt.Errorf(
			"Could not create Turing Machine configuration, the initial state was not defined",
		)
	}

	_, exists := definition.StateDefinition.States[definition.StateDefinition.Initial]
	if !exists {
		return MachineConfiguration{}, fmt.Errorf(
			"Could not create Turing Machine configuration, "+
				"the specified initial state %v was not found in the list of state definitions",
			definition.StateDefinition.Initial,
		)
	}
	machineConfiguration.InitialState = definition.StateDefinition.Initial

	if definition.StateDefinition.Halting.IsNil() {
		return MachineConfiguration{}, fmt.Errorf(
			"Could not create Turing Machine configuration, the final state was not defined",
		)
	}

	finalStates := ConvertSingleOrArrayEitherToArray(definition.StateDefinition.Halting)
	for i := 0; i < len(finalStates); i++ {
		_, exists = definition.StateDefinition.States[finalStates[i]]
		if !exists {
			return MachineConfiguration{}, fmt.Errorf(
				"Could not create Turing Machine configuration, "+
					"the specified final state '%v' was not found in the list of state definitions",
				finalStates[i],
			)
		}
	}
	machineConfiguration.HaltingStates = definition.StateDefinition.Halting

	return machineConfiguration, nil
}
