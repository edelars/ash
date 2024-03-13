package commands

import (
	"ash/internal/dto"
)

func NewPseudoCommand(r string, setInputFunc func(r []rune), description, displayName string) *Command {
	return NewCommandWithExtendedInfo(r,
		func(_ dto.InternalContextIface, _ []string) dto.ExecResult {
			setInputFunc([]rune(r))
			return dto.CommandExecResultNewUserInput
		}, false, description, displayName)
}
