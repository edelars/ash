package list

import (
	"ash/internal/commands"
	"ash/internal/dto"
)

func NewRemoveLeftSymbol(deleteLeftSymbolAndMoveCursor func()) *commands.Command {
	return commands.NewCommand(":RemoveLeftSymbol",
		func(_ dto.InternalContextIface) {
			deleteLeftSymbolAndMoveCursor()
		})
}
