package list

import (
	"ash/internal/commands"
	"ash/internal/dto"
)

func NewExecuteCommand() *commands.Command {
	return commands.NewCommand(":Execute",
		func(internalC dto.InternalContextIface) {
			for _, cmd := range internalC.GetExecutionList() {
				cmd.GetExecFunc()(internalC)
				return // TODO
			}
		})
}
