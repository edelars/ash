package dto

import (
	"context"
)

type CommandRouterSearchResult interface {
	GetDataByPattern(pattern PatternIface) []CommandManagerSearchResult
}

type CommandManagerSearchResult interface {
	GetSourceName() string
	GetCommands() []CommandIface
	GetPattern() PatternIface
}

type CommandIface interface {
	GetMathWeight() int8 // 0-100%
	SetMathWeight(weight int8)
	GetExecFunc() ExecF
	GetName() string
}

type ExecF func(ctx context.Context, internalContext InternalContextIface, inputChan chan []byte, outputChan chan []byte)

type PatternIface interface {
	GetPattern() string
	IsPrecisionSearch() bool
}

type InternalContextIface interface {
	GetEnvList() []string
	GetEnv(envName string) string
	GetCurrentDir() string
	WithLastKeyPressed(b byte) InternalContextIface
	WithCurrentInputBuffer(b []byte) InternalContextIface
	GetCurrentInputBuffer() []byte
	GetLastKeyPressed() byte
	GetCTX() context.Context
	GetInputChan() chan byte
	GetOutputChan() chan byte
	GetErrChan() chan error
}
