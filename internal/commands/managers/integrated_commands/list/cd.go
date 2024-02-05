package list

import (
	"os"

	"ash/internal/commands"
	"ash/internal/dto"
)

func NewCDCommand() *commands.Command {
	return commands.NewCommand("cd",
		func(internalC dto.InternalContextIface, args []string) int {
			el := internalC.GetExecutionList()
			if len(el) == 1 && len(args) == 1 {
				os.Chdir(args[0])
			}
			return 0
		}, true)
}
