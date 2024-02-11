package dto

import (
	"context"
	"io"

	"github.com/nsf/termbox-go"
)

type CommandRouterSearchResult interface {
	GetDataByPattern(pattern PatternIface) []CommandManagerSearchResult
}

type CommandManagerSearchResult interface {
	GetSourceName() string
	GetCommands() []CommandIface
	GetPattern() PatternIface
	Founded() int
	GetPriority() uint8
}

type CommandIface interface {
	GetMathWeight() uint8 // 0-100%
	SetMathWeight(weight uint8)
	GetExecFunc() ExecutionFunction
	GetName() string // execution name
	WithArgs(args []string) CommandIface
	GetArgs() []string
	MustPrepareExecutionList() bool // current user input ready for exec and need to prepare exec list
	GetDisplayName() string         // display name
	GetDescription() string         // second disply field for autocomplete
}

type ExecutionFunction func(internalC InternalContextIface, args []string) ExecResult // command result. 0 ok - done, -1 there will be a new user command (ex: for backspace)

type ExecResult int8

// Result status of execution, default CommandExecResultStatusOk
const (
	CommandExecResultStatusOk = ExecResult(0)

	CommandExecResultMainExit     = ExecResult(11) // ash is exiting
	CommandExecResultNewUserInput = ExecResult(12) // put new user input
	CommandExecResultNotDoAnyting = ExecResult(13)
)

// type ExecutionFunction func(internalC InternalContextIface, args []string) ExecResult // command result. 0 ok - done, -1 there will be a new user command (ex: for backspace)

type CommandExecResult interface {
	GetExitCode() int
	GetResultUserInput() []rune
}

type PatternIface interface {
	GetPattern() string
	IsPrecisionSearch() bool
}

type InternalContextIface interface {
	GetEnvList() []string
	GetEnv(envName string) string
	GetCurrentDir() string
	WithLastKeyPressed(b byte) InternalContextIface
	WithCurrentInputBuffer(b []rune) InternalContextIface
	GetCurrentInputBuffer() []rune
	GetLastKeyPressed() byte
	GetCTX() context.Context
	GetInputEventChan() chan termbox.Event
	GetErrChan() chan error
	WithExecutionList(executionList []CommandIface) InternalContextIface
	GetExecutionList() []CommandIface
	GetPrintFunction() func(msg string)
	// console I/O
	GetOutputWriter() io.Writer
	GetInputReader() io.Reader
	WithOutputWriter(io.Writer) InternalContextIface
	WithInputReader(io.Reader) InternalContextIface
}
