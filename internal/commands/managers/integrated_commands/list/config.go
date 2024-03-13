package list

import (
	"ash/internal/commands"
	"ash/internal/dto"

	"gopkg.in/yaml.v3"
)

const (
	cmdNameConfig = "_config"
	cmdDescConfig = "Displays current configuration"
)

func NewConfigCommand(cfg CfgManager) *commands.Command {
	return commands.NewCommandWithExtendedInfo(cmdNameConfig,
		func(iContext dto.InternalContextIface, _ []string) dto.ExecResult {
			data, err := yaml.Marshal(cfg.GetConfig())
			if err == nil {
				w := iContext.GetOutputWriter()
				w.Write(data)
			}
			return dto.CommandExecResultStatusOk
		}, true, cmdDescConfig, cmdNameConfig)
}

type CfgManager interface {
	GetConfig() interface{}
}
