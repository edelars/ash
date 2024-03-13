package list

import (
	"fmt"

	"ash/internal/commands"
	"ash/internal/dto"
)

const (
	cmdNameVersion = "_version"
	cmdDescVersion    = "Displays version of ash"
)

func NewVersionCommand(version, buildTime, commit, branchName string) *commands.Command {
	return commands.NewCommandWithExtendedInfo(cmdNameVersion,
		func(iContext dto.InternalContextIface, _ []string) dto.ExecResult {
			iContext.GetPrintFunction()(fmt.Sprintf("Version: %s, buildTime: %s, commit: %s, branchName: %s/n", version, buildTime, commit, branchName))
			return dto.CommandExecResultStatusOk
		}, true, cmdDescVersion, cmdNameVersion)
}
