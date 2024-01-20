package commands

import "ash/internal/dto"

func NewCommand(name string, execFunc dto.ExecF) *Command {
	return &Command{
		execFunc: execFunc,
		name:     name,
	}
}

type Command struct {
	weight   int8
	execFunc dto.ExecF
	name     string
	args     string
}

func (c *Command) GetMathWeight() int8 {
	return c.weight
}

func (c *Command) SetMathWeight(weight int8) {
	c.weight = weight
}

func (c *Command) GetExecFunc() dto.ExecF {
	return c.execFunc
}

func (c *Command) GetName() string {
	return c.name
}

func (c *Command) WithArgs(args string) dto.CommandIface {
	res := *c
	res.args = args
	return &res
}

func (c *Command) GetArgs() string {
	return c.args
}
