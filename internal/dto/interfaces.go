package dto

import (
	"context"

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
	WithArgs(args string) CommandIface
	GetArgs() string
}

type ExecF func(internalC InternalContextIface)

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
	GetOutputChan() chan byte
	GetErrChan() chan error
	WithExecutionList(executionList []CommandIface) InternalContextIface
	GetExecutionList() []CommandIface
	GetPrintFunction() func(msg string)
}
