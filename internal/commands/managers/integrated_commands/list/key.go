package list

import (
	"fmt"

	"ash/internal/commands"
	"ash/internal/dto"

	"github.com/nsf/termbox-go"
)

func NewKeyCommand() *commands.Command {
	return commands.NewCommand("_key",
		func(internalC dto.InternalContextIface) {
			internalC.GetPrintFunction()("press Enter (13) key to break")
			for {
				ev := <-internalC.GetInputEventChan()
				switch ev.Type {
				case termbox.EventKey:
					if ev.Ch != 0 {
						internalC.GetPrintFunction()(fmt.Sprintf("got key: %d", ev.Ch))
					} else {
						switch ev.Ch {
						case 13:
							return
						default:
							internalC.GetPrintFunction()(fmt.Sprintf("got key: %d", ev.Key))
						}
					}
				}

			}
		})
}
