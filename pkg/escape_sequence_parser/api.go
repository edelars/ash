package escape_sequence_parser

type EscapeSequenceParserIface interface {
	ParseEscapeSequence(b []byte) []EscapeSequenceResultIface
}

type EscapeSequenceResultIface interface {
	GetAction() EscapeAction
	GetArgs() [][]byte
	//  Valid only if Action == EscapeActionSetColor, bool true - background, false - foreground
	//	EscapeColorDefault - set foreground and background to the default color
	GetColorFormat() (EscapeColor, bool)
	// valid only if action == EscapeActionNone
	GetRaw() []byte
	// trying to parse args and get x,y from it. If fail or empty args - result will be 1,1
	// not valid for action == EscapeActionNone
	GetIntsFromArgs() (n1, n2 int)
}

type EscapeAction byte

const (
	EscapeActionNone           = 0x00
	EscapeActionCursorPosition = 0x48 // H
	EscapeActionCursorUp       = 0x41 // A
	EscapeActionCursorDown     = 0x42 // B
	EscapeActionCursorForward  = 0x43 // C
	EscapeActionCursorBackward = 0x44 // D
	EscapeActionCursorNextLine = 0x45 // E
	EscapeActionCursorPrevLine = 0x46 // F
	EscapeActionCursorLeft     = 0x47 // G
	EscapeActionCursorTop      = 0x64 // d

	EscapeActionClearScreen = 0x63 // c

	escapeActionEraseRightLeftLine = 0x4b // K
	EscapeActionEraseRightLine     = 0x51 // Q
	EscapeActionEraseLeftLine      = 0x57 // W
	EscapeActionEraseLine          = 0x59 // Y

	escapeActionEraseRightLeftScreen = 0x4a // J
	EscapeActionEraseRightScreen     = 0x5a // Z
	EscapeActionEraseLeftScreen      = 0x52 // R
	EscapeActionEraseScreen          = 0x49 // I

	escapeActionPrivateControlSequence = 0x3f //?
	EscapeActionCursorShow             = 0x68 // h
	EscapeActionCursorHide             = 0x6c // l

	EscapeActionTextInsertChar = 0x40 // "@"
	EscapeActionTextDeleteChar = 0x50 // "P"
	EscapeActionTextEraseChar  = 0x58 // "X"
	EscapeActionTextInsertLine = 0x4c // "L"
	EscapeActionTextDeleteLine = 0x4d // "M"

	EscapeActionScrollUp   = 0x53 // "S"
	EscapeActionScrollDown = 0x54 // "T"

	EscapeActionSetColor = 0x6d // "m"

	escapeActionSequenceHeader            = 0x1b // '/e'
	escapeActionControlSequenceIntroducer = 0x5b // '['
	escapeActionStringTerminator          = 0x5c // '\'
	escapeActionSemicolonDelimiter        = 0x3b // ';'

	escapeActionSetScreenMode = 0x3d // '='

)

type EscapeColor uint64

const (
	EscapeColorDefault EscapeColor = iota
	EscapeColorRed
	EscapeColorGreen
	EscapeColorYellow
	EscapeColorBlue
	EscapeColorMagenta
	EscapeColorCyan
	EscapeColorWhite
	EscapeColorBrightBlack
	EscapeColorBrightRed
	EscapeColorBrightGreen
	EscapeColorBrightYellow
	EscapeColorBrightBlue
	EscapeColorBrightMagenta
	EscapeColorBrightCyan
	EscapeColorBrightWhite

	EscapeFormatBold
	EscapeFormatUnderline
	EscapeFormatItalic

	escapeColorMaxAttr = 65536
)

func (e *escapeParser) ParseEscapeSequence(b []byte) (res []EscapeSequenceResultIface) {
	var controlSequence, brokenSequence bool

mainLool:
	for _, i := range b {
		if e.currentResult != nil && e.terminated {
			res = append(res, e.currentResult)
			e.currentResult = nil
			e.terminated = false
		}

		if brokenSequence && i != escapeActionSequenceHeader {
			e.setUpdateCurrentInputWithRaw(EscapeActionNone, i)
			continue mainLool
		}

		switch i {
		case escapeActionSequenceHeader:
			e.terminated = false
			controlSequence, brokenSequence = false, false
			e.currentResult = newEscapeParserResult(EscapeActionNone)
			continue mainLool
		case escapeActionControlSequenceIntroducer:
			controlSequence = true
			continue mainLool
		case escapeActionStringTerminator:
			e.terminated = true
			continue mainLool
		case EscapeActionCursorPosition,
			EscapeActionCursorUp,
			EscapeActionCursorDown,
			EscapeActionCursorForward,
			EscapeActionCursorBackward,
			EscapeActionCursorNextLine,
			EscapeActionCursorPrevLine,
			EscapeActionCursorLeft,
			EscapeActionCursorTop,
			escapeActionEraseRightLeftLine,
			escapeActionEraseRightLeftScreen,
			EscapeActionCursorHide,
			EscapeActionTextInsertChar,
			EscapeActionTextDeleteChar,
			EscapeActionTextEraseChar,
			EscapeActionTextInsertLine,
			EscapeActionTextDeleteLine,
			EscapeActionScrollUp,
			EscapeActionScrollDown,
			EscapeActionSetColor:
			e.setCurrentAction(EscapeAction(i))
			e.terminated = true
		case isDigit(i):
			if e.currentResult == nil {
				res = append(res, newEscapeParserResult(EscapeActionNone).WithRaw(i))
				continue mainLool
			}
			e.currentResult.addToLastArg(i)
		case EscapeActionCursorShow: // h
			if e.currentResult == nil || (e.currentResult != nil && e.currentResult.action == escapeActionPrivateControlSequence) {
				e.setCurrentAction(EscapeAction(i))
			}
			e.terminated = true

		case escapeActionSemicolonDelimiter:
			if e.currentResult != nil {
				e.currentResult.addEmptyArg()
			}
		case EscapeActionClearScreen:
			if !controlSequence {
				e.setCurrentAction(EscapeActionClearScreen)
				e.terminated = true
			} else {
				e.setUpdateCurrentInputWithRaw(EscapeActionNone, i)
				continue mainLool
			}
		case escapeActionPrivateControlSequence:
			e.setCurrentAction(escapeActionPrivateControlSequence)
			e.terminated = false
		case escapeActionSetScreenMode:
			if controlSequence {
				e.setCurrentAction(escapeActionSetScreenMode)
				e.terminated = false
			} else {
				brokenSequence = true
				e.setUpdateCurrentInputWithRaw(EscapeActionNone, i)
			}
		default:
			e.setUpdateCurrentInputWithRaw(EscapeActionNone, i)
			brokenSequence = true
			continue mainLool
		}

		if e.terminated {
			res = append(res, e.currentResult)
			e.currentResult = nil
			e.terminated = false
		}
	}
	if e.currentResult != nil {
		res = append(res, e.currentResult)
		e.currentResult = nil
	}
	return
}

func NewEscapeSequenceParser() escapeParser {
	return escapeParser{terminated: true}
}
