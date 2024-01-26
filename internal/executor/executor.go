package executor

import (
	"errors"
	"fmt"
	"strings"

	"ash/internal/commands"
	"ash/internal/dto"
)

var (
	errTooManyCmdFounds = errors.New("Too many commands found")
	errCmdNotFounds     = errors.New("Command not found")
)

type commandExecutor struct {
	commandRouter     routerIface
	keyBindingManager keyBindingsIface
}

func NewCommandExecutor(commandRouter routerIface, keyBindingManager keyBindingsIface) commandExecutor {
	return commandExecutor{
		commandRouter:     commandRouter,
		keyBindingManager: keyBindingManager,
	}
}

func (r commandExecutor) Execute(iContext dto.InternalContextIface) {
	if mainCommand := r.keyBindingManager.GetCommandByKey(uint16(iContext.GetLastKeyPressed())); mainCommand != nil {
		var err error
		if iContext, err = r.prepareExecutionList(iContext); err != nil {
			iContext.GetPrintFunction()(fmt.Sprintf("Error execute: %s", err.Error()))
		}
		mainCommand.GetExecFunc()(iContext)
	} else {
		iContext.GetPrintFunction()(fmt.Sprintf("Error, unknown key: %d", iContext.GetLastKeyPressed()))
	}
}

func (r commandExecutor) prepareExecutionList(iContext dto.InternalContextIface) (dto.InternalContextIface, error) {
	var executionList []dto.CommandIface

	pattrensArr, argsArr := splitToArrays(string(iContext.GetCurrentInputBuffer()))
	crsr := r.commandRouter.SearchCommands(pattrensArr...)
	for counter, pattern := range pattrensArr {
		cmsr := crsr.GetDataByPattern(pattern)
		switch len(cmsr) {
		case 1:
			commands := cmsr[0].GetCommands()
			if len(commands) != 1 {
				return iContext, fmt.Errorf("%w : %s", errTooManyCmdFounds, pattern.GetPattern())
			}
			executionList = append(executionList, commands[0].WithArgs(argsArr[counter]))
		case 0:
			return iContext, fmt.Errorf("%w : %s", errCmdNotFounds, pattern.GetPattern())
		default:
			return iContext, fmt.Errorf("%w : %s", errTooManyCmdFounds, pattern.GetPattern())
		}
	}

	return iContext.WithExecutionList(executionList), nil
}

type routerIface interface {
	SearchCommands(patterns ...dto.PatternIface) dto.CommandRouterSearchResult
}

type keyBindingsIface interface {
	GetCommandByKey(key uint16) dto.CommandIface
}

// git clone http://ya.ru | grep ok | >> out.txt
// Parse this to 3 cmd list without trailing " " space and put args for cmd in args
func splitToArrays(s string) (res []dto.PatternIface, args []string) {
	for _, v := range strings.Split(s, "|") {
		v = strings.TrimSpace(v)
		if len(v) == 0 {
			continue
		}
		cmd, arg, _ := strings.Cut(strings.TrimSpace(v), " ")
		res = append(res, commands.NewPattern(strings.TrimSpace(cmd), true))
		args = append(args, strings.TrimSpace(arg))
	}
	return res, args
}
