package main

import (
	settings "etm/common/settings"
	utl "etm/common/utils"
	cfg "etm/configuration"
	def "etm/definition"
	"etm/executor"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func bootstrap() executor.Executor {
	config := settings.LoadSettings()
	machineDefinition, err := def.LoadMachineDefinition(config.General.TuringMachineConfigurationFile)
	utl.PanicOnError(err)

	machine, err := cfg.New(machineDefinition, config)
	utl.PanicOnError(err)

	return executor.New(machine)
}

func main() {
	executor := bootstrap()

	config := settings.LoadSettings()
	var file *os.File
	if config.General.LogFile != "" {
		file, _ = os.Create(config.General.LogFile)
		defer file.Close()
	} else {
		file = nil
	}

	pe := printExecutor(file)
	pe(&executor)
	err := executor.Run(nil, pe)
	fmt.Print("\n", err)
}

func printExecutor(file *os.File) func(*executor.Executor) {
	return func(exc *executor.Executor) {
		dispEx := exc.ToDisplay()
		fmt.Print(
			"\033[2K\r",
			dispEx.Tape, " ",
			dispEx.CurrentState.Name, " ",
			dispEx.Index,
		)
		if file != nil {
			file.Write([]byte(fmt.Sprintf("%v %v %v\n", dispEx.Tape, dispEx.CurrentState.Name, dispEx.Index)))
		}
	}
}

func printYaml(inter interface{}) {
	ym, _ := yaml.Marshal(inter)
	fmt.Println(string(ym))
}
