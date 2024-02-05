package list

import (
	"errors"

	"ash/internal/commands"
	"ash/internal/dto"
)

func NewAliasCommand(iContext dto.InternalContextIface) *commands.Command {
	return commands.NewCommand("alias",
		func(internalC dto.InternalContextIface, _ []string) int {
			panic("TODO alias")
			internalC.GetPrintFunction()("we`re done")
			internalC.GetErrChan() <- errors.New("ash exiting")
			return 0
		}, true)
}
