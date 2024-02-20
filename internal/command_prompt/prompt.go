package command_prompt

import (
	"errors"

	"ash/internal/configuration"
	"ash/internal/dto"

	"github.com/nsf/termbox-go"
)

type CommandPrompt struct {
	template      []promptItem
	currentBuffer []rune
	stopChan      chan struct{}
	execAdapter   executionAdapter
}

func NewCommandPrompt(template string) CommandPrompt {
	if template == "" {
		template = `[{"value": "ash> ", "bold": true}]`
	}
	return CommandPrompt{template: parsePromptConfigString([]byte(template)), stopChan: make(chan struct{})}
}

// For cases when user input update needed
func (c *CommandPrompt) GetUserInputFunc() func(r []rune) {
	return func(r []rune) {
		c.currentBuffer = r
	}
}

func (c *CommandPrompt) Stop() {
	c.stopChan <- struct{}{}
	defer close(c.stopChan)
}

// Delete rune from user currentBuffer if possible. If not will be error
// Position start from 0
func (c *CommandPrompt) DeleteFromCurrentBuffer(position int) error {
	if position >= len(c.currentBuffer) || position < 0 {
		return errors.New("position is out of buffer")
	}
	b := c.currentBuffer[:position]
	if len(c.currentBuffer) > position+1 {
		b = append(b, c.currentBuffer[position+1])
	}

	c.currentBuffer = b
	return nil
}

func (c *CommandPrompt) DeleteLastSymbolFromCurrentBuffer() error {
	return c.DeleteFromCurrentBuffer(len(c.currentBuffer) - 1)
}

func (c *CommandPrompt) Run(iContext dto.InternalContextIface, exec Executor, cfg configuration.ConfigLoader, enterKey uint16) error {
	promptChan := make(chan struct{}, 1)
	defer close(promptChan)

	c.execAdapter = newExecAdapter(exec, enterKey)
	promptChan <- struct{}{}

	var lastExitStatus dto.ExecResult
mainLoop:
	for {
		select {
		case ev := <-iContext.GetInputEventChan():
			switch ev.Type {
			case termbox.EventKey:
				// panic(fmt.Sprintf("%d %v", ev.Ch, ev.Key))
				if ev.Ch != 0 {
					c.currentBuffer = append(c.currentBuffer, ev.Ch)
					iContext.GetPrintFunction()(string(ev.Ch))
				} else if ev.Key == 32 {
					c.currentBuffer = append(c.currentBuffer, rune(' '))
					iContext.GetPrintFunction()(" ")
				} else {
					v := dto.VariableSet{Name: dto.VariableLastExitCode, Value: string(rune(lastExitStatus))}
					ictx := iContext.WithLastKeyPressed(uint16(ev.Key)).WithCurrentInputBuffer(c.currentBuffer).WithVariables([]dto.VariableSet{v})
					r := exec.Execute(ictx)
					switch r {
					case dto.CommandExecResultNewUserInput:
						iContext.GetPrintFunction()("\n")
					case dto.CommandExecResultNotDoAnyting:
						continue mainLoop
					case dto.CommandExecResultMainExit:
						iContext.GetErrChan() <- errors.New("done")
					default:
						c.currentBuffer = nil
						lastExitStatus = r
						// iContext.GetPrintFunction()("\n")
					}
					promptChan <- struct{}{}
				}
			}
		case <-promptChan:
			iContext.GetCellsPrintFunction()(c.generatePrompt(iContext))
			if c.currentBuffer != nil {
				iContext.GetPrintFunction()(string(c.currentBuffer))
			}
		case <-c.stopChan:
			break mainLoop
		}
	}
	return nil
}

type Executor interface {
	Execute(internalC dto.InternalContextIface) dto.ExecResult
}
