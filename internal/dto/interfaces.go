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
}

type CommandIface interface {
	GetMathWeight() int8 // 0-100%
	SetMathWeight(weight int8)
	GetExecFunc() ExecF
	GetName() string
	WithArgs(args []string) CommandIface
	GetArgs() []string
	MustPrepareExecutionList() bool // current user input ready for exec and need to prepare exec list
}

type ExecF func(internalC InternalContextIface, args []string) int // command result. 0 ok - done, -1 there will be a new user command (ex: for backspace)

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
