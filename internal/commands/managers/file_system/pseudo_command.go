package file_system

import (
	"ash/internal/commands"
	"ash/internal/dto"
)

func NewPseudoCommand(r string, setInputFunc func(r []rune), fileInfo, displayName string) *commands.Command {
	return commands.NewCommandWithExtendedInfo(r,
		func(_ dto.InternalContextIface, _ []string) dto.ExecResult {
			setInputFunc([]rune(r))
			return dto.CommandExecResultNewUserInput
		}, false, fileInfo, displayName)
}
