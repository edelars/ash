package file_system

import (
	"ash/internal/commands"
)

func NewPseudoCommand(fileToExec string) *commands.Command {
	return commands.NewCommand(fileToExec,
		nil, true)
}
