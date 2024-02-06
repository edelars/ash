package commands

type CommandExecutionResult struct {
	code  int
	input []rune
}

func NewCommandExecutionResult(code int, userInput []rune) CommandExecutionResult {
	return CommandExecutionResult{
		code:  code,
		input: userInput,
	}
}

func (c CommandExecutionResult) GetExitCode() int {
	return c.code
}

func (c CommandExecutionResult) GetResultUserInput() []rune {
	return c.input
}
