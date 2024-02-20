package internal_context

import (
	"io"
	"os"

	"ash/internal/dto"

	"github.com/nsf/termbox-go"
)

type InternalContext struct {
	currentKeyPressed  uint16
	inputManager       inputManagerIface
	errsChan           chan error
	outputWriter       io.Writer
	inputReader        io.Reader
	currentInputBuffer []rune
	executionList      []dto.CommandIface
	printFunc          func(msg string)
	printCellFunc      func(c []termbox.Cell)
	variables          map[dto.Variable]string
}

func (i InternalContext) GetVariable(v dto.Variable) string {
	switch v {
	case dto.VariableCurDir:
		return i.GetCurrentDir()
	default:
		s, _ := i.variables[v]
		return s
	}
}

func (i InternalContext) WithVariables(vars []dto.VariableSet) dto.InternalContextIface {
	for _, v := range vars {
		i.variables[v.Name] = v.Value
	}
	return i
}

func (i InternalContext) WithOutputWriter(w io.Writer) dto.InternalContextIface {
	i.outputWriter = w
	return i
}

func (i InternalContext) WithInputReader(r io.Reader) dto.InternalContextIface {
	i.inputReader = r
	return i
}

type inputManagerIface interface {
	GetInputEventChan() chan termbox.Event
}

func NewInternalContext(im inputManagerIface, errs chan error, printFunc func(msg string), outputWriter io.Writer, inputReader io.Reader, printCellFunc func(c []termbox.Cell)) *InternalContext {
	return &InternalContext{
		inputManager:  im,
		errsChan:      errs,
		printFunc:     printFunc,
		outputWriter:  outputWriter,
		inputReader:   inputReader,
		printCellFunc: printCellFunc,
		variables:     make(map[dto.Variable]string),
	}
}

func (i InternalContext) GetEnv(envName string) string {
	return os.Getenv(envName)
}

func (i InternalContext) GetCurrentDir() string {
	s, _ := os.Getwd()
	return s
}

func (i InternalContext) WithLastKeyPressed(b uint16) dto.InternalContextIface {
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

func (i InternalContext) GetErrChan() chan error {
	return i.errsChan
}

func (i InternalContext) GetLastKeyPressed() uint16 {
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

func (i InternalContext) GetCellsPrintFunction() func(cells []termbox.Cell) {
	return i.printCellFunc
}

func (i InternalContext) GetInputEventChan() chan termbox.Event {
	return i.inputManager.GetInputEventChan()
}

func (i InternalContext) GetOutputWriter() io.Writer {
	return i.outputWriter
}

func (i InternalContext) GetInputReader() io.Reader {
	return i.inputReader
}
