package file_system

import (
	"context"
	"fmt"
	"os/exec"

	"ash/internal/commands"
	"ash/internal/dto"

	"ash/pkg/termbox"
)

func NewSystemCommand(fileToExec, description string) *commands.Command {
	return commands.NewCommandWithExtendedInfo(fileToExec,
		func(iContext dto.InternalContextIface, args []string) dto.ExecResult {
			ctx, cancelFunc := context.WithCancel(context.Background())
			defer cancelFunc()

			go func() {
				select {
				case <-iContext.GetExecTerminateChan():
					cancelFunc()
				case <-ctx.Done():
					return
				}
			}()

			cmd := exec.CommandContext(ctx, fileToExec, args...)
			cmd.Stdin = iContext.GetInputReader()
			cmd.Stderr = iContext.GetOutputWriter()
			cmd.Stdout = iContext.GetOutputWriter()

			if err := cmd.Run(); err != nil {
				iContext.GetPrintFunction()(fmt.Sprintf("command: %s :%s\n", fileToExec, err.Error()))
			}
			if err := termbox.ReInit(); err != nil {
				panic(err)
			}

			return dto.ExecResult(cmd.ProcessState.ExitCode())
		}, true, description, "")
}
