package list

import (
	"ash/internal/commands"
	"ash/internal/data_source"
	"ash/internal/dto"
	"ash/internal/pseudo_graphics"
	"ash/internal/pseudo_graphics/windows/selection_window"
)

func NewAutocompleteCommand(dr pseudo_graphics.Drawer, searchFunc func(pattern dto.PatternIface) []dto.CommandManagerSearchResult) *commands.Command {
	return commands.NewCommand(":Autocomplete",
		func(internalC dto.InternalContextIface) int {
			doneChan := make(chan struct{}, 1)
			defer close(doneChan)

			sFunc := func(pattern []rune) dto.DataSource {
				p := commands.NewPattern(string(pattern), false)
				return data_source.NewDataSource(searchFunc(p))
			}
			rFunc := func(cmd dto.CommandIface, userInput []rune) {
				doneChan <- struct{}{}
				internalC.GetPrintFunction()(" selected: " + cmd.GetName())
			}

			pWindow := selection_window.NewSelectionWindow(internalC.GetCurrentInputBuffer(), sFunc, rFunc)
			dr.Draw(&pWindow, internalC, doneChan)
			return -1
		}, false)
}
