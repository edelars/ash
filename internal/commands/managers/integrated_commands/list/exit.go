package list

import (
	"errors"

	"ash/internal/commands"
	"ash/internal/dto"
)

func NewExitCommand() *commands.Command {
	return commands.NewCommand("exit",
		func(internalC dto.InternalContextIface) int {
			internalC.GetPrintFunction()("we`re done")
			internalC.GetErrChan() <- errors.New("ash exiting")
			return 0
		}, true)
}
