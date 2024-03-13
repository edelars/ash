package integrated

import (
	"ash/internal/commands"
	"ash/internal/commands/managers/integrated_commands/list"
)

func NewIntegratedManager(configManager list.CfgManager, version, buildTime, commit, branchName string) (im commands.CommandManagerIface) {
	return commands.NewCommandManager(
		"Internal commands",
		100,
		false,
		list.NewExitCommand(),
		list.NewCDCommand(),
		list.NewEchoCommand(),
		list.NewExportCommand(),
		list.NewKeyCommand(),
		list.NewLogoutCommand(),
		list.NewConfigCommand(configManager),
		list.NewVersionCommand(version, buildTime, commit, branchName),
	)
}
