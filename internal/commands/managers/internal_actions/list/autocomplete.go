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

			cmdChan := make(chan dto.CommandIface, 1)
			defer close(cmdChan)

			sFunc := func(pattern []rune) dto.DataSource {
				p := commands.NewPattern(string(pattern), false)
				return data_source.NewDataSource(searchFunc(iContext, p))
			}
			rFunc := func(cmd dto.CommandIface, newUserInput []rune) {
				doneChan <- struct{}{}
				cmdChan <- cmd
			}

			pWindow := selection_window.NewSelectionWindow(iContext.GetCurrentInputBuffer(), sFunc, rFunc)
			dr.Draw(&pWindow, iContext, doneChan)
			if cmd := <-cmdChan; cmd != nil {
				cmd.GetExecFunc()(iContext, nil)
				return 0
			}
			return -1
		}, false)
}
