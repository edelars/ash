package list

import (
	"fmt"

	"ash/internal/commands"
	"ash/internal/dto"

	"github.com/nsf/termbox-go"
)

func NewKeyCommand() *commands.Command {
	return commands.NewCommand("_key",
		func(internalC dto.InternalContextIface, _ []string) dto.ExecResult {
			internalC.GetPrintFunction()("press Enter (13) key to break\n")
			for {
				ev := <-internalC.GetInputEventChan()
				switch ev.Type {
				case termbox.EventKey:
					if ev.Ch != 0 {
						internalC.GetPrintFunction()(fmt.Sprintf("got key: %d\n", ev.Ch))
					} else {
						switch ev.Ch {
						case 13:
							return dto.CommandExecResultStatusOk
						default:
							internalC.GetPrintFunction()(fmt.Sprintf("got key: %d\n", ev.Key))
						}
					}
				}

			}
		}, true)
}
