package internal_actions

import (
	"ash/internal/commands"
	"ash/internal/commands/managers/internal_actions/list"
	"ash/internal/configuration"
	"ash/internal/dto"
	"ash/internal/pseudo_graphics"
	"ash/internal/storage"
)

func NewInternalActionsManager(dr pseudo_graphics.Drawer, searchFunc func(iContext dto.InternalContextIface, pattern dto.PatternIface) []dto.CommandManagerSearchResult, inputSet func(r []rune), autocomplOpts configuration.AutocompleteOpts, historyAddFunc func(data storage.DataIface), colorsAdapter dto.ColorsAdapterIface) (im commands.CommandManagerIface) {
	return commands.NewCommandManager(
		"Actions",
		1,
		false,
		list.NewExecuteCommand(configuration.CmdExecute, historyAddFunc),
		list.NewAutocompleteCommand(configuration.CmdAutocomplete, dr, searchFunc, inputSet, autocomplOpts, colorsAdapter),
		list.NewAutocompleteCommand(configuration.CmdKeyUp, dr, searchFunc, inputSet, autocomplOpts, colorsAdapter),
	)
}
