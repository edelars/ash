package list

import (
	"os"

	"ash/internal/commands"
	"ash/internal/configuration/envs_loader"
	"ash/internal/dto"
)

func NewExportCommand() *commands.Command {
	return commands.NewCommand("export",
		func(internalC dto.InternalContextIface) int {
			el := internalC.GetExecutionList()
			if len(el) == 1 {
				if a, b, err := envs_loader.ParseEnvString(el[0].GetArgs()); err == nil {
					os.Setenv(a, b)
				} else {
					internalC.GetPrintFunction()("Fail to set ENV: " + el[0].GetArgs())
				}
			}
			return 0
		}, true)
}
