package list

import (
	"ash/internal/commands"
	"ash/internal/dto"
)

func NewEchoCommand() *commands.Command {
	return commands.NewCommand("echo",
		func(internalC dto.InternalContextIface, args []string) dto.ExecResult {
			el := internalC.GetExecutionList()
			if len(el) == 1 && len(args) == 1 {
				internalC.GetPrintFunction()(args[0] + "\n")
			}
			return dto.CommandExecResultStatusOk
		}, true)
}
