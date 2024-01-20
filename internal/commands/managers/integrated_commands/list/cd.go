package list

import (
	"os"

	"ash/internal/commands"
	"ash/internal/dto"
)

func NewCDCommand() *commands.Command {
	return commands.NewCommand("cd",
		func(internalC dto.InternalContextIface) {
			el := internalC.GetExecutionList()
			if len(el) == 1 {
				os.Chdir(el[0].GetArgs())
			}
		})
}
