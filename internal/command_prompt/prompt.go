package command_prompt

type CommandPrompt struct{}

func NewCommandPromt() CommandPrompt {
	return CommandPrompt{}
}

func (c CommandPrompt) GetPrompt(outputChan chan byte) {
	p := "\n\rash> "
	for _, v := range p {
		outputChan <- byte(v)
	}
}
