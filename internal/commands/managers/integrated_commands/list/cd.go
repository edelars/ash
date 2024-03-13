package list

import (
	"os"

	"ash/internal/commands"
	"ash/internal/dto"
)

const (
	cmdNameCd = "cd"
	cmdDescCd = "Change current directory"
)

func NewCDCommand() *commands.Command {
	return commands.NewCommandWithExtendedInfo(cmdNameCd,
		func(iContext dto.InternalContextIface, args []string) dto.ExecResult {
			el := iContext.GetExecutionList()
			if len(el) == 1 && len(args) == 1 {
				os.Chdir(args[0])
			}
			return dto.CommandExecResultStatusOk
		}, true, cmdDescCd, cmdNameCd)
}
