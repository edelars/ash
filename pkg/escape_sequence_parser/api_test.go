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
