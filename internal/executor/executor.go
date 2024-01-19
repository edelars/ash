package executor

import (
	"ash/internal/dto"
)

type CommandExecutor struct {
	commandRouter     routerIface
	keyBindingManager keyBindingsIface
}

func NewCommandExecutor(commandRouter routerIface, keyBindingManager keyBindingsIface) CommandExecutor {
	return CommandExecutor{
		commandRouter:     commandRouter,
		keyBindingManager: keyBindingManager,
	}
}

func (r CommandExecutor) Execute(internalC dto.InternalContextIface) {
	// mainCommand := r.keyBindingManager.GetCommandByKey(int(internalC.GetLastKeyPressed()))
	internalC.GetOutputChan() <- byte(44)
}

type routerIface interface {
	SearchCommands(patterns ...dto.PatternIface) dto.CommandRouterSearchResult
}

type keyBindingsIface interface {
	GetCommandByKey(key int) dto.CommandIface
}
