package list

import (
	"fmt"

	"ash/internal/commands"
	"ash/internal/dto"
)

const (
	cmdNameEcho = "echo"
	cmdDescEcho = "Displays OS Env or internal Var"
)

func NewEchoCommand() *commands.Command {
	return commands.NewCommandWithExtendedInfo(cmdNameEcho,
		func(iContext dto.InternalContextIface, args []string) dto.ExecResult {
			el := iContext.GetExecutionList()
			if len(el) == 1 && len(args) == 1 {
				var res string

				if res = iContext.GetVariable(dto.Variable(args[0])); res == "" {
					res = iContext.GetEnv(args[0])
				}

				if res != "" {
					iContext.GetPrintFunction()(fmt.Sprintf("Founded: %s:%s\n", args[0], res))
				}
			}
			return dto.CommandExecResultStatusOk
		}, true, cmdDescEcho, cmdNameEcho)
}
