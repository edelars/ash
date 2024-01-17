package list

import (
	"context"

	"ash/internal/commands"
	"ash/internal/dto"
)

func NewExitCommand() *commands.Command {
	return commands.NewCommand("exit",
		func(ctx context.Context, internalContext dto.InternalContextIface, inputChan chan []byte, outputChan chan []byte) {
			panic("exit")
		})
}
