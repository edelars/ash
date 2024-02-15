package command_prompt

import (
	"context"
	"io"
	"reflect"
	"testing"

	"ash/internal/dto"

	"github.com/go-playground/colors"
	"github.com/nsf/termbox-go"
	"github.com/stretchr/testify/assert"
)

func Test_parsePromptConfigString(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		args    args
		wantRes []promptItem
	}{
		{
			name: "1",
			args: args{
				b: []byte(` [{"value": "exe", "color": "", "bold": true,"underline": true }, {"value": "ddd", "color": "3" }]`),
			},
			wantRes: []promptItem{{Value: "exe", Bold: true, Underline: true}, {Value: "ddd", Color: "3"}},
		},
		{
			name: "2",
			args: args{
				b: []byte(` [{"value": "exe", "color": 0, "bold": true,"under": true }, {"valu`),
			},
			wantRes: []promptItem{{Value: constErrParse}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := parsePromptConfigString(tt.args.b); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("parsePromptConfigString() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func Test_hasPrefix(t *testing.T) {
	type args struct {
		s string
		p string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "$",
			args: args{
				s: "$sdfsd",
				p: "$",
			},
			want: "$sdfsd",
		},
		{
			name: "a",
			args: args{
				s: "$sdfsd",
				p: "a",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hasPrefix(tt.args.s, tt.args.p); got != tt.want {
				t.Errorf("hasPrefix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_stringToCells(t *testing.T) {
	type args struct {
		s     string
		color colors.Color
		b     bool
		u     bool
	}
	tests := []struct {
		name    string
		args    args
		wantRes []termbox.Cell
	}{
		{
			name: "1",
			args: args{
				s:     "asd",
				color: nil,
				b:     true,
				u:     false,
			},
			wantRes: []termbox.Cell{{Ch: 'a', Fg: 512}, {Ch: 's', Fg: 512}, {Ch: 'd', Fg: 512}},
		},
		{
			name: "2",
			args: args{
				s:     "asd",
				color: nil,
				b:     false,
				u:     false,
			},
			wantRes: []termbox.Cell{{Ch: 'a', Fg: 0}, {Ch: 's', Fg: 0}, {Ch: 'd', Fg: 0}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := stringToCells(tt.args.s, tt.args.color, tt.args.b, tt.args.u); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("stringToCells() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestCommandPrompt_generatePieceOfPrompt(t *testing.T) {
	e := exeImpl{}
	v := vstorImpl{}
	h := CommandPrompt{execAdapter: &e}

	// 1 test
	res, err := h.generatePieceOfPrompt(&v, promptItem{})
	assert.Equal(t, errEmptyValue, err)
	assert.Equal(t, 11, len(res))

	// 2 test
	res, err = h.generatePieceOfPrompt(&v, promptItem{Value: "$sdfs"})
	assert.NoError(t, err)
	assert.Equal(t, 5, len(res))
	assert.Equal(t, true, v.b)
	assert.Equal(t, false, e.b)

	// 3 test
	v.b = false
	res, err = h.generatePieceOfPrompt(&v, promptItem{Value: "%sdfs"})
	assert.NoError(t, err)
	assert.Equal(t, 5, len(res))
	assert.Equal(t, false, v.b)
	assert.Equal(t, true, e.b)

	// 4 test
	v.b = false
	e.b = false
	res, err = h.generatePieceOfPrompt(&v, promptItem{Value: "sdfs"})
	assert.NoError(t, err)
	assert.Equal(t, 4, len(res))
	assert.Equal(t, false, v.b)
	assert.Equal(t, false, e.b)
}

type exeImpl struct {
	b bool
}

func (exeimpl *exeImpl) ExecCmd(_ dto.InternalContextIface, cmd string) (string, error) {
	exeimpl.b = true
	return cmd, nil
}

type vstorImpl struct {
	b bool
}

func (vstorimpl *vstorImpl) GetCellsPrintFunction() func(cells []termbox.Cell) {
	panic("not implemented") // TODO: Implement
}

func (vstorimpl *vstorImpl) GetEnvList() []string {
	panic("not implemented") // TODO: Implement
}

func (vstorimpl *vstorImpl) GetEnv(envName string) string {
	panic("not implemented") // TODO: Implement
}

func (vstorimpl *vstorImpl) GetCurrentDir() string {
	panic("not implemented") // TODO: Implement
}

func (vstorimpl *vstorImpl) WithLastKeyPressed(b byte) dto.InternalContextIface {
	panic("not implemented") // TODO: Implement
}

func (vstorimpl *vstorImpl) WithCurrentInputBuffer(b []rune) dto.InternalContextIface {
	panic("not implemented") // TODO: Implement
}

func (vstorimpl *vstorImpl) GetCurrentInputBuffer() []rune {
	panic("not implemented") // TODO: Implement
}

func (vstorimpl *vstorImpl) GetLastKeyPressed() byte {
	panic("not implemented") // TODO: Implement
}

func (vstorimpl *vstorImpl) GetCTX() context.Context {
	panic("not implemented") // TODO: Implement
}

func (vstorimpl *vstorImpl) GetInputEventChan() chan termbox.Event {
	panic("not implemented") // TODO: Implement
}

func (vstorimpl *vstorImpl) GetErrChan() chan error {
	panic("not implemented") // TODO: Implement
}

func (vstorimpl *vstorImpl) WithExecutionList(executionList []dto.CommandIface) dto.InternalContextIface {
	panic("not implemented") // TODO: Implement
}

func (vstorimpl *vstorImpl) GetExecutionList() []dto.CommandIface {
	panic("not implemented") // TODO: Implement
}

func (vstorimpl *vstorImpl) GetPrintFunction() func(msg string) {
	panic("not implemented") // TODO: Implement
}

// console I/O
func (vstorimpl *vstorImpl) GetOutputWriter() io.Writer {
	panic("not implemented") // TODO: Implement
}

func (vstorimpl *vstorImpl) GetInputReader() io.Reader {
	panic("not implemented") // TODO: Implement
}

func (vstorimpl *vstorImpl) WithOutputWriter(_ io.Writer) dto.InternalContextIface {
	panic("not implemented") // TODO: Implement
}

func (vstorimpl *vstorImpl) WithInputReader(_ io.Reader) dto.InternalContextIface {
	panic("not implemented") // TODO: Implement
}

func (vstorimpl *vstorImpl) GetVariable(v string) string {
	vstorimpl.b = true
	return v
}
