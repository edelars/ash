package command_prompt

import "fmt"

type CommandPrompt struct {
	template string
}

func NewCommandPrompt(template string) CommandPrompt {
	if template == "" {
		template = "ash> "
	}
	return CommandPrompt{template: template}
}

func (c CommandPrompt) GetPrompt(outputChan chan byte) {
	p := fmt.Sprintf("\n\r%s", c.template)
	for _, v := range p {
		outputChan <- byte(v)
	}
}
