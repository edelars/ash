package list

import (
	"ash/internal/commands"
	"ash/internal/dto"
)

func NewEchoCommand() *commands.Command {
	return commands.NewCommand("echo",
		func(internalC dto.InternalContextIface, args []string) int {
			el := internalC.GetExecutionList()
			if len(el) == 1 && len(args) == 1 {
				internalC.GetPrintFunction()(args[0] + "\n")
			}
			return 0
		}, true)
}
