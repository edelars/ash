package list

import (
	"ash/internal/commands"
	"ash/internal/dto"
)

func NewAliasCommand(iContext dto.InternalContextIface) *commands.Command {
	return commands.NewCommand("alias",
		func(internalC dto.InternalContextIface, _ []string) dto.ExecResult {
			panic("TODO alias")
			// internalC.GetPrintFunction()("we`re done")
			// internalC.GetErrChan() <- errors.New("ash exiting")
			return dto.CommandExecResultStatusOk
		}, true)
}
