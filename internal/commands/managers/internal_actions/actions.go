package internal_actions

import (
	"ash/internal/commands"
	"ash/internal/commands/managers/internal_actions/list"
	"ash/internal/configuration"
	"ash/internal/dto"
	"ash/internal/pseudo_graphics"
)

func NewInternalActionsManager(dr pseudo_graphics.Drawer, searchFunc func(iContext dto.InternalContextIface, pattern dto.PatternIface) []dto.CommandManagerSearchResult, inputSet func(r []rune), autocomplOpts configuration.AutocompleteOpts) (im commands.CommandManagerIface) {
	return commands.NewCommandManager(
		"Actions",
		1,
		list.NewExecuteCommand(),
		list.NewAutocompleteCommand(dr, searchFunc, inputSet, autocomplOpts),
	)
}
