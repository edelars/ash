package list

import (
	"ash/internal/commands"
	"ash/internal/dto"
)

func NewRemoveLeftSymbol(deleteLeftSymbolAndMoveCursor func(), DeleteLastSymbolFromCurrentBuffer func() error) *commands.Command {
	return commands.NewCommand(":RemoveLeftSymbol",
		func(_ dto.InternalContextIface) int {
			if err := DeleteLastSymbolFromCurrentBuffer(); err == nil {
				deleteLeftSymbolAndMoveCursor()
			}
			return -1
		},
		false)
}
