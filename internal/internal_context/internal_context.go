package internal_context

import (
	"context"

	"ash/internal/dto"

	"github.com/nsf/termbox-go"
)

type InternalContext struct {
	im                 inputManager
	errs               chan error
	currentKeyPressed  byte
	ctx                context.Context
	currentInputBuffer []rune
	executionList      []dto.CommandIface
	printFunc          func(msg string)
}
type inputManager interface {
	GetInputEventChan() chan termbox.Event
}

func NewInternalContext(ctx context.Context, im inputManager, errs chan error, printFunc func(msg string)) InternalContext {
	return InternalContext{
		ctx:       ctx,
		im:        im,
		errs:      errs,
		printFunc: printFunc,
	}
}

func (i InternalContext) GetEnvList() []string {
	panic("not implemented") // TODO: Implement
}

func (i InternalContext) GetEnv(envName string) string {
	panic("not implemented") // TODO: Implement
}

func (i InternalContext) GetCurrentDir() string {
	panic("not implemented") // TODO: Implement
}

func (i InternalContext) WithLastKeyPressed(b byte) dto.InternalContextIface {
	i.currentKeyPressed = b
	return i
}

func (i InternalContext) WithCurrentInputBuffer(b []rune) dto.InternalContextIface {
	i.currentInputBuffer = b
	return i
}

func (i InternalContext) GetCurrentInputBuffer() []rune {
	return i.currentInputBuffer
}

func (i InternalContext) GetCTX() context.Context {
	return i.ctx
}

func (i InternalContext) GetErrChan() chan error {
	return i.errs
}

func (i InternalContext) GetLastKeyPressed() byte {
	return i.currentKeyPressed
}

func (i InternalContext) WithExecutionList(executionList []dto.CommandIface) dto.InternalContextIface {
	i.executionList = executionList
	return i
}

func (i InternalContext) GetExecutionList() []dto.CommandIface {
	return i.executionList
}

func (i InternalContext) GetPrintFunction() func(msg string) {
	return i.printFunc
}

func (i InternalContext) GetInputEventChan() chan termbox.Event {
	return i.im.GetInputEventChan()
}
