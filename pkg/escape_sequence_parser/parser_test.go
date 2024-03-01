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
		{
			name: "25l",
			fields: fields{
				action: EscapeActionCursorHide,
				args:   [][]byte{{0x32, 0x35}},
			},
			want: EscapeActionCursorHide,
		},
		{
			name: "24l - wrong",
			fields: fields{
				action: EscapeActionCursorHide,
				args:   [][]byte{{0x32, 0x34}},
			},
			want: escapeActionPrivateControlSequence,
		},
		{
			name: "25h",
			fields: fields{
				action: EscapeActionCursorShow,
				args:   [][]byte{{0x32, 0x35}},
			},
			want: EscapeActionCursorShow,
		},
		{
			name: "24h - wrong",
			fields: fields{
				action: EscapeActionCursorShow,
				args:   [][]byte{{0x32, 0x34}},
			},
			want: escapeActionPrivateControlSequence,
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

func Test_escapeParserResult_GetColor(t *testing.T) {
	type fields struct {
		action EscapeAction
		args   [][]byte
	}
	tests := []struct {
		name      string
		fields    fields
		wantColor EscapeColor
		wantBack  bool
	}{
		{
			name: "0m",
			fields: fields{
				action: EscapeActionSetColor,
				args:   [][]byte{{0x30}},
			},
			wantColor: EscapeColorDefault,
			wantBack:  false,
		},
		{
			name: "7m - wrong",
			fields: fields{
				action: EscapeActionSetColor,
				args:   [][]byte{{0x37}},
			},
			wantColor: EscapeColorDefault,
			wantBack:  false,
		},

		{
			name: "30m",
			fields: fields{
				action: EscapeActionSetColor,
				args:   [][]byte{{0x33, 0x30}},
			},
			wantColor: EscapeColorDefault,
			wantBack:  false,
		},

		{
			name: "31m",
			fields: fields{
				action: EscapeActionSetColor,
				args:   [][]byte{{0x33, 0x31}},
			},
			wantColor: EscapeColorRed,
			wantBack:  false,
		},

		{
			name: "32m",
			fields: fields{
				action: EscapeActionSetColor,
				args:   [][]byte{{0x33, 0x32}},
			},
			wantColor: EscapeColorGreen,
			wantBack:  false,
		},

		{
			name: "33m",
			fields: fields{
				action: EscapeActionSetColor,
				args:   [][]byte{{0x33, 0x33}},
			},
			wantColor: EscapeColorYellow,
			wantBack:  false,
		},

		{
			name: "34m",
			fields: fields{
				action: EscapeActionSetColor,
				args:   [][]byte{{0x33, 0x34}},
			},
			wantColor: EscapeColorBlue,
			wantBack:  false,
		},

		{
			name: "35m",
			fields: fields{
				action: EscapeActionSetColor,
				args:   [][]byte{{0x33, 0x35}},
			},
			wantColor: EscapeColorMagenta,
			wantBack:  false,
		},

		{
			name: "36m",
			fields: fields{
				action: EscapeActionSetColor,
				args:   [][]byte{{0x33, 0x36}},
			},
			wantColor: EscapeColorCyan,
			wantBack:  false,
		},

		{
			name: "37m",
			fields: fields{
				action: EscapeActionSetColor,
				args:   [][]byte{{0x33, 0x37}},
			},
			wantColor: EscapeColorWhite,
			wantBack:  false,
		},

		{
			name: "30;1m",
			fields: fields{
				action: EscapeActionSetColor,
				args:   [][]byte{{0x33, 0x30}, {0x31}},
			},
			wantColor: EscapeColorBrightBlack,
			wantBack:  false,
		},

		{
			name: "31;1m",
			fields: fields{
				action: EscapeActionSetColor,
				args:   [][]byte{{0x33, 0x31}, {0x31}},
			},
			wantColor: EscapeColorBrightRed,
			wantBack:  false,
		},

		{
			name: "32;1m",
			fields: fields{
				action: EscapeActionSetColor,
				args:   [][]byte{{0x33, 0x32}, {0x31}},
			},
			wantColor: EscapeColorBrightGreen,
			wantBack:  false,
		},

		{
			name: "33;1m",
			fields: fields{
				action: EscapeActionSetColor,
				args:   [][]byte{{0x33, 0x33}, {0x31}},
			},
			wantColor: EscapeColorBrightYellow,
			wantBack:  false,
		},

		{
			name: "34;1m",
			fields: fields{
				action: EscapeActionSetColor,
				args:   [][]byte{{0x33, 0x34}, {0x31}},
			},
			wantColor: EscapeColorBrightBlue,
			wantBack:  false,
		},

		{
			name: "35;1m",
			fields: fields{
				action: EscapeActionSetColor,
				args:   [][]byte{{0x33, 0x35}, {0x31}},
			},
			wantColor: EscapeColorBrightMagenta,
			wantBack:  false,
		},

		{
			name: "36;1m",
			fields: fields{
				action: EscapeActionSetColor,
				args:   [][]byte{{0x33, 0x36}, {0x31}},
			},
			wantColor: EscapeColorBrightCyan,
			wantBack:  false,
		},

		{
			name: "37;1m",
			fields: fields{
				action: EscapeActionSetColor,
				args:   [][]byte{{0x33, 0x37}, {0x31}},
			},
			wantColor: EscapeColorBrightWhite,
			wantBack:  false,
		},

		{
			name: "40m",
			fields: fields{
				action: EscapeActionSetColor,
				args:   [][]byte{{0x34, 0x30}},
			},
			wantColor: EscapeColorDefault,
			wantBack:  true,
		},

		{
			name: "41m",
			fields: fields{
				action: EscapeActionSetColor,
				args:   [][]byte{{0x34, 0x31}},
			},
			wantColor: EscapeColorRed,
			wantBack:  true,
		},

		{
			name: "42m",
			fields: fields{
				action: EscapeActionSetColor,
				args:   [][]byte{{0x34, 0x32}},
			},
			wantColor: EscapeColorGreen,
			wantBack:  true,
		},

		{
			name: "43m",
			fields: fields{
				action: EscapeActionSetColor,
				args:   [][]byte{{0x34, 0x33}},
			},
			wantColor: EscapeColorYellow,
			wantBack:  true,
		},

		{
			name: "44m",
			fields: fields{
				action: EscapeActionSetColor,
				args:   [][]byte{{0x34, 0x34}},
			},
			wantColor: EscapeColorBlue,
			wantBack:  true,
		},

		{
			name: "45m",
			fields: fields{
				action: EscapeActionSetColor,
				args:   [][]byte{{0x34, 0x35}},
			},
			wantColor: EscapeColorMagenta,
			wantBack:  true,
		},

		{
			name: "46m",
			fields: fields{
				action: EscapeActionSetColor,
				args:   [][]byte{{0x34, 0x36}},
			},
			wantColor: EscapeColorCyan,
			wantBack:  true,
		},

		{
			name: "47m",
			fields: fields{
				action: EscapeActionSetColor,
				args:   [][]byte{{0x34, 0x37}},
			},
			wantColor: EscapeColorWhite,
			wantBack:  true,
		},

		{
			name: "40;1m",
			fields: fields{
				action: EscapeActionSetColor,
				args:   [][]byte{{0x34, 0x30}, {0x31}},
			},
			wantColor: EscapeColorBrightBlack,
			wantBack:  true,
		},

		{
			name: "41;1m",
			fields: fields{
				action: EscapeActionSetColor,
				args:   [][]byte{{0x34, 0x31}, {0x31}},
			},
			wantColor: EscapeColorBrightRed,
			wantBack:  true,
		},

		{
			name: "42;1m",
			fields: fields{
				action: EscapeActionSetColor,
				args:   [][]byte{{0x34, 0x32}, {0x31}},
			},
			wantColor: EscapeColorBrightGreen,
			wantBack:  true,
		},

		{
			name: "43;1m",
			fields: fields{
				action: EscapeActionSetColor,
				args:   [][]byte{{0x34, 0x33}, {0x31}},
			},
			wantColor: EscapeColorBrightYellow,
			wantBack:  true,
		},

		{
			name: "44;1m",
			fields: fields{
				action: EscapeActionSetColor,
				args:   [][]byte{{0x34, 0x34}, {0x31}},
			},
			wantColor: EscapeColorBrightBlue,
			wantBack:  true,
		},

		{
			name: "45;1m",
			fields: fields{
				action: EscapeActionSetColor,
				args:   [][]byte{{0x34, 0x35}, {0x31}},
			},
			wantColor: EscapeColorBrightMagenta,
			wantBack:  true,
		},

		{
			name: "46;1m",
			fields: fields{
				action: EscapeActionSetColor,
				args:   [][]byte{{0x34, 0x36}, {0x31}},
			},
			wantColor: EscapeColorBrightCyan,
			wantBack:  true,
		},

		{
			name: "47;1m",
			fields: fields{
				action: EscapeActionSetColor,
				args:   [][]byte{{0x34, 0x37}, {0x31}},
			},
			wantColor: EscapeColorBrightWhite,
			wantBack:  true,
		},

		{
			name: "47;31m - wrong",
			fields: fields{
				action: EscapeActionSetColor,
				args:   [][]byte{{0x34, 0x37}, {0x33, 0x31}},
			},
			wantColor: EscapeColorDefault,
			wantBack:  false,
		},

		{
			name: "38;5;123m",
			fields: fields{
				action: EscapeActionSetColor,
				args:   [][]byte{{0x33, 0x38}, {0x35}, {0x31, 0x32, 0x33}},
			},
			wantColor: EscapeColor(379), // 256 + value
			wantBack:  false,
		},

		{
			name: "34;5;1m",
			fields: fields{
				action: EscapeActionSetColor,
				args:   [][]byte{{0x34, 0x38}, {0x35}, {0x31}},
			},
			wantColor: EscapeColor(257), // 256 + value
			wantBack:  true,
		},

		{
			name: "38;2;1;2;3m",
			fields: fields{
				action: EscapeActionSetColor,
				args:   [][]byte{{0x33, 0x38}, {0x32}, {0x31}, {0x32}, {0x33}},
			},
			wantColor: EscapeColor(2203351973888),
			wantBack:  false,
		},

		{
			name: "48;2;10;20;30m",
			fields: fields{
				action: EscapeActionSetColor,
				args:   [][]byte{{0x34, 0x38}, {0x32}, {0x31, 0x30}, {0x32, 0x30}, {0x33, 0x30}},
			},
			wantColor: EscapeColor(2242310438912),
			wantBack:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &escapeParserResult{
				action: tt.fields.action,
				args:   tt.fields.args,
			}
			got, got1 := e.GetColorFormat()
			if got != tt.wantColor {
				t.Errorf("escapeParserResult.GetColor() got color = %v, want %v", got, tt.wantColor)
			}
			if got1 != tt.wantBack {
				t.Errorf("escapeParserResult.GetColor() got back = %v, want %v", got1, tt.wantBack)
			}
		})
	}
}
