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

	var writer bytes.Buffer

	// 1 test
	ictx := internal_context.NewInternalContext(context.Background(), nil, nil, nil, &writer, reader).WithExecutionList([]dto.CommandIface{commands.NewCommand("1", nil, false), commands.NewCommand("2", nil, false), commands.NewCommand("3", nil, false)})
	res := executeCommands(ictx, nil, execCmdImpl)

	assert.Equal(t, 0, res)
	assert.Equal(t, b, writer.Bytes())

	// 2 test
	writer.Reset()
	ictx = internal_context.NewInternalContext(context.Background(), nil, nil, nil, &writer, reader).WithExecutionList([]dto.CommandIface{commands.NewCommand("1", nil, false), commands.NewCommand("err", nil, false), commands.NewCommand("3", nil, false)})
	res = executeCommands(ictx, nil, execCmdImpl)

	assert.Equal(t, 1, res)
	assert.Equal(t, "err", writer.String())
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
	if cmd.GetName() == "err" {
		st.code = 1
		st.output = []byte("err")
	}
	w.Write(doneBuf)
	w.Close()
	return st
}
