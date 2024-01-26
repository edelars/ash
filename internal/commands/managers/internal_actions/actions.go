package internal_actions

import (
	"ash/internal/commands"
	"ash/internal/commands/managers/internal_actions/list"
	"ash/internal/dto"
	"ash/internal/pseudo_graphics"
)

func NewInternalActionsManager(dr pseudo_graphics.Drawer, searchFunc func(pattern dto.PatternIface) []dto.CommandManagerSearchResult) (im commands.CommandManagerIface) {
	return commands.NewCommandManager(
		"actions",
		list.NewExecuteCommand(),
		list.NewAutocompleteCommand(dr, searchFunc),
	)
}
