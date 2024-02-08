package list

import (
	"ash/internal/commands"
	"ash/internal/dto"
)

func NewLogoutCommand() *commands.Command {
	return commands.NewCommand("logout",
		func(_ dto.InternalContextIface, _ []string) dto.ExecResult {
			return dto.CommandExecResultMainExit
		}, true)
}
