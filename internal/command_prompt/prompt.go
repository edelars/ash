package command_prompt

import (
	"errors"
	"fmt"

	"ash/internal/configuration"
	"ash/internal/dto"

	"github.com/nsf/termbox-go"
)

type CommandPrompt struct {
	template      string
	currentBuffer []rune
}

func NewCommandPrompt(template string) CommandPrompt {
	if template == "" {
		template = "ash> "
	}
	return CommandPrompt{template: template}
}

func (c *CommandPrompt) getPrompt() string {
	return fmt.Sprintf("%s", c.template)
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

func (c *CommandPrompt) Run(iContext dto.InternalContextIface, exec Executor, cfg configuration.ConfigLoader) error {
	promptChan := make(chan struct{}, 1)
	defer close(promptChan)

	promptChan <- struct{}{}

mainLoop:
	for {
		select {
		case ev := <-iContext.GetInputEventChan():
			switch ev.Type {
			case termbox.EventKey:

				if ev.Ch != 0 {
					c.currentBuffer = append(c.currentBuffer, ev.Ch)
					iContext.GetPrintFunction()(string(ev.Ch))
				} else if ev.Key == 32 {
					c.currentBuffer = append(c.currentBuffer, rune(' '))
					iContext.GetPrintFunction()(" ")
				} else {
					if res := exec.Execute(iContext.WithLastKeyPressed(byte(ev.Key)).WithCurrentInputBuffer(c.currentBuffer)); res >= -1 {
						c.currentBuffer = nil
						promptChan <- struct{}{}
					}
				}
			}
		case <-promptChan:
			iContext.GetPrintFunction()(c.getPrompt())
		case <-iContext.GetCTX().Done():
			break mainLoop
		}
	}
	return nil
}

type Executor interface {
	Execute(internalC dto.InternalContextIface) int
}
