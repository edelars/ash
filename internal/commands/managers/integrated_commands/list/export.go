package list

import (
	"os"

	"ash/internal/commands"
	"ash/internal/configuration/envs_loader"
	"ash/internal/dto"
)

func NewExportCommand() *commands.Command {
	return commands.NewCommand("export",
		func(internalC dto.InternalContextIface, args []string) dto.ExecResult {
			el := internalC.GetExecutionList()
			if len(el) == 1 && len(args) == 1 {
				if a, b, err := envs_loader.ParseEnvString(args[0]); err == nil {
					os.Setenv(a, b)
				} else {
					internalC.GetPrintFunction()("Fail to set ENV: " + args[0])
				}
			}
			return dto.CommandExecResultStatusOk
		}, true)
}
