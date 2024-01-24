package internal_context

import (
	"context"

	"ash/internal/dto"
)

type InternalContext struct {
	inputChan          chan byte
	outputChan         chan byte
	errs               chan error
	currentKeyPressed  byte
	ctx                context.Context
	currentInputBuffer []byte
	executionList      []dto.CommandIface
}

func NewInternalContext(ctx context.Context, inputChan chan byte, outputChan chan byte, errs chan error) InternalContext {
	return InternalContext{
		ctx:        ctx,
		inputChan:  inputChan,
		outputChan: outputChan,
		errs:       errs,
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

func (i InternalContext) WithCurrentInputBuffer(b []byte) dto.InternalContextIface {
	i.currentInputBuffer = b
	return i
}

func (i InternalContext) GetCurrentInputBuffer() []byte {
	return i.currentInputBuffer
}

func (i InternalContext) GetCTX() context.Context {
	return i.ctx
}

func (i InternalContext) GetInputChan() chan byte {
	return i.inputChan
}

func (i InternalContext) GetOutputChan() chan byte {
	return i.outputChan
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
	return func(msg string) {
		msg = "\n\r" + msg
		for _, b := range []byte(msg) {
			i.GetOutputChan() <- b
		}
	}
}
