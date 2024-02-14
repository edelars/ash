package history

import (
	"fmt"

	"ash/internal/commands"
	"ash/internal/dto"
	"ash/internal/storage"
)

const (
	constManagerName = "History"
	limit            = 30
)

type historyManager struct {
	storage     storage.StorageIface
	inputSet    func(r []rune)
	resultLimit int
}

func NewHistoryManager(storage storage.StorageIface, inputSet func(r []rune)) historyManager {
	return historyManager{
		storage:     storage,
		inputSet:    inputSet,
		resultLimit: limit,
	}
}

// Only the first patterns[0] will be proceed. This manager used only for autocomplete/history and no reason for multisearch
// PrecisionSearch = true:
// do nothing
//
// PrecisionSearch = false:
//
//		if pattern == "" give top history for current dir and from top used commands
//	 if pattern != "" search "pattern*" in storage.
//		.
func (m *historyManager) SearchCommands(iContext dto.InternalContextIface, resultChan chan dto.CommandManagerSearchResult, patterns ...dto.PatternIface) {
	var data []dto.CommandIface

	defer func() {
		commandManager := commands.NewCommandManager(constManagerName, 90, true, data...)
		commandManager.SearchCommands(iContext, resultChan, patterns...)
	}()

	var res []storage.StorageResult

	for _, pattern := range patterns {
		if pattern.IsPrecisionSearch() {
			continue
		}
		switch pattern.GetPattern() {
		case "":
			res = m.storage.GetTopHistoryByDirs(iContext.GetCurrentDir(), m.resultLimit)
		default:
			res = m.storage.GetTopHistoryByPattern(pattern.GetPattern(), m.resultLimit)
		}
		break // only first item
	}
	data = m.convertStorageItems(res)
}

func (m *historyManager) convertStorageItems(items []storage.StorageResult) (data []dto.CommandIface) {
	for _, v := range items {
		data = append(data, commands.NewPseudoCommand(v.GetCommand(), m.inputSet, generateDescription(v.GetDir(), v.GetUsedCount()), v.GetCommand()))
	}
	return data
}

func generateDescription(dir string, used int) string {
	return fmt.Sprintf("from:%s used:%d", dir, used)
}
