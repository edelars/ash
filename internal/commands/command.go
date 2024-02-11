package commands

import "ash/internal/dto"

func NewCommand(name string, execFunc dto.ExecutionFunction, mustPrepareExecutionList bool) *Command {
	return &Command{
		execFunc:                 execFunc,
		name:                     name,
		mustPrepareExecutionList: mustPrepareExecutionList,
	}
}

type Command struct {
	mustPrepareExecutionList bool
	weight                   uint8
	execFunc                 dto.ExecutionFunction
	name                     string
	displayName              string
	args                     []string
}

func (c *Command) GetMathWeight() uint8 {
	return c.weight
}

func (c *Command) SetMathWeight(weight uint8) {
	c.weight = weight
}

func (c *Command) GetExecFunc() dto.ExecutionFunction {
	return c.execFunc
}

func (c *Command) GetName() string {
	return c.name
}

func (c *Command) GetDisplayName() string {
	if c.displayName == "" {
		return c.name
	}
	return c.displayName
}

func (c *Command) SetDisplayName(s string) {
	c.displayName = s
}

func (c *Command) WithArgs(args []string) dto.CommandIface {
	res := *c
	res.args = args
	return &res
}

func (c *Command) GetArgs() []string {
	return c.args
}

func (c *Command) MustPrepareExecutionList() bool {
	return c.mustPrepareExecutionList
}
