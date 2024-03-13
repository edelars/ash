package aliases

import (
	"ash/internal/commands"
	"ash/internal/dto"
)

func NewAliasCommand(fileToExec, description, displayName string, executor executorIface, executeKey uint16, iContext dto.InternalContextIface) *commands.Command {
	return commands.NewCommandWithExtendedInfo(displayName,
		func(_ dto.InternalContextIface, _ []string) dto.ExecResult {
			iCtx := iContext.WithLastKeyPressed(executeKey).WithCurrentInputBuffer([]rune(fileToExec))
			return executor.Execute(iCtx)
		}, true, description, displayName)
}

type executorIface interface {
	Execute(internalC dto.InternalContextIface) dto.ExecResult
}
