package list

import (
	"errors"

	"ash/internal/commands"
	"ash/internal/dto"
)

func NewExitCommand() *commands.Command {
	return commands.NewCommand("exit",
		func(iContext dto.InternalContextIface, _ []string) int {
			iContext.GetPrintFunction()("we`re done")
			iContext.GetErrChan() <- errors.New("ash exiting")
			return 0
		}, true)
}
