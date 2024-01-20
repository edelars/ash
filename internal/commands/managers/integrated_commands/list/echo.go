package list

import (
	"os"

	"ash/internal/commands"
	"ash/internal/dto"
)

func NewEchoCommand() *commands.Command {
	return commands.NewCommand("echo",
		func(internalC dto.InternalContextIface) {
			el := internalC.GetExecutionList()
			if len(el) == 1 {
				os.Chdir(el[0].GetArgs())
			}
		})
}
