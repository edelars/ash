package list

import (
	"ash/internal/commands"
	"ash/internal/dto"

	"gopkg.in/yaml.v3"
)

func NewConfigCommand(cfg CfgManager) *commands.Command {
	return commands.NewCommand("_config",
		func(internalC dto.InternalContextIface, _ []string) dto.ExecResult {
			data, err := yaml.Marshal(cfg.GetConfig())
			if err == nil {
				w := internalC.GetOutputWriter()
				w.Write(data)
			}
			return dto.CommandExecResultStatusOk
		}, true)
}

type CfgManager interface {
	GetConfig() interface{}
}
