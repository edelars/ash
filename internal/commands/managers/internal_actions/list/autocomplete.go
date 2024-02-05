package list

import (
	"ash/internal/commands"
	"ash/internal/data_source"
	"ash/internal/dto"
	"ash/internal/pseudo_graphics"
	"ash/internal/pseudo_graphics/windows/selection_window"
)

func NewAutocompleteCommand(dr pseudo_graphics.Drawer, searchFunc func(iContext dto.InternalContextIface, pattern dto.PatternIface) []dto.CommandManagerSearchResult) *commands.Command {
	return commands.NewCommand(":Autocomplete",
		func(iContext dto.InternalContextIface, _ []string) int {
			doneChan := make(chan struct{}, 1)
			defer close(doneChan)

			sFunc := func(pattern []rune) dto.DataSource {
				p := commands.NewPattern(string(pattern), false)
				return data_source.NewDataSource(searchFunc(iContext, p))
			}
			rFunc := func(cmd dto.CommandIface, _ []rune) {
				doneChan <- struct{}{}
				iContext.GetPrintFunction()(" selected: " + cmd.GetName())
			}

			pWindow := selection_window.NewSelectionWindow(iContext.GetCurrentInputBuffer(), sFunc, rFunc)
			dr.Draw(&pWindow, iContext, doneChan)
			return -1
		}, false)
}
