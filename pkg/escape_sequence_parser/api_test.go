package escape_sequence_parser

import (
	"reflect"
	"testing"
)

func Test_escapeParser_ParseEscapeSequence(t *testing.T) {
	type fields struct {
		terminated    bool
		currentResult *escapeParserResult
	}
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantRes []EscapeSequenceResultIface
	}{
		{
			name:   "asa",
			fields: fields{},
			args: args{
				b: []byte("asa"),
			},
			wantRes: []EscapeSequenceResultIface{
				&escapeParserResult{
					action: EscapeActionNone,
					args:   [][]byte{{0x61}},
				},
				&escapeParserResult{
					action: EscapeActionNone,
					args:   [][]byte{{0x73}},
				},
				&escapeParserResult{
					action: EscapeActionNone,
					args:   [][]byte{{0x61}},
				},
			},
		},

		{
			name:   "\\e + asa",
			fields: fields{},
			args: args{
				b: []byte{0x1b, 0x61, 0x73, 0x61},
			},
			wantRes: []EscapeSequenceResultIface{
				&escapeParserResult{
					action: EscapeActionNone,
					args:   [][]byte{{0x61}},
				},
				&escapeParserResult{
					action: EscapeActionNone,
					args:   [][]byte{{0x73}},
				},
				&escapeParserResult{
					action: EscapeActionNone,
					args:   [][]byte{{0x61}},
				},
			},
		},

		{
			name:   "10;10Hzzzz",
			fields: fields{},
			args: args{
				b: []byte{0x1b, 0x5b, 0x31, 0x30, 0x3b, 0x31, 0x30, 0x48, 0x7a, 0x7a, 0x7a, 0x7a},
			},
			wantRes: []EscapeSequenceResultIface{
				&escapeParserResult{
					action: EscapeActionCursorPosition,
					args:   [][]byte{{0x31, 0x30}, {0x31, 0x30}},
				},
				&escapeParserResult{
					action: EscapeActionNone,
					args:   [][]byte{{0x7a}},
				},
				&escapeParserResult{
					action: EscapeActionNone,
					args:   [][]byte{{0x7a}},
				},
				&escapeParserResult{
					action: EscapeActionNone,
					args:   [][]byte{{0x7a}},
				},
				&escapeParserResult{
					action: EscapeActionNone,
					args:   [][]byte{{0x7a}},
				},
			},
		},

		{
			name:   "1H",
			fields: fields{},
			args: args{
				b: []byte{0x1b, 0x5b, 0x31, 0x48},
			},
			wantRes: []EscapeSequenceResultIface{
				&escapeParserResult{
					action: EscapeActionCursorPosition,
					args:   [][]byte{{0x31}},
				},
			},
		},

		{
			name:   "1A",
			fields: fields{},
			args: args{
				b: []byte{0x1b, 0x5b, 0x31, 0x41},
			},
			wantRes: []EscapeSequenceResultIface{
				&escapeParserResult{
					action: EscapeActionCursorUp,
					args:   [][]byte{{0x31}},
				},
			},
		},

		{
			name:   "1B",
			fields: fields{},
			args: args{
				b: []byte{0x1b, 0x5b, 0x31, 0x42},
			},
			wantRes: []EscapeSequenceResultIface{
				&escapeParserResult{
					action: EscapeActionCursorDown,
					args:   [][]byte{{0x31}},
				},
			},
		},

		{
			name:   "1C",
			fields: fields{},
			args: args{
				b: []byte{0x1b, 0x5b, 0x31, 0x43},
			},
			wantRes: []EscapeSequenceResultIface{
				&escapeParserResult{
					action: EscapeActionCursorForward,
					args:   [][]byte{{0x31}},
				},
			},
		},

		{
			name:   "1D",
			fields: fields{},
			args: args{
				b: []byte{0x1b, 0x5b, 0x31, 0x44},
			},
			wantRes: []EscapeSequenceResultIface{
				&escapeParserResult{
					action: EscapeActionCursorBackward,
					args:   [][]byte{{0x31}},
				},
			},
		},

		{
			name:   "1E",
			fields: fields{},
			args: args{
				b: []byte{0x1b, 0x5b, 0x31, 0x45},
			},
			wantRes: []EscapeSequenceResultIface{
				&escapeParserResult{
					action: EscapeActionCursorNextLine,
					args:   [][]byte{{0x31}},
				},
			},
		},

		{
			name:   "1F",
			fields: fields{},
			args: args{
				b: []byte{0x1b, 0x5b, 0x31, 0x46},
			},
			wantRes: []EscapeSequenceResultIface{
				&escapeParserResult{
					action: EscapeActionCursorPrevLine,
					args:   [][]byte{{0x31}},
				},
			},
		},

		{
			name:   "1G",
			fields: fields{},
			args: args{
				b: []byte{0x1b, 0x5b, 0x31, 0x47},
			},
			wantRes: []EscapeSequenceResultIface{
				&escapeParserResult{
					action: EscapeActionCursorLeft,
					args:   [][]byte{{0x31}},
				},
			},
		},

		{
			name:   "1d",
			fields: fields{},
			args: args{
				b: []byte{0x1b, 0x5b, 0x31, 0x64},
			},
			wantRes: []EscapeSequenceResultIface{
				&escapeParserResult{
					action: EscapeActionCursorTop,
					args:   [][]byte{{0x31}},
				},
			},
		},

		{
			name:   "c",
			fields: fields{},
			args: args{
				b: []byte{0x1b, 0x63},
			},
			wantRes: []EscapeSequenceResultIface{
				&escapeParserResult{
					action: EscapeActionClearScreen,
				},
			},
		},

		{
			name:   "wrong c",
			fields: fields{},
			args: args{
				b: []byte{0x1b, 0x5b, 0x63},
			},
			wantRes: []EscapeSequenceResultIface{
				&escapeParserResult{
					action: EscapeActionNone,
					args:   [][]byte{{0x63}},
				},
			},
		},
		{
			name:   "ZZc",
			fields: fields{},
			args: args{
				b: []byte{0x1b, 0x5a, 0x5a, 0x63},
			},
			wantRes: []EscapeSequenceResultIface{
				&escapeParserResult{
					action: EscapeActionNone,
					args:   [][]byte{{0x5a}},
				},
				&escapeParserResult{
					action: EscapeActionNone,
					args:   [][]byte{{0x5a}},
				},
				&escapeParserResult{
					action: EscapeActionNone,
					args:   [][]byte{{0x63}},
				},
			},
		},

		{
			name:   "cZZ",
			fields: fields{},
			args: args{
				b: []byte{0x1b, 0x63, 0x5a, 0x5a},
			},
			wantRes: []EscapeSequenceResultIface{
				&escapeParserResult{
					action: EscapeActionClearScreen,
				},

				&escapeParserResult{
					action: EscapeActionNone,
					args:   [][]byte{{0x5a}},
				},
				&escapeParserResult{
					action: EscapeActionNone,
					args:   [][]byte{{0x5a}},
				},
			},
		},

		{
			name:   "1J",
			fields: fields{},
			args: args{
				b: []byte{0x1b, 0x5b, 0x31, 0x4a},
			},
			wantRes: []EscapeSequenceResultIface{
				&escapeParserResult{
					action: escapeActionEraseDownUpScreen,
					args:   [][]byte{{0x31}},
				},
			},
		},

		{
			name:   "1G",
			fields: fields{},
			args: args{
				b: []byte{0x1b, 0x5b, 0x31, 0x4b},
			},
			wantRes: []EscapeSequenceResultIface{
				&escapeParserResult{
					action: escapeActionEraseRightLeftLine,
					args:   [][]byte{{0x31}},
				},
			},
		},

		{
			name:   "?25h",
			fields: fields{},
			args: args{
				b: []byte{0x1b, 0x3f, 0x32, 0x35, 0x68},
			},
			wantRes: []EscapeSequenceResultIface{
				&escapeParserResult{
					action: EscapeActionCursorShow,
					args:   [][]byte{{0x32, 0x35}},
				},
			},
		},

		{
			name:   "?25l",
			fields: fields{},
			args: args{
				b: []byte{0x1b, 0x3f, 0x32, 0x35, 0x6c},
			},
			wantRes: []EscapeSequenceResultIface{
				&escapeParserResult{
					action: EscapeActionCursorHide,
					args:   [][]byte{{0x32, 0x35}},
				},
			},
		},

		{
			name:   "@",
			fields: fields{},
			args: args{
				b: []byte{0x1b, 0x5b, 0x40},
			},
			wantRes: []EscapeSequenceResultIface{
				&escapeParserResult{
					action: EscapeActionTextInsertChar,
				},
			},
		},

		{
			name:   "P",
			fields: fields{},
			args: args{
				b: []byte{0x1b, 0x5b, 0x50},
			},
			wantRes: []EscapeSequenceResultIface{
				&escapeParserResult{
					action: EscapeActionTextDeleteChar,
				},
			},
		},

		{
			name:   "X",
			fields: fields{},
			args: args{
				b: []byte{0x1b, 0x5b, 0x58},
			},
			wantRes: []EscapeSequenceResultIface{
				&escapeParserResult{
					action: EscapeActionTextEraseChar,
				},
			},
		},

		{
			name:   "L",
			fields: fields{},
			args: args{
				b: []byte{0x1b, 0x5b, 0x4c},
			},
			wantRes: []EscapeSequenceResultIface{
				&escapeParserResult{
					action: EscapeActionTextInsertLine,
				},
			},
		},

		{
			name:   "M",
			fields: fields{},
			args: args{
				b: []byte{0x1b, 0x5b, 0x4d},
			},
			wantRes: []EscapeSequenceResultIface{
				&escapeParserResult{
					action: EscapeActionTextDeleteLine,
				},
			},
		},

		{
			name:   "S",
			fields: fields{},
			args: args{
				b: []byte{0x1b, 0x5b, 0x53},
			},
			wantRes: []EscapeSequenceResultIface{
				&escapeParserResult{
					action: EscapeActionScrollUp,
				},
			},
		},

		{
			name:   "T",
			fields: fields{},
			args: args{
				b: []byte{0x1b, 0x5b, 0x54},
			},
			wantRes: []EscapeSequenceResultIface{
				&escapeParserResult{
					action: EscapeActionScrollDown,
				},
			},
		},

		{
			name:   "m",
			fields: fields{},
			args: args{
				b: []byte{0x1b, 0x5b, 0x6d},
			},
			wantRes: []EscapeSequenceResultIface{
				&escapeParserResult{
					action: EscapeActionSetColor,
				},
			},
		},

		{
			name:   "=15h",
			fields: fields{},
			args: args{
				b: []byte{0x1b, 0x5b, 0x3d, 0x31, 0x35, 0x68},
			},
			wantRes: []EscapeSequenceResultIface{
				&escapeParserResult{
					action: escapeActionSetScreenMode,
					args:   [][]byte{{0x31, 0x35}},
				},
			},
		},

		{
			name:   "m rgb [38;2;1;2;3m", // 5B 33 38 3B 32 3B 31 3B 32 3B 33 6D
			fields: fields{},
			args: args{
				b: []byte{0x1b, 0x5b, 0x33, 0x38, 0x3b, 0x32, 0x3b, 0x31, 0x3b, 0x32, 0x3b, 0x33, 0x6d},
			},
			wantRes: []EscapeSequenceResultIface{
				&escapeParserResult{
					action: EscapeActionSetColor,
					args:   [][]byte{{0x33, 0x38}, {0x32}, {0x31}, {0x32}, {0x33}},
				},
			},
		},

		{
			name:   "m rgb [48;2;1;2;3m", // 5B 33 48 3B 32 3B 31 3B 32 3B 33 6D
			fields: fields{},
			args: args{
				b: []byte{0x1b, 0x5b, 0x34, 0x38, 0x3b, 0x32, 0x3b, 0x31, 0x3b, 0x32, 0x3b, 0x33, 0x6d},
			},
			wantRes: []EscapeSequenceResultIface{
				&escapeParserResult{
					action: EscapeActionSetColor,
					args:   [][]byte{{0x34, 0x38}, {0x32}, {0x31}, {0x32}, {0x33}},
				},
			},
		},

		{
			name:   "m 256 [38;5;208m", // 5B 33 38 3B 35 3B 32 30 38 6D
			fields: fields{},
			args: args{
				b: []byte{0x1b, 0x5b, 0x33, 0x38, 0x3b, 0x35, 0x3b, 0x32, 0x30, 0x38, 0x6d},
			},
			wantRes: []EscapeSequenceResultIface{
				&escapeParserResult{
					action: EscapeActionSetColor,
					args:   [][]byte{{0x33, 0x38}, {0x35}, {0x32, 0x30, 0x38}},
				},
			},
		},

		{
			name:   "m 256 [48;5;208m", // 5B 34 38 3B 35 3B 32 30 38 6D
			fields: fields{},
			args: args{
				b: []byte{0x1b, 0x5b, 0x34, 0x38, 0x3b, 0x35, 0x3b, 0x32, 0x30, 0x38, 0x6d},
			},
			wantRes: []EscapeSequenceResultIface{
				&escapeParserResult{
					action: EscapeActionSetColor,
					args:   [][]byte{{0x34, 0x38}, {0x35}, {0x32, 0x30, 0x38}},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &escapeParser{
				terminated:    tt.fields.terminated,
				currentResult: tt.fields.currentResult,
			}
			if gotRes := e.ParseEscapeSequence(tt.args.b); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("escapeParser.ParseEscapeSequence() = %+v, want %+v", gotRes, tt.wantRes)
			}
		})
	}
}
