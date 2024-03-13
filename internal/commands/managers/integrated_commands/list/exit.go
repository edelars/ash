package list

import (
	"ash/internal/commands"
	"ash/internal/dto"
)

const (
	cmdNameExit = "exit"
	cmdDescExit = "Escape from ash"
)

func NewExitCommand() *commands.Command {
	return commands.NewCommandWithExtendedInfo(cmdNameExit,
		func(_ dto.InternalContextIface, _ []string) dto.ExecResult {
			return dto.CommandExecResultMainExit
		}, true, cmdDescExit, cmdNameExit)
}
