package file_system

import (
	"os/exec"

	"ash/internal/commands"
	"ash/internal/dto"
)

func NewSystemCommand(fileToExec string) *commands.Command {
	return commands.NewCommand(fileToExec,
		func(iContext dto.InternalContextIface, args []string) dto.ExecResult {
			ctx := iContext.GetCTX()
			cmd := exec.CommandContext(ctx, fileToExec, args...)
			cmd.Stdin = iContext.GetInputReader()
			cmd.Stderr = iContext.GetOutputWriter()
			cmd.Stdout = iContext.GetOutputWriter()
			if err := cmd.Run(); err != nil {
				// iContext.GetPrintFunction()(fmt.Sprintf("could not run command:%s :%s ", fileToExec, err.Error()))
			}
			// iContext.GetPrintFunction()(cmd.Path)
			return dto.ExecResult(cmd.ProcessState.ExitCode())
		}, true)
}
