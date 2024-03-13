package commands

import (
	"ash/internal/dto"
)

type CommandManager struct {
	supportEmptyPattern bool
	priority            uint8 // for autocomplete sort
	mainName            string
	data                []dto.CommandIface
}

func (m CommandManager) SearchCommands(_ dto.InternalContextIface, resultChan chan dto.CommandManagerSearchResult, patterns ...dto.PatternIface) {
	for _, pattern := range patterns {
		var founded foundedData
		var arr []dto.CommandIface

		if pattern.GetPattern() != "" || (pattern.GetPattern() == "" && m.supportEmptyPattern) {

			if pattern.IsPrecisionSearch() {
				founded = m.precisionSearchInCommands(pattern.GetPattern())
			} else {
				founded = m.searchPatternInCommands(pattern.GetPattern())
			}
			for c, p := range founded {
				c.SetMathWeight(p + c.GetMathWeight()) // plus basic init weight
				arr = append(arr, c)
			}
		}

		resultChan <- &searchResult{
			name:         m.mainName,
			commandsData: arr,
			patternValue: pattern,
			priority:     m.priority,
		}
	}
}

func (m CommandManager) AddCommands(cmds ...dto.CommandIface) {
}

func NewCommandManager(mainName string, priority uint8, supportEmptyPattern bool, cmds ...dto.CommandIface) (cm CommandManager) {
	cm.data = append(cm.data, cmds...)
	cm.mainName = mainName
	cm.priority = priority
	cm.supportEmptyPattern = supportEmptyPattern
	return cm
}

type foundedData map[dto.CommandIface]uint8

// Search pattern in commands:
// [cmd] - [pattern] = [result %]
// exit - ext = 75%
// gettalk = ttk = 40%
// cd - cd = 100%
// cd - vv = 0%
// cd - cc = 50%
func (m CommandManager) searchPatternInCommands(searchPattern string) foundedData {
	searchPatternRunes := []rune(searchPattern)
	founded := make(foundedData)

	for _, cmd := range m.data {

		if searchPattern == "" && m.supportEmptyPattern {
			founded[cmd] = 100
			continue
		}

		percentCorrect := uint8(100)
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
			if percentCorrect > step {
				percentCorrect = percentCorrect - step
			} else {
				percentCorrect = 0
			}
		}

		if percentCorrect > 0 && f {
			founded[cmd] = percentCorrect
		}

	}
	return founded
}

func getStepValue(s string) uint8 {
	runeCount := len([]rune(s))
	if runeCount == 0 {
		return 0
	}
	return uint8(100 / runeCount)
}

func (m CommandManager) precisionSearchInCommands(searchName string) foundedData {
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
	priority     uint8
}

func (searchresult *searchResult) GetPriority() uint8 {
	return searchresult.priority
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
