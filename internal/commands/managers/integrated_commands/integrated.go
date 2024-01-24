package integrated

import (
	"ash/internal/commands"
	"ash/internal/commands/managers/integrated_commands/list"
)

func NewIntegratedManager(configManager list.CfgManager) (im commands.CommandManagerIface) {
	return commands.NewCommandManager(
		list.NewExitCommand(),
		list.NewCDCommand(),
		list.NewEchoCommand(),
		list.NewExportCommand(),
		list.NewKeyCommand(),
		list.NewLogoutCommand(),
		list.NewConfigCommand(configManager),
	)
}
