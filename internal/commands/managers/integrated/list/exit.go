package list

import (
	"context"

	"ash/internal/commands"
	"ash/internal/internal_context"
)

func NewExitCommand() *commands.Command {
	return commands.NewCommand("exit",
		func(ctx context.Context, internalContext internal_context.InternalContextIface, inputChan chan []byte, outputChan chan []byte) {
			panic("exit")
		})
}
