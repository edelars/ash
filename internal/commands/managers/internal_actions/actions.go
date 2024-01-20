package internal_actions

import (
	"ash/internal/commands"
	"ash/internal/commands/managers/internal_actions/list"
)

func NewInternalAcgionsManager() (im commands.CommandManagerIface) {
	return commands.NewCommandManager(list.NewExecuteCommand())
}
