package main

import (
	settings "etm/common/settings"
	utl "etm/common/utils"
	cfg "etm/configuration"
	def "etm/definition"
	"etm/executor"
	"fmt"

	"gopkg.in/yaml.v3"
)

func bootstrap() executor.Executor {
	config := settings.LoadSettings()
	machineDefinition, err := def.LoadMachineDefinition(config.General.TuringMachineConfigurationFile)
	utl.PanicOnError(err)

	machine, err := cfg.CreateMachineConfiguration(machineDefinition, config)
	utl.PanicOnError(err)

	return executor.New(machine, config.Executor.InitialTape, config.Executor.InitialIndex)
}

func main() {
	executor := bootstrap()

	printExecutor(&executor)
	executor.Run(nil, printExecutor)
}

func printExecutor(ex *executor.Executor) {
	fmt.Println(
		ex.ToDisplay().Tape,
		ex.ToDisplay().CurrentState.Name,
		ex.ToDisplay().Index,
	)
}

func printYaml(inter interface{}) {
	ym, _ := yaml.Marshal(inter)
	fmt.Println(string(ym))
}
