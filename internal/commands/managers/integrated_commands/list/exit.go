package list

import (
	"ash/internal/commands"
	"ash/internal/dto"
)

func NewExitCommand() *commands.Command {
	return commands.NewCommandWithExtendedInfo("exit",
		func(iContext dto.InternalContextIface, _ []string) dto.ExecResult {
			return dto.CommandExecResultMainExit
		}, true, "Exit from ash", "exit")
}
