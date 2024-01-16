package commands

import (
	"context"

	"ash/internal/internal_context"
)

type ExecF func(ctx context.Context, internalContext internal_context.InternalContextIface, inputChan chan []byte, outputChan chan []byte)

type CommandIface interface {
	GetMathWeight() int8 // 0-100%
	SetMathWeight(weight int8)
	GetExecFunc() ExecF
	GetName() string
}

func NewCommand(name string, execFunc ExecF) *Command {
	return &Command{
		execFunc: execFunc,
		name:     name,
	}
}

type Command struct {
	weight   int8
	execFunc ExecF
	name     string
}

func (c *Command) GetMathWeight() int8 {
	return c.weight
}

func (c *Command) SetMathWeight(weight int8) {
	c.weight = weight
}

func (c *Command) GetExecFunc() ExecF {
	return c.execFunc
}

func (c *Command) GetName() string {
	return c.name
}
