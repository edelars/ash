package file_system

import (
	"context"
	"io"
	"os"
	"os/exec"

	"ash/internal/commands"
	"ash/internal/dto"

	"ash/pkg/termbox"

	"github.com/creack/pty"
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
			// cmd.Stdin = iContext.GetInputReader()
			// cmd.Stderr = iContext.GetOutputWriter()
			// cmd.Stdout = iContext.GetOutputWriter()
			cmd.Env = os.Environ()

			// termbox.SetInputMode(termbox.InputEsc)
			// ptmx, err := pty.Start(cmd)
			ptmx, err := pty.StartWithSize(cmd, &pty.Winsize{Rows: 167, Cols: 47})
			if err != nil {
				panic(err)
			}

			defer func() { _ = ptmx.Close() }() // Best effort.

			pty.InheritSize(os.Stdin, ptmx)
			go func() { _, _ = io.Copy(ptmx, iContext.GetInputReader()) }()
			_, _ = io.Copy(iContext.GetOutputWriter(), ptmx)

			// pty.Setsize(ptmx, &pty.Winsize{Rows: 10, Cols: 10})
			// M
			// cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true, Foreground: false}

			// if err := cmd.Run(); err != nil {
			// 	iContext.GetPrintFunction()(fmt.Sprintf("command: %s :%s\n", fileToExec, err.Error()))
			// }
			if err := termbox.ReInit(); err != nil {
				panic(err)
			}
			// termbox.SetInputMode(termbox.InputMouse)

			return dto.ExecResult(cmd.ProcessState.ExitCode())
		}, true, description, "")
}
