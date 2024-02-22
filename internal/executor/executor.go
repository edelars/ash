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

func (r *commandExecutor) Execute(iContext dto.InternalContextIface) dto.ExecResult {
	if mainCommand := r.keyBindingManager.GetCommandByKey(iContext.GetLastKeyPressed()); mainCommand != nil {
		var err error
		if mainCommand.MustPrepareExecutionList() {
			iContext.GetPrintFunction()("\n")
			if iContext, err = r.prepareExecutionList(iContext); err != nil {
				iContext.GetPrintFunction()(fmt.Sprintf("Error execute: %s\n", err.Error()))
			}
		}
		r := mainCommand.GetExecFunc()(iContext, mainCommand.GetArgs())
		return r
	} else {
		return dto.CommandExecResultNotDoAnyting
		// iContext.GetPrintFunction()(fmt.Sprintf("Error, unknown key: %d", iContext.GetLastKeyPressed()))
	}
	// return dto.CommandExecResultStatusOk
}

func (r *commandExecutor) prepareExecutionList(iContext dto.InternalContextIface) (dto.InternalContextIface, error) {
	var executionList []dto.CommandIface

	pattrensArr, argsArr := splitToArrays(string(iContext.GetCurrentInputBuffer()))
	crsr := r.commandRouter.SearchCommands(iContext, pattrensArr...)
	for counter, pattern := range pattrensArr {
		cmsr := crsr.GetDataByPattern(pattern)
		switch len(cmsr) {
		default:
			commands := cmsr[0].GetCommands()
			// if len(commands) != 1 {
			// return iContext, fmt.Errorf("%w : %s commands: %d", errTooManyCmdFounds, pattern.GetPattern(), len(commands))
			// }
			executionList = append(executionList, commands[0].WithArgs(splitArgsStringToArr(argsArr[counter])))
		case 0:
			return iContext, fmt.Errorf("%w : %s", errCmdNotFounds, pattern.GetPattern())
			// default:
			// return iContext, fmt.Errorf("%w : %s source: %d %s %s", errTooManyCmdFounds, pattern.GetPattern(), len(cmsr), cmsr[0].GetSourceName(), cmsr[1].GetSourceName())
		}
	}

	return iContext.WithExecutionList(executionList), nil
}

type routerIface interface {
	SearchCommands(iContext dto.InternalContextIface, patterns ...dto.PatternIface) dto.CommandRouterSearchResult
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

// "-l -a" to {"-l","-a"}
func splitArgsStringToArr(a string) (res []string) {
	for _, a := range strings.Split(a, " ") {
		if a = strings.TrimSpace(a); a != "" {
			res = append(res, a)
		}
	}
	return res
}
