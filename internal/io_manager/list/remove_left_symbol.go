package list

import (
	"ash/internal/commands"
	"ash/internal/dto"
)

func NewRemoveLeftSymbol(cmdName string, deleteLeftSymbolAndMoveCursor func(), DeleteLastSymbolFromCurrentBuffer func() error) *commands.Command {
	return commands.NewCommand(cmdName,
		func(_ dto.InternalContextIface, _ []string) dto.ExecResult {
			if err := DeleteLastSymbolFromCurrentBuffer(); err == nil {
				deleteLeftSymbolAndMoveCursor()
			}
			return dto.CommandExecResultNotDoAnyting
		},
		false)
}
