package keys_bindings

import (
	"ash/internal/dto"
)

type KeyBindingsManager struct {
	bindings map[int]dto.CommandIface
}

func NewKeyBindingsManager(configLoader configLoaderIface, commandManager commandManagerIface) KeyBindingsManager {
	return KeyBindingsManager{}
}

type (
	configLoaderIface   interface{}
	commandManagerIface interface{}
)

func (k KeyBindingsManager) GetCommandByKey(key int) dto.CommandIface {
	return nil
}
