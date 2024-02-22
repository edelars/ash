package list

import (
	"fmt"

	"ash/internal/commands"
	"ash/internal/dto"

	"ash/pkg/termbox"
)

const (
	cmdNameKey = "_key"
	cmdDescKey = "Displays input key number"
)

func NewKeyCommand() *commands.Command {
	return commands.NewCommandWithExtendedInfo(cmdNameKey,
		func(iContext dto.InternalContextIface, _ []string) dto.ExecResult {
			iContext.GetPrintFunction()("press Enter (13) key to break\n")
			for {
				ev := <-iContext.GetInputEventChan()
				switch ev.Type {
				case termbox.EventKey:
					if ev.Ch != 0 {
						iContext.GetPrintFunction()(fmt.Sprintf("got ch key: %d\n", ev.Ch))
					} else {
						switch ev.Key {
						case 13:
							return dto.CommandExecResultStatusOk
						default:
							iContext.GetPrintFunction()(fmt.Sprintf("got key: %d, mod: %d\n", ev.Key, ev.Mod))
						}
					}
				}

			}
		}, true, cmdDescKey, cmdNameKey)
}
