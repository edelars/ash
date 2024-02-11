package file_system

import (
	"ash/internal/commands"
	"ash/internal/dto"
)

func NewPseudoCommand(r string, setInputFunc func(r []rune)) *commands.Command {
	return commands.NewCommand(r,
		func(_ dto.InternalContextIface, _ []string) dto.ExecResult {
			setInputFunc([]rune(r))
			return dto.CommandExecResultNewUserInput
		}, false)
}
