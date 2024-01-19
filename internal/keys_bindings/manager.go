package keys_bindings

import (
	"ash/internal/commands"
	"ash/internal/dto"
)

type KeyBindingsManager struct {
	bindings map[int]dto.CommandIface
}

func NewKeyBindingsManager(configLoader configLoaderIface, commandRouter commandRouterIface) KeyBindingsManager {
	kb := KeyBindingsManager{bindings: make(map[int]dto.CommandIface)}

	var patterns []dto.PatternIface

	m := make(map[string]int)

	for _, kb := range configLoader.GetKeysBindings() {
		patterns = append(patterns, commands.NewPattern(kb.Action, true))
		m[kb.Action] = kb.Key
	}
	sr := commandRouter.SearchCommands(patterns...)

	for _, pattern := range patterns {
		cmsr := sr.GetDataByPattern(pattern)
		if len(cmsr) == 1 {
			if commands := cmsr[0].GetCommands(); len(commands) == 1 {
				if key, ok := m[commands[0].GetName()]; ok {
					kb.bindings[key] = commands[0]
				}
			}
		}
	}
	return kb
}

type (
	configLoaderIface interface {
		GetKeysBindings() []struct {
			Key    int
			Action string
		}
	}
	commandRouterIface interface {
		SearchCommands(patterns ...dto.PatternIface) dto.CommandRouterSearchResult
	}
)

func (k KeyBindingsManager) GetCommandByKey(key int) dto.CommandIface {
	return k.bindings[key]
}
