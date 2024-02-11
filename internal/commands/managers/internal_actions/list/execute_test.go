package list

import (
	"bytes"
	"context"
	"io"
	"testing"

	"ash/internal/commands"
	"ash/internal/dto"
	"ash/internal/internal_context"

	"github.com/stretchr/testify/assert"
)

func Test_executeCommands(t *testing.T) {
	b := []byte("123")
	reader := bytes.NewReader(b)

	writer := bytes.Buffer{}

	// 1 test
	ictx := internal_context.NewInternalContext(context.Background(), nil, nil, func(s string) { return }, &writer, reader).WithExecutionList([]dto.CommandIface{commands.NewCommand("1", nil, false), commands.NewCommand("2", nil, false), commands.NewCommand("3", nil, false)})
	res := executeCommands(ictx, nil, execCmdImpl)

	assert.Equal(t, dto.CommandExecResultStatusOk, res)
	assert.Equal(t, "123123", writer.String())

	// 2 test
	for i := 0; i < 100; i++ {
		reader.Reset(b)
		writer.Reset()
		ictx = internal_context.NewInternalContext(context.Background(), nil, nil, func(s string) { return }, &writer, reader).WithExecutionList([]dto.CommandIface{commands.NewCommand("1", nil, false), commands.NewCommand("err", nil, false), commands.NewCommand("3", nil, false)})
		res = executeCommands(ictx, nil, execCmdImpl)

		assert.Equal(t, dto.CommandExecResultMainExit, res)
		assert.Equal(t, "err", writer.String())
	}

	// 3 test
	for i := 0; i < 100; i++ {

		reader.Reset(b)
		writer.Reset()
		ictx = internal_context.NewInternalContext(context.Background(), nil, nil, func(s string) { return }, &writer, reader).WithExecutionList([]dto.CommandIface{commands.NewCommand("err", nil, false)})
		res = executeCommands(ictx, nil, execCmdImpl)

		assert.Equal(t, dto.CommandExecResultMainExit, res)
		assert.Equal(t, "err", writer.String())
	}

	// 4 test
	for i := 0; i < 100; i++ {
		reader.Reset(b)
		writer.Reset()
		ictx = internal_context.NewInternalContext(context.Background(), nil, nil, func(s string) { return }, &writer, reader).WithExecutionList([]dto.CommandIface{commands.NewCommand("1", nil, false), commands.NewCommand("2", nil, false), commands.NewCommand("3", nil, false)})
		res = executeCommands(ictx, nil, execCmdImpl)

		assert.Equal(t, dto.CommandExecResultStatusOk, res)
		assert.Equal(t, "123123", writer.String())
	}

	// 5 test
	for i := 0; i < 100; i++ {
		writer.Reset()
		reader.Reset(b)
		ictx = internal_context.NewInternalContext(context.Background(), nil, nil, func(s string) { return }, &writer, reader).WithExecutionList([]dto.CommandIface{commands.NewCommand("1", nil, false), commands.NewCommand("2", nil, false), commands.NewCommand("3", nil, false), commands.NewCommand("4", nil, false), commands.NewCommand("5", nil, false)})
		res = executeCommands(ictx, nil, execCmdImpl)

		assert.Equal(t, dto.CommandExecResultStatusOk, res)
		assert.Equal(t, "12312345", writer.String())
	}
}

func execCmdImpl(iContext dto.InternalContextIface, r io.Reader, w *io.PipeWriter, cmd dto.CommandIface) (st stResult) {
	var doneBuf []byte
	buf := make([]byte, 1024)

	for {
		n, err := r.Read(buf)
		if err == io.EOF {
			break
		}
		doneBuf = buf[:n]
	}
	st.code = dto.CommandExecResultStatusOk
	if cmd.GetName() == "err" {
		st.code = dto.CommandExecResultMainExit
		st.output = []byte("err")
	}
	w.Write(append(doneBuf, []byte(cmd.GetName())...))
	w.Close()
	return st
}
