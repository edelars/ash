package list

import (
	"ash/internal/commands"
	"ash/internal/dto"
)

const (
	cmdNameAlias = "alias"
	cmdDescAlias = "Set alias for the short command call"
)

func NewAliasCommand(iContext dto.InternalContextIface) *commands.Command {
	return commands.NewCommandWithExtendedInfo(cmdNameAlias,
		func(iContext dto.InternalContextIface, _ []string) dto.ExecResult {
			panic("TODO alias")
			// internalC.GetPrintFunction()("we`re done")
			// internalC.GetErrChan() <- errors.New("ash exiting")
			return dto.CommandExecResultStatusOk
		}, true, cmdDescAlias, cmdNameAlias)
}
