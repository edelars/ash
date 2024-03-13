package escape_sequence_parser

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_escapeParser_ParseEscapeSequence_assert(t *testing.T) {
	// 1 test
	h := NewEscapeSequenceParser()
	res1 := h.ParseEscapeSequence([]byte{0x1b, 0x5b, 0x33, 0x39, 0x3b, 0x31, 0x32, 0x31}) // [39;121
	assert.Equal(t, 0, len(res1))
	assert.NotNil(t, h.currentResult)
	res2 := h.ParseEscapeSequence([]byte{0x48, 0x32, 0x39, 0x37, 0x33}) // H2973

	assert.Equal(t, 1, len(res2))
	assert.Equal(t, EscapeAction(EscapeActionCursorPosition), res2[0].GetAction())
	n1, n2 := res2[0].GetIntsFromArgs()
	assert.Equal(t, 39, n1)
	assert.Equal(t, 121, n2)
	assert.NotNil(t, h.currentResult)
	assert.Equal(t, EscapeAction(EscapeActionNone), h.currentResult.action)
	assert.Equal(t, []byte{0x32, 0x39, 0x37, 0x33}, h.currentResult.GetRaw())
}

func Test_escapeParser_ParseEscapeSequence(t *testing.T) {
	type fields struct {
		terminated    bool
		currentResult *escapeParserResult
	}
	type args struct {
		b []byte
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		wantRes        []EscapeSequenceResultIface
		wantCurrentRes EscapeSequenceResultIface
	}{
		{
			name: "asa",
			wantCurrentRes: &escapeParserResult{
				action: EscapeActionNone,
				args:   [][]byte{{0x61, 0x73, 0x61}},
			},

			args: args{
				b: []byte("asa"),
			},
		},

		{
			name: "\\e + asa",
			wantCurrentRes: &escapeParserResult{
				action: EscapeActionNone,
				args:   [][]byte{{0x61, 0x73, 0x61}},
			},

			args: args{
				b: []byte{0x1b, 0x61, 0x73, 0x61},
			},
		},

		{
			name: "10;10Hzzzz",
			wantCurrentRes: &escapeParserResult{
				action: EscapeActionNone,
				args:   [][]byte{{0x7a, 0x7a, 0x7a, 0x7a}},
			},

			args: args{
				b: []byte{0x1b, 0x5b, 0x31, 0x30, 0x3b, 0x31, 0x30, 0x48, 0x7a, 0x7a, 0x7a, 0x7a},
			},
			wantRes: []EscapeSequenceResultIface{
				&escapeParserResult{
					action: EscapeActionCursorPosition,
					args:   [][]byte{{0x31, 0x30}, {0x31, 0x30}},
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
					action: EscapeActionCursorMoveToLine,
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
			name: "wrong n",
			wantCurrentRes: &escapeParserResult{
				action: EscapeActionNone,
				args:   [][]byte{{0x6e}},
			},

			args: args{
				b: []byte{0x1b, 0x5b, 0x6e},
			},
		},

		{
			name: "wrong c",
			wantCurrentRes: &escapeParserResult{
				action: EscapeActionNone,
				args:   [][]byte{{0x63}},
			},

			args: args{
				b: []byte{0x1b, 0x5b, 0x63},
			},
		},
		{
			name: "ZZc",
			wantCurrentRes: &escapeParserResult{
				action: EscapeActionNone,
				args:   [][]byte{{0x5a, 0x5a, 0x63}},
			},

			args: args{
				b: []byte{0x1b, 0x5a, 0x5a, 0x63},
			},
		},

		{
			name: "cZZ",
			wantCurrentRes: &escapeParserResult{
				action: EscapeActionNone,
				args:   [][]byte{{0x5a, 0x5a}},
			},

			args: args{
				b: []byte{0x1b, 0x63, 0x5a, 0x5a},
			},
			wantRes: []EscapeSequenceResultIface{
				&escapeParserResult{
					action: EscapeActionClearScreen,
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
					action: escapeActionEraseRightLeftScreen,
					args:   [][]byte{{0x31}},
				},
			},
		},

		{
			name: "1Jaaa",
			wantCurrentRes: &escapeParserResult{
				action: EscapeActionNone,
				args:   [][]byte{{0x61, 0x61, 0x61}},
			},

			args: args{
				b: []byte{0x1b, 0x5b, 0x31, 0x4a, 0x61, 0x61, 0x61},
			},
			wantRes: []EscapeSequenceResultIface{
				&escapeParserResult{
					action: escapeActionEraseRightLeftScreen,
					args:   [][]byte{{0x31}},
				},
			},
		},

		{
			name: "2Jaaa",
			wantCurrentRes: &escapeParserResult{
				action: EscapeActionNone,
				args:   [][]byte{{0x61, 0x61, 0x61}},
			},

			args: args{
				b: []byte{0x1b, 0x5b, 0x32, 0x4a, 0x61, 0x61, 0x61},
			},
			wantRes: []EscapeSequenceResultIface{
				&escapeParserResult{
					action: escapeActionEraseRightLeftScreen,
					args:   [][]byte{{0x32}},
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
			name:   "?1049h",
			fields: fields{},
			args: args{
				b: []byte{0x1b, 0x5b, 0x3f, 0x31, 0x30, 0x34, 0x39, 0x68},
			},
			wantRes: []EscapeSequenceResultIface{
				&escapeParserResult{
					action: escapeActionPrivateControlSequence,
					args:   [][]byte{{0x31, 0x30, 0x34, 0x39}},
				},
			},
		},

		{
			name:   "?1h",
			fields: fields{},
			args: args{
				b: []byte{0x1b, 0x3f, 0x31, 0x68},
			},
			wantRes: []EscapeSequenceResultIface{
				&escapeParserResult{
					action: escapeActionPrivateControlSequence,
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
					action: escapeActionPrivateControlSequence,
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
			name:   "ESC =",
			fields: fields{},
			args: args{
				b: []byte{0x1b, 0x3d},
			},
			wantRes: []EscapeSequenceResultIface{
				&escapeParserResult{
					action: EscapeActionSetKeypadApplicationMode,
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

		{
			name:   "7",
			fields: fields{},
			args: args{
				b: []byte{0x1b, 0x37},
			},
			wantRes: []EscapeSequenceResultIface{
				&escapeParserResult{
					action: EscapeActionSaveCursorPositionInMemory,
				},
			},
		},

		{
			name:   "8",
			fields: fields{},
			args: args{
				b: []byte{0x1b, 0x38},
			},
			wantRes: []EscapeSequenceResultIface{
				&escapeParserResult{
					action: EscapeActionRestoreCursorPositionInMemory,
				},
			},
		},

		{
			name:   "M reverse",
			fields: fields{},
			args: args{
				b: []byte{0x1b, 0x4d},
			},
			wantRes: []EscapeSequenceResultIface{
				&escapeParserResult{
					action: EscapeActionSetReverseIndex,
				},
			},
		},

		{
			name: "10;10fzzzz",
			wantCurrentRes: &escapeParserResult{
				action: EscapeActionNone,
				args:   [][]byte{{0x7a, 0x7a, 0x7a, 0x7a}},
			},

			args: args{
				b: []byte{0x1b, 0x5b, 0x31, 0x30, 0x3b, 0x31, 0x30, 0x66, 0x7a, 0x7a, 0x7a, 0x7a},
			},
			wantRes: []EscapeSequenceResultIface{
				&escapeParserResult{
					action: EscapeActionCursorPosition,
					args:   [][]byte{{0x31, 0x30}, {0x31, 0x30}},
				},
			},
		},

		{
			name:   "ESC [ 3 SP q",
			fields: fields{},
			args: args{
				b: []byte{0x1b, 0x5b, 0x33, 0x20, 0x71},
			},
			wantRes: []EscapeSequenceResultIface{
				&escapeParserResult{
					action: EscapeActionSetCursorOption,
					args:   [][]byte{{0x33}},
				},
			},
		},

		{
			name:   "ESC [ 0 SP q",
			fields: fields{},
			args: args{
				b: []byte{0x1b, 0x5b, 0x30, 0x20, 0x71},
			},
			wantRes: []EscapeSequenceResultIface{
				&escapeParserResult{
					action: EscapeActionSetCursorOption,
					args:   [][]byte{{0x30}},
				},
			},
		},

		{
			name:   "ESC [ 6 n",
			fields: fields{},
			args: args{
				b: []byte{0x1b, 0x5b, 0x36, 0x6e},
			},
			wantRes: []EscapeSequenceResultIface{
				&escapeParserResult{
					action: EscapeActionReportCursorPosition,
					args:   [][]byte{{0x36}},
				},
			},
		},

		{
			name:   "ESC [ 0 c",
			fields: fields{},
			args: args{
				b: []byte{0x1b, 0x5b, 0x30, 0x63},
			},
			wantRes: []EscapeSequenceResultIface{
				&escapeParserResult{
					action: EscapeActionReportDeviceAttr,
					args:   [][]byte{{0x30}},
				},
			},
		},

		{
			name:   "ESC [ 1;12r",
			fields: fields{},
			args: args{
				b: []byte{0x1b, 0x5b, 0x31, 0x3b, 0x31, 0x32, 0x72},
			},
			wantRes: []EscapeSequenceResultIface{
				&escapeParserResult{
					action: EscapeActionSetSrollingRegion,
					args:   [][]byte{{0x31}, {0x31, 0x32}},
				},
			},
		},

		{
			name:   "ESC (0",
			fields: fields{},
			args: args{
				b: []byte{0x1b, 0x28, 0x30},
			},
			wantRes: []EscapeSequenceResultIface{
				&escapeParserResult{
					action: EscapeActionSetEnablesDECLineDrawingMode,
				},
			},
		},

		{
			name: "ESC (B",
			args: args{
				b: []byte{0x1b, 0x28, 0x42},
			},
			wantRes: []EscapeSequenceResultIface{
				&escapeParserResult{
					action: EscapeActionSetEnablesASCIIMode,
				},
			},
		},

		{
			name: "ESC [1;159H11:2",
			wantCurrentRes: &escapeParserResult{
				action: EscapeActionNone,
				args:   [][]byte{{0x31, 0x31, 0x3a, 0x32}},
			},

			args: args{
				b: []byte{0x1b, 0x5b, 0x31, 0x3b, 0x31, 0x35, 0x39, 0x48, 0x31, 0x31, 0x3a, 0x32},
			},
			wantRes: []EscapeSequenceResultIface{
				&escapeParserResult{
					action: EscapeActionCursorPosition,
					args:   [][]byte{{0x31}, {0x31, 0x35, 0x39}},
				},
			},
		},

		{
			name: "ESC [1;159H11:29:54",
			wantCurrentRes: &escapeParserResult{
				action: EscapeActionNone,
				args:   [][]byte{{0x31, 0x31, 0x3a, 0x32, 0x39, 0x3a, 0x35, 0x34, 0x0d}},
			},

			args: args{
				b: []byte{0x1b, 0x5b, 0x31, 0x3b, 0x31, 0x35, 0x39, 0x48, 0x31, 0x31, 0x3a, 0x32, 0x39, 0x3a, 0x35, 0x34, 0x0d},
			},
			wantRes: []EscapeSequenceResultIface{
				&escapeParserResult{
					action: EscapeActionCursorPosition,
					args:   [][]byte{{0x31}, {0x31, 0x35, 0x39}},
				},
			},
		},

		{
			name: "ESC [2JProcesses: 748 total, ",
			wantCurrentRes: &escapeParserResult{
				action: EscapeActionNone,
				args: [][]byte{
					{
						0x50,
						0x72,
						0x6f,
						0x63,
						0x65,
						0x73,
						0x73,
						0x65,
						0x73,
						0x3a,
						0x20,
						0x37,
						0x34,
						0x38,
						0x20,
						0x74,
						0x6f,
						0x74,
						0x61,
						0x6c,
						0x2c,
						0x20,
					},
				},
			},

			args: args{
				b: []byte{
					0x1b,
					0x5b,
					0x32,
					0x4a,
					0x50,
					0x72,
					0x6f,
					0x63,
					0x65,
					0x73,
					0x73,
					0x65,
					0x73,
					0x3a,
					0x20,
					0x37,
					0x34,
					0x38,
					0x20,
					0x74,
					0x6f,
					0x74,
					0x61,
					0x6c,
					0x2c,
					0x20,
				},
			},
			wantRes: []EscapeSequenceResultIface{
				&escapeParserResult{
					action: escapeActionEraseRightLeftScreen,
					args:   [][]byte{{0x32}},
				},
			},
		},

		{
			name: "ESC [Hes[2JProcesses: ",
			wantCurrentRes: &escapeParserResult{
				action: EscapeActionNone,
				args: [][]byte{{
					0x50,
					0x72,
					0x6f,
					0x63,
					0x65,
					0x73,
					0x73,
					0x65,
					0x73,
					0x3a,
					0x20,
				}},
			},

			args: args{
				b: []byte{
					0x1b,
					0x5b,
					0x48,
					0x65,
					0x73,
					0x1b,
					0x5b,
					0x32,
					0x4a,
					0x50,
					0x72,
					0x6f,
					0x63,
					0x65,
					0x73,
					0x73,
					0x65,
					0x73,
					0x3a,
					0x20,
				},
			},
			wantRes: []EscapeSequenceResultIface{
				&escapeParserResult{
					action: EscapeActionCursorPosition,
				},
				&escapeParserResult{
					action: EscapeActionNone,
					args:   [][]byte{{0x65, 0x73}},
				},
				&escapeParserResult{
					action: escapeActionEraseRightLeftScreen,
					args:   [][]byte{{0x32}},
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
				t.Errorf("escapeParser.ParseEscapeSequence() len = %+v, want len %+v", len(gotRes), len(tt.wantRes))
				for _, v := range gotRes {
					t.Errorf("escapeParser.ParseEscapeSequence() got action: %v", v.GetAction())
					t.Errorf("escapeParser.ParseEscapeSequence() got len args: %d", len(v.GetArgs()))

				}
				for _, v := range tt.wantRes {
					t.Errorf("escapeParser.ParseEscapeSequence() want action: %v", v.GetAction())
					t.Errorf("escapeParser.ParseEscapeSequence() want len args: %d", len(v.GetArgs()))
				}
			}
			if e.currentResult == nil && tt.wantCurrentRes == nil {
				return
			}
			if !reflect.DeepEqual(e.currentResult, tt.wantCurrentRes) {
				t.Errorf(
					"escapeParser.ParseEscapeSequence() e.currentResult = %+v, want %+v",
					e.currentResult,
					tt.wantCurrentRes,
				)
			}
		})
	}
}
