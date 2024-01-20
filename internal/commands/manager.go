package commands

import (
	"ash/internal/dto"
)

type commandManager struct {
	data []dto.CommandIface
}

func (m commandManager) SearchCommands(resultChan chan dto.CommandManagerSearchResult, patterns ...dto.PatternIface) {
	founded := make(map[dto.CommandIface]int8)
	for _, pattern := range patterns {
		if pattern.IsPrecisionSearch() {
			founded = m.precisionSearchInCommands(pattern.GetPattern(), founded)
		} else {
			founded = m.searchPatternInCommands(pattern.GetPattern(), founded)
		}
		var arr []dto.CommandIface
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

func NewCommandManager(cmds ...dto.CommandIface) (cm commandManager) {
	cm.data = append(cm.data, cmds...)
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
func (m commandManager) searchPatternInCommands(searchPattern string, founded foundedData) foundedData {
	for _, cmd := range m.data {
		percentCorrect := int8(100)
		step := getStepValue(cmd.GetName())
		searchPatternRunes := []rune(searchPattern)

		var counterPattern int

		for _, cmdRune := range cmd.GetName() {
			if counterPattern >= len(searchPatternRunes) || cmdRune != searchPatternRunes[counterPattern] {
				percentCorrect = percentCorrect - step
			} else {
				counterPattern++
			}
		}
		if percentCorrect > 0 {
			founded[cmd] = int8(percentCorrect)
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

func (m commandManager) precisionSearchInCommands(searchName string, founded foundedData) foundedData {
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
