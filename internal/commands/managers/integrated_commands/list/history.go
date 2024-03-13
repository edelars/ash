package list

import (
	"ash/internal/commands"
	"ash/internal/dto"
)

const (
	cmdNameHistory = "history"
	cmdDescHistory = "Shows command history"
)

func NewHistoryCommand(iContext dto.InternalContextIface) *commands.Command {
	return commands.NewCommandWithExtendedInfo(cmdNameHistory,
		func(iContext dto.InternalContextIface, _ []string) dto.ExecResult {
			panic("TODO history")
			return dto.CommandExecResultStatusOk
		}, true, cmdDescHistory, cmdNameHistory)
}
