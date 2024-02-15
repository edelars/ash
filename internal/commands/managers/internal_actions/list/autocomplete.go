package list

import (
	"ash/internal/commands"
	"ash/internal/configuration"
	"ash/internal/data_source"
	"ash/internal/dto"
	"ash/internal/pseudo_graphics"
	"ash/internal/pseudo_graphics/windows/selection_window"
)

func NewAutocompleteCommand(cmdName string, dr pseudo_graphics.Drawer, searchFunc func(iContext dto.InternalContextIface, pattern dto.PatternIface) []dto.CommandManagerSearchResult, setInputFunc func(r []rune), autocomplOpts configuration.AutocompleteOpts) *commands.Command {
	return commands.NewCommand(cmdName,
		func(iContext dto.InternalContextIface, _ []string) dto.ExecResult {
			doneChan := make(chan struct{}, 1)
			defer close(doneChan)

			cmdChan := make(chan dto.CommandIface, 1)
			defer close(cmdChan)

			var r []rune

			sFunc := func(pattern []rune) dto.DataSource {
				p := commands.NewPattern(string(pattern), false)
				return data_source.NewDataSource(searchFunc(iContext, p))
			}
			rFunc := func(cmd dto.CommandIface, newUserInput []rune) {
				r = newUserInput
				doneChan <- struct{}{}
				cmdChan <- cmd
			}

			pWindow := selection_window.NewSelectionWindow(iContext.GetCurrentInputBuffer(), sFunc, rFunc, autocomplOpts)
			dr.Draw(&pWindow, iContext, doneChan)
			if cmd := <-cmdChan; cmd != nil {
				return cmd.GetExecFunc()(iContext, nil)
			}
			setInputFunc(r)
			return dto.CommandExecResultNewUserInput
		}, false)
}
