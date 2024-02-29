package escape_sequence_parser

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_escapeParserResult_addToLastArg(t *testing.T) {
	h := escapeParserResult{}
	h.addToLastArg(0x1b)
	assert.Equal(t, 1, len(h.args))
	assert.Equal(t, 1, len(h.args[0]))
	assert.Equal(t, uint8(27), h.args[0][0])

	h.addToLastArg(0x1e)
	assert.Equal(t, 1, len(h.args))
	assert.Equal(t, 2, len(h.args[0]))
	assert.Equal(t, uint8(27), h.args[0][0])
	assert.Equal(t, uint8(30), h.args[0][1])

	h.args = append(h.args, []byte{})
	h.addToLastArg(0x1f)
	assert.Equal(t, 2, len(h.args))
	assert.Equal(t, 2, len(h.args[0]))
	assert.Equal(t, uint8(27), h.args[0][0])
	assert.Equal(t, uint8(30), h.args[0][1])
	assert.Equal(t, uint8(31), h.args[1][0])
}

func Test_escapeParserResult_WithRaw(t *testing.T) {
	type fields struct {
		action EscapeAction
		args   [][]byte
	}
	type args struct {
		b byte
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *escapeParserResult
	}{
		{
			name:   "1",
			fields: fields{},
			args: args{
				b: 63,
			},
			want: &escapeParserResult{
				args: [][]byte{{63}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &escapeParserResult{
				action: tt.fields.action,
				args:   tt.fields.args,
			}
			if got := e.WithRaw(tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("escapeParserResult.WithRaw() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_escapeParserResult_addEmptyArg(t *testing.T) {
	h := escapeParserResult{}
	h.addEmptyArg()
	assert.Equal(t, 1, len(h.args))
	assert.Equal(t, 0, len(h.args[0]))
}

func Test_escapeParser_setCurrentAction(t *testing.T) {
	h := escapeParser{}
	h.setCurrentAction(5)
	assert.Equal(t, EscapeAction(5), h.currentResult.action)

	h = escapeParser{}
	h.currentResult = &escapeParserResult{}
	h.setCurrentAction(5)
	assert.Equal(t, EscapeAction(5), h.currentResult.action)
}

func Test_isDigit(t *testing.T) {
	type args struct {
		i byte
	}
	tests := []struct {
		name string
		args args
		want byte
	}{
		{
			name: "0x30",
			args: args{
				i: 0x30,
			},
			want: 0x30,
		},
		{
			name: "0x5f",
			args: args{
				i: 0x5f,
			},
			want: 0x00,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isDigit(tt.args.i); got != tt.want {
				t.Errorf("isDigit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_escapeParserResult_GetAction(t *testing.T) {
	type fields struct {
		action EscapeAction
		args   [][]byte
	}
	tests := []struct {
		name   string
		fields fields
		want   EscapeAction
	}{
		{
			name: "K",
			fields: fields{
				action: escapeActionEraseRightLeftLine,
			},
			want: EscapeActionEraseRight,
		},
		{
			name: "1K",
			fields: fields{
				action: escapeActionEraseRightLeftLine,
				args:   [][]byte{{0x31}},
			},
			want: EscapeActionEraseLeft,
		},
		{
			name: "2K",
			fields: fields{
				action: escapeActionEraseRightLeftLine,
				args:   [][]byte{{0x32}},
			},
			want: EscapeActionEraseLine,
		},

		{
			name: "J",
			fields: fields{
				action: escapeActionEraseDownUpScreen,
			},
			want: EscapeActionEraseDown,
		},
		{
			name: "1J",
			fields: fields{
				action: escapeActionEraseDownUpScreen,
				args:   [][]byte{{0x31}},
			},
			want: EscapeActionEraseUp,
		},
		{
			name: "2J",
			fields: fields{
				action: escapeActionEraseDownUpScreen,
				args:   [][]byte{{0x32}},
			},
			want: EscapeActionEraseScreen,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &escapeParserResult{
				action: tt.fields.action,
				args:   tt.fields.args,
			}
			if got := e.GetAction(); got != tt.want {
				t.Errorf("escapeParserResult.GetAction() = %v, want %v", got, tt.want)
			}
		})
	}
}
