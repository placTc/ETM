package executor

import (
	"encoding/json"
	"errors"
	. "etm/common"
	"etm/common/utils"
	cfg "etm/configuration"
	"fmt"
	"slices"
	"time"
)

type Executor struct {
	machine      cfg.MachineConfiguration
	currentState singleState
	tape         []string
	index        int64
	halted       bool
}

type singleState struct {
	Name  string
	State cfg.State
}

func New(machine cfg.MachineConfiguration) Executor {
	executor := Executor{}
	executor.machine = machine
	executor.currentState = singleState{Name: machine.InitialState, State: machine.StateMap[machine.InitialState]}
	executor.tape = machine.Tape.InitialTape
	executor.index = machine.Tape.InitialIndex

	return executor
}

func (ex Executor) IsHalted() bool {
	return ex.halted
}

func (ex *Executor) Step() error {
	if !ex.halted {
		if currentStateIsHaltState(ex) {
			ex.halted = true
			return errors.New("Turing machine has halted.")
		}

		currentStateAction, exists := ex.currentState.State[ex.tape[ex.index]]
		if !exists {
			ex.halted = true
			return fmt.Errorf("Got unexpected symbol '%v' for state '%v'. Turing machine halted.", ex.tape[ex.index], ex.currentState.Name)
		}

		if currentStateAction.Write != "" {
			ex.tape[ex.index] = currentStateAction.Write
		}
		move(ex, currentStateAction.Move)
		ex.currentState = singleState{
			Name:  currentStateAction.Transition,
			State: ex.machine.StateMap[currentStateAction.Transition],
		}

		if ex.currentState.State == nil {

		}

		return nil
	} else {
		return errors.New("Turing machine has halted.")
	}

}

func currentStateIsHaltState(ex *Executor) bool {
	return slices.Contains(
		utils.ConvertSingleOrArrayEitherToArray(ex.machine.HaltingStates),
		ex.currentState.Name,
	)
}

func (ex *Executor) Run(prestep func(*Executor), poststep func(*Executor)) error {
	for !ex.halted {
		if prestep != nil {
			prestep(ex)
		}
		err := ex.Step()
		if err != nil {
			return err
		}
		if poststep != nil {
			poststep(ex)
		}
		time.Sleep(time.Duration(ex.machine.ExecutionDelay) * time.Millisecond)
	}
	
	return nil
}

func move(ex *Executor, move Move) {
	if move == MoveLeft {
		moveLeft(ex)
	} else if move == MoveRight {
		moveRight(ex)
	}
}

func moveRight(ex *Executor) {
	if ex.index == int64(len(ex.tape)-1) {
		ex.tape = append(ex.tape, ex.machine.BlankSymbol)
	}
	ex.index += 1
}

func moveLeft(ex *Executor) {
	if ex.index == 0 {
		ex.tape = append(ex.tape, "")
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
