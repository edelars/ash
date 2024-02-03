package list

import (
	"encoding/json"
	"fmt"

	"ash/internal/commands"
	"ash/internal/dto"
)

func NewConfigCommand(cfg CfgManager) *commands.Command {
	return commands.NewCommand("_config",
		func(internalC dto.InternalContextIface) int {
			output, _ := json.MarshalIndent(cfg.GetConfig(), "", "\t")
			internalC.GetPrintFunction()(fmt.Sprintf("%s", output))
			return 0
		}, true)
}

type CfgManager interface {
	GetConfig() interface{}
}
