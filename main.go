package main

import (
	cfg "etm/common/config"
	"etm/common/utils"
	def "etm/definition"
	ex "etm/executor"
	mc "etm/machineConfiguration"
	"fmt"

	"gopkg.in/yaml.v3"
)

func bootstrap() ex.Executor {
	config := cfg.LoadConfig()
	machineDefinition, err := def.LoadMachineDefinition(config.General.TuringMachineConfigurationFile)
	utils.PanicOnError(err)

	machine, err := mc.CreateMachineConfiguration(machineDefinition, config)
	utils.PanicOnError(err)

	return ex.InitializeExecutor(machine, config.Executor.InitialTape, config.Executor.InitialIndex)
}

func main() {
	executor := bootstrap()

	printExecutor(&executor)
	executor.Run(nil, printExecutor)
}

func printExecutor(ex *ex.Executor) {
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
