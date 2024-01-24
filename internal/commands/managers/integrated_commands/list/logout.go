package list

import (
	"errors"

	"ash/internal/commands"
	"ash/internal/dto"
)

func NewLogoutCommand() *commands.Command {
	return commands.NewCommand("logout",
		func(internalC dto.InternalContextIface) {
			internalC.GetErrChan() <- errors.New("ash exiting")
		})
}
