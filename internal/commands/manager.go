package commands

import (
	"ash/internal/dto"
)

type commandManager struct {
	mainName string
	data     []dto.CommandIface
}

func (m commandManager) SearchCommands(resultChan chan dto.CommandManagerSearchResult, patterns ...dto.PatternIface) {
	for _, pattern := range patterns {
		var founded foundedData

		if pattern.IsPrecisionSearch() {
			founded = m.precisionSearchInCommands(pattern.GetPattern())
		} else {
			founded = m.searchPatternInCommands(pattern.GetPattern())
		}
		var arr []dto.CommandIface
		for c := range founded {
			arr = append(arr, c)
		}

		resultChan <- &searchResult{
			name:         m.mainName,
			commandsData: arr,
			patternValue: pattern,
		}
	}
}

func NewCommandManager(mainName string, cmds ...dto.CommandIface) (cm commandManager) {
	cm.data = append(cm.data, cmds...)
	cm.mainName = mainName
	return cm
}

type foundedData map[dto.CommandIface]int8

// Search pattern in commands:
// [cmd] - [pattern] = [result %]
// exit - ext = 75%
// gettalk = ttk = 40%
// cd - cd = 100%
// cd - vv = 0%
// cd - cc = 50%
func (m commandManager) searchPatternInCommands(searchPattern string) foundedData {
	searchPatternRunes := []rune(searchPattern)
	founded := make(foundedData)

	for _, cmd := range m.data {

		percentCorrect := int8(100)
		step := getStepValue(cmd.GetName())

		cmdRunes := []rune(cmd.GetName())
		var pointerPattern int
		var f bool

	patternLoop:
		for pointerCmd := 0; pointerCmd < len(cmdRunes); pointerCmd++ {
			for i := pointerPattern; i < len(searchPatternRunes); i++ {
				if searchPatternRunes[i] == cmdRunes[pointerCmd] {
					pointerPattern = i + 1
					f = true
					continue patternLoop
				}
			}
			percentCorrect = percentCorrect - step
		}

		if percentCorrect > 0 && f {
			founded[cmd] = percentCorrect
		}
	}
	return founded
}

func getStepValue(s string) int8 {
	runeCount := len([]rune(s))
	if runeCount == 0 {
		return 0
	}
	return int8(100 / runeCount)
}

func (m commandManager) precisionSearchInCommands(searchName string) foundedData {
	founded := make(foundedData)
	for _, cmd := range m.data {
		if cmd.GetName() == searchName {
			founded[cmd] = 100
		}
	}
	return founded
}

type searchResult struct {
	name         string
	commandsData []dto.CommandIface
	patternValue dto.PatternIface
}

func (searchresult *searchResult) GetSourceName() string {
	return searchresult.name
}

func (searchresult *searchResult) GetCommands() []dto.CommandIface {
	return searchresult.commandsData
}

func (searchresult *searchResult) GetPattern() dto.PatternIface {
	return searchresult.patternValue
}

func (searchresult *searchResult) Founded() int {
	return len(searchresult.commandsData)
}
