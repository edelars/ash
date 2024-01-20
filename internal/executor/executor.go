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

type CommandExecutor struct {
	commandRouter     routerIface
	keyBindingManager keyBindingsIface
}

func NewCommandExecutor(commandRouter routerIface, keyBindingManager keyBindingsIface) CommandExecutor {
	return CommandExecutor{
		commandRouter:     commandRouter,
		keyBindingManager: keyBindingManager,
	}
}

func (r CommandExecutor) Execute(internalC dto.InternalContextIface) {
	// mainCommand := r.keyBindingManager.GetCommandByKey(int(internalC.GetLastKeyPressed()))
	internalC.GetOutputChan() <- byte(44)
}

func (r CommandExecutor) prepareExecutionList(internalC dto.InternalContextIface) (dto.InternalContextIface, error) {
	var executionList []dto.CommandIface

	pattrensArr, argsArr := splitToArray(string(internalC.GetCurrentInputBuffer()))
	crsr := r.commandRouter.SearchCommands(pattrensArr...)
	for counter, pattern := range pattrensArr {
		cmsr := crsr.GetDataByPattern(pattern)
		switch len(cmsr) {
		case 1:
			commands := cmsr[0].GetCommands()
			if len(commands) != 1 {
				return nil, fmt.Errorf("%w : %s", errTooManyCmdFounds, pattern.GetPattern())
			}
			executionList = append(executionList, commands[0].WithArgs(argsArr[counter]))
		case 0:
			return nil, fmt.Errorf("%w : %s", errCmdNotFounds, pattern.GetPattern())
		default:
			return nil, fmt.Errorf("%w : %s", errTooManyCmdFounds, pattern.GetPattern())
		}
	}

	return internalC.WithExecutionList(executionList), nil
}

type routerIface interface {
	SearchCommands(patterns ...dto.PatternIface) dto.CommandRouterSearchResult
}

type keyBindingsIface interface {
	GetCommandByKey(key int) dto.CommandIface
}

// git clone http://ya.ru | grep ok | >> out.txt
// Parse this to 3 cmd list without trailing " " space and put args for cmd in args
func splitToArray(s string) (res []dto.PatternIface, args []string) {
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
