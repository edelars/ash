package list

import (
	"fmt"

	"ash/internal/commands"
	"ash/internal/dto"
)

func NewKeyCommand() *commands.Command {
	return commands.NewCommand("_key",
		func(internalC dto.InternalContextIface) {
			internalC.GetPrintFunction()("press Enter (13) key to break")
			for {
				key := <-internalC.GetInputChan()
				switch key {
				case 13:
					return
				default:
					internalC.GetPrintFunction()(fmt.Sprintf("got key: %d", int(key)))
				}
			}
		})
}
