package list

import (
	"errors"

	"ash/internal/commands"
	"ash/internal/dto"
)

func NewExecuteCommand() *commands.Command {
	return commands.NewCommand(":Execute",
		func(internalC dto.InternalContextIface) {
			internalC.GetErrChan() <- errors.New("ash exiting")
		})
}
