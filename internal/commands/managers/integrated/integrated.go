package integrated

import (
	"ash/internal/commands"
	"ash/internal/commands/managers/integrated/list"
)

type IntergatedManager struct {
	data []commands.CommandIface
}

func (m IntergatedManager) SearchCommands(resultChan chan commands.CommandManagerSearchResult, patterns ...commands.PatternIface) {
	founded := make(map[commands.CommandIface]struct{})
	for _, pattern := range patterns {
		for _, r := range pattern.GetPattern() {
			founded = m.searchRuneInCommands(r, founded)
		}
		arr := make([]commands.CommandIface, len(founded))
		for c := range founded {
			arr = append(arr, c)
		}
		resultChan <- &searchResult{
			name:         "internal",
			commandsData: arr,
			patternValue: pattern,
		}
	}
}

func NewIntegratedManager() (im IntergatedManager) {
	im.data = append(im.data, list.NewExitCommand())
	return im
}

type foundedData map[commands.CommandIface]struct{}

func (m IntergatedManager) searchRuneInCommands(r rune, f foundedData) foundedData {
	return f
}

type searchResult struct {
	name         string
	commandsData []commands.CommandIface
	patternValue commands.PatternIface
}

func (searchresult *searchResult) GetSourceName() string {
	return searchresult.name
}

func (searchresult *searchResult) GetCommands() []commands.CommandIface {
	return searchresult.commandsData
}

func (searchresult *searchResult) GetPattern() commands.PatternIface {
	return searchresult.patternValue
}
