package executor

import (
	"encoding/json"
	"errors"
	"etm/common/utils"
	. "etm/configuration"
	"slices"
	"time"
)

type Executor struct {
	machine      MachineConfiguration
	currentState singleState
	tape         []string
	index        int64
	halted       bool
}

type singleState struct {
	Name  string
	State State
}

func InitializeExecutor(machine MachineConfiguration, initialTape []string, initialTapeIndex int64) Executor {
	executor := Executor{}
	executor.machine = machine
	executor.currentState = singleState{Name: machine.InitialState, State: machine.StateMap[machine.InitialState]}
	executor.tape = initialTape
	for i := range executor.tape {
		executor.tape[i] = executor.machine.BlankSymbol
	}
	executor.index = initialTapeIndex

	return executor
}

func (ex Executor) IsHalted() bool {
	return ex.halted
}

func (ex *Executor) Step() error {
	if !ex.halted {
		currentStateAction := ex.currentState.State[ex.tape[ex.index]]
		ex.currentState = singleState{
			Name:  currentStateAction.Transition,
			State: ex.machine.StateMap[currentStateAction.Transition],
		}
		if slices.Contains(
			utils.ConvertSingleOrArrayEitherToArray(ex.machine.PermittedFinalStates),
			ex.currentState.Name,
		) {
			ex.halted = true
			return nil
		}
		ex.tape[ex.index] = currentStateAction.Write
		move(ex, currentStateAction.Move)

		return nil
	} else {
		return errors.New("Turing machine has halted.")
	}

}

func (ex *Executor) Run(prestep func(*Executor), poststep func(*Executor)) {
	for !ex.halted {
		if prestep != nil {
			prestep(ex)
		}
		ex.Step()
		if poststep != nil {
			poststep(ex)
		}
		time.Sleep(time.Duration(ex.machine.ExecutionDelay) * time.Millisecond)
	}
}

func move(ex *Executor, move Move) {
	if move == MoveLeft {
		moveRight(ex)
	} else if move == MoveRight {
		moveLeft(ex)
	}
}

func moveRight(ex *Executor) {
	if ex.index == int64(len(ex.tape)-1) {
		ex.tape = append(ex.tape, ex.machine.BlankSymbol)
	} else {
		ex.index += 1
	}
}

func moveLeft(ex *Executor) {
	if ex.index == 0 {
		ex.tape = append(ex.tape, ex.machine.BlankSymbol)
		copy(ex.tape[1:], ex.tape)
		ex.tape[0] = ex.machine.BlankSymbol
	} else {
		ex.index -= 1
	}
}

func (ex *Executor) ToDisplay() executorDisplay {
	return executorDisplay{
		CurrentState: ex.currentState,
		Tape:         ex.tape,
		Index:        ex.index,
		Halted:       ex.halted,
	}
}

type executorDisplay struct {
	CurrentState singleState
	Tape         []string
	Index        int64
	Halted       bool
}

func (input Executor) MarshalJSON() ([]byte, error) {
	return json.Marshal(input.ToDisplay())
}

func (input Executor) MarshalYAML() (interface{}, error) {
	return input.ToDisplay(), nil
}
