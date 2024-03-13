package list

import (
	"os"

	"ash/internal/commands"
	"ash/internal/configuration/envs_loader"
	"ash/internal/dto"
)

const (
	cmdNameExport = "export"
	cmdDescExport = "Set OS env"
)

func NewExportCommand() *commands.Command {
	return commands.NewCommandWithExtendedInfo(cmdNameExport,
		func(iContext dto.InternalContextIface, args []string) dto.ExecResult {
			el := iContext.GetExecutionList()
			if len(el) == 1 && len(args) == 1 {
				if a, b, err := envs_loader.ParseEnvString(args[0]); err == nil {
					os.Setenv(a, b)
				} else {
					iContext.GetPrintFunction()("Fail to set ENV: " + args[0] + "\n")
					return 1
				}
			}
			return dto.CommandExecResultStatusOk
		}, true, cmdDescExport, cmdNameExport)
}
