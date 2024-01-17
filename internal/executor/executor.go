package executor

import (
	"ash/internal/dto"
)

type CommandExecutor struct {
	commandRouter routerIface
}

func NewCommandExecutor(commandRouter routerIface) CommandExecutor {
	return CommandExecutor{
		commandRouter: commandRouter,
	}
}

func (r CommandExecutor) Execute(internalC dto.InternalContextIface) {
	internalC.GetOutputChan() <- byte(44)
}

type routerIface interface {
	SearchCommands(patterns ...dto.PatternIface) dto.CommandRouterSearchResult
}
