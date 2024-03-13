package list

import (
	"ash/internal/commands"
	"ash/internal/dto"
)

const (
	cmdNameLogout = "logout"
)

func NewLogoutCommand() *commands.Command {
	return commands.NewCommandWithExtendedInfo(cmdNameLogout,
		func(_ dto.InternalContextIface, _ []string) dto.ExecResult {
			return dto.CommandExecResultMainExit
		}, true, cmdDescExit, cmdNameLogout)
}
