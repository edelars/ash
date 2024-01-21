package integrated

import (
	"ash/internal/commands"
	"ash/internal/commands/managers/integrated_commands/list"
)

func NewIntegratedManager() (im commands.CommandManagerIface) {
	return commands.NewCommandManager(list.NewExitCommand(), list.NewCDCommand(), list.NewEchoCommand(), list.NewExportCommand(), list.NewKeyCommand())
}
