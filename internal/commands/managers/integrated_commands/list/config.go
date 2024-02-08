package list

import (
	"encoding/json"
	"fmt"

	"ash/internal/commands"
	"ash/internal/dto"
)

func NewConfigCommand(cfg CfgManager) *commands.Command {
	return commands.NewCommand("_config",
		func(internalC dto.InternalContextIface, _ []string) dto.ExecResult {
			output, _ := json.MarshalIndent(cfg.GetConfig(), "", "\t")
			internalC.GetPrintFunction()(fmt.Sprintf("%s\n", output)) // TODO via writer
			return dto.CommandExecResultStatusOk
		}, true)
}

type CfgManager interface {
	GetConfig() interface{}
}
