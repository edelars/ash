package list

import (
	"ash/internal/commands"
	"ash/internal/dto"
)

func NewExecuteCommand() *commands.Command {
	return commands.NewCommand(":Execute",
		func(internalC dto.InternalContextIface) int {
			for _, cmd := range internalC.GetExecutionList() {
				cmd.GetExecFunc()(internalC)
			}
			return 0 // TODO
		}, true)
}
