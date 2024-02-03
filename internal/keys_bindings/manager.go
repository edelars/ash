package keys_bindings

import (
	"ash/internal/commands"
	"ash/internal/dto"
)

type KeyBindingsManager struct {
	bindings map[uint16]dto.CommandIface
}

func NewKeyBindingsManager(configLoader configLoaderIface, commandRouter commandRouterIface) KeyBindingsManager {
	kb := KeyBindingsManager{bindings: make(map[uint16]dto.CommandIface)}

	var patterns []dto.PatternIface

	m := make(map[string]uint16)

	for _, kb := range configLoader.GetKeysBindings() {
		patterns = append(patterns, commands.NewPattern(kb.Action, true))
		m[kb.Action] = kb.Key
	}
	sr := commandRouter.SearchCommands(patterns...)
patternLoop:
	for _, pattern := range patterns {
		cmsr := sr.GetDataByPattern(pattern)

		for _, v := range cmsr {
			if commands := v.GetCommands(); len(commands) > 0 {
				if key, ok := m[commands[0].GetName()]; ok {
					kb.bindings[key] = commands[0]
					continue patternLoop
				}
			}
		}
	}
	return kb
}

type (
	configLoaderIface interface {
		GetKeysBindings() []struct {
			Key    uint16
			Action string
		}
	}
	commandRouterIface interface {
		SearchCommands(patterns ...dto.PatternIface) dto.CommandRouterSearchResult
	}
)

func (k KeyBindingsManager) GetCommandByKey(key uint16) dto.CommandIface {
	return k.bindings[key]
}
