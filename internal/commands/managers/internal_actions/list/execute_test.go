package list

import (
	"bytes"
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
	ictx := internal_context.NewInternalContext(nil, nil, func(s string) { return }, &writer, reader, nil).WithExecutionList([]dto.CommandIface{commands.NewCommand("1", nil, false), commands.NewCommand("2", nil, false), commands.NewCommand("3", nil, false)})
	res := executeCommands(ictx, nil, execCmdImpl)

	assert.Equal(t, dto.CommandExecResultStatusOk, res)
	assert.Equal(t, "123123", writer.String())

	// 2 test
	for i := 0; i < 100; i++ {
		reader.Reset(b)
		writer.Reset()
		ictx = internal_context.NewInternalContext(nil, nil, func(s string) { return }, &writer, reader, nil).WithExecutionList([]dto.CommandIface{commands.NewCommand("1", nil, false), commands.NewCommand("err", nil, false), commands.NewCommand("3", nil, false)})
		res = executeCommands(ictx, nil, execCmdImpl)

		assert.Equal(t, dto.CommandExecResultMainExit, res)
		assert.Equal(t, "err", writer.String())
	}

	// 3 test
	for i := 0; i < 100; i++ {

		reader.Reset(b)
		writer.Reset()
		ictx = internal_context.NewInternalContext(nil, nil, func(s string) { return }, &writer, reader, nil).WithExecutionList([]dto.CommandIface{commands.NewCommand("err", nil, false)})
		res = executeCommands(ictx, nil, execCmdImpl)

		assert.Equal(t, dto.CommandExecResultMainExit, res)
		assert.Equal(t, "err", writer.String())
	}

	// 4 test
	for i := 0; i < 100; i++ {
		reader.Reset(b)
		writer.Reset()
		ictx = internal_context.NewInternalContext(nil, nil, func(s string) { return }, &writer, reader, nil).WithExecutionList([]dto.CommandIface{commands.NewCommand("1", nil, false), commands.NewCommand("2", nil, false), commands.NewCommand("3", nil, false)})
		res = executeCommands(ictx, nil, execCmdImpl)

		assert.Equal(t, dto.CommandExecResultStatusOk, res)
		assert.Equal(t, "123123", writer.String())
	}

	// 5 test
	for i := 0; i < 100; i++ {
		writer.Reset()
		reader.Reset(b)
		ictx = internal_context.NewInternalContext(nil, nil, func(s string) { return }, &writer, reader, nil).WithExecutionList([]dto.CommandIface{commands.NewCommand("1", nil, false), commands.NewCommand("2", nil, false), commands.NewCommand("3", nil, false), commands.NewCommand("4", nil, false), commands.NewCommand("5", nil, false)})
		res = executeCommands(ictx, nil, execCmdImpl)

		assert.Equal(t, dto.CommandExecResultStatusOk, res)
		assert.Equal(t, "12312345", writer.String())
	}

	// 6 test
	b = []byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vivamus scelerisque purus eu erat egestas volutpat. Ut sodales erat ut bibendum scelerisque. Mauris ac dui accumsan augue varius vehicula vel eu ex. Etiam aliquam orci ex, eget porttitor metus pulvinar vitae. Suspendisse scelerisque condimentum dolor eu varius. Nulla ut porttitor turpis. Sed venenatis ultrices sem quis congue. Sed mi erat, vestibulum vel aliquam et, ultrices et lorem. Vestibulum vulputate ligula in libero feugiat euismod. Pellentesque consectetur leo vel urna semper sollicitudin. Sed porta libero eu bibendum egestas. Donec eget vestibulum sem. Suspendisse tempus sit amet sem vel gravida. Proin luctus libero nec lacus tincidunt, a ultrices nibh viverra./n Suspendisse ut nisl tempus, pretium dui a, blandit ligula. Curabitur nec pretium metus. Donec dui dolor, venenatis id commodo eget, dictum ac ex. Vivamus rhoncus euismod mauris et vulputate. Cras luctus aliquam vestibulum. Aenean id lectus tortor. Etiam consectetur libero in risus integer.")
	reader = bytes.NewReader(b)
	writer.Reset()
	reader.Reset(b)
	ictx = internal_context.NewInternalContext(nil, nil, func(s string) { return }, &writer, reader, nil).WithExecutionList([]dto.CommandIface{commands.NewCommand("1", nil, false)})
	res = executeCommands(ictx, nil, execCmdImpl2)

	assert.Equal(t, dto.CommandExecResultStatusOk, res)
	assert.Equal(t, string(b), writer.String())

	// 7 test
	writer.Reset()
	reader.Reset(b)
	ictx = internal_context.NewInternalContext(nil, nil, func(s string) { return }, &writer, reader, nil).WithExecutionList([]dto.CommandIface{commands.NewCommand("1", nil, false), commands.NewCommand("2", nil, false), commands.NewCommand("3", nil, false), commands.NewCommand("4", nil, false), commands.NewCommand("5", nil, false)})
	res = executeCommands(ictx, nil, execCmdImpl2)

	assert.Equal(t, dto.CommandExecResultStatusOk, res)
	assert.Equal(t, string(b), writer.String())
}

func execCmdImpl(iContext dto.InternalContextIface, r io.Reader, w *io.PipeWriter, cmd dto.CommandIface) (st stResult) {
	var doneBuf []byte
	buf := make([]byte, 1024)

	for {
		n, err := r.Read(buf)
		if err == io.EOF {
			break
		}
		doneBuf = append(doneBuf, buf[:n]...)
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

func execCmdImpl2(iContext dto.InternalContextIface, r io.Reader, w *io.PipeWriter, cmd dto.CommandIface) (st stResult) {
	var doneBuf []byte
	buf := make([]byte, 1024)

	for {
		n, err := r.Read(buf)
		if err == io.EOF {
			break
		}
		doneBuf = append(doneBuf, buf[:n]...)
	}
	st.code = dto.CommandExecResultStatusOk
	w.Write(doneBuf)
	w.Close()
	return st
}
