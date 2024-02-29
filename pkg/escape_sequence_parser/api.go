package escape_sequence_parser

type EscapeSequenceParserIface interface {
	ParseEscapeSequence(b []byte) EscapeSequenceResultIface
}

type EscapeSequenceResultIface interface {
	GetAction() EscapeAction
	GetArgs() []string
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
	EscapeActionEraseRight         = 0x51 // Q
	EscapeActionEraseLeft          = 0x57 // W
	EscapeActionEraseLine          = 0x58 // X

	escapeActionEraseDownUpScreen = 0x4a // J
	EscapeActionEraseDown         = 0x5a // Z
	EscapeActionEraseUp           = 0x52 // R
	EscapeActionEraseScreen       = 0x4c // L

	// EscapeActionEraseRight         = 0x4b // K

)

func (e *escapeParser) ParseEscapeSequence(b []byte) (res []EscapeSequenceResultIface) {
	var controlSequence, brokenSequence bool

mainLool:
	for _, i := range b {
		if e.currentResult != nil && e.terminated {
			res = append(res, e.currentResult)
			e.currentResult = nil
		}

		if brokenSequence && i != eTypeSequenceHeader {
			res = append(res, newEscapeParserResult(EscapeActionNone).WithRaw(i))
			continue mainLool
		}

		switch i {
		case eTypeSequenceHeader:
			e.terminated = false
			controlSequence, brokenSequence = false, false
			e.currentResult = newEscapeParserResult(EscapeActionNone)
			continue mainLool
		case eTypeControlSequenceIntroducer:
			controlSequence = true
			continue mainLool
		case eTypeStringTerminator:
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
			escapeActionEraseDownUpScreen:
			e.setCurrentAction(EscapeAction(i))
			e.terminated = true
		case isDigit(i):
			if e.currentResult == nil {
				res = append(res, newEscapeParserResult(EscapeActionNone).WithRaw(i))
				continue mainLool
			}
			e.currentResult.addToLastArg(i)
		case eTypeSemicolonDelimiter:
			if e.currentResult != nil {
				e.currentResult.addEmptyArg()
			}
		case EscapeActionClearScreen:
			if !controlSequence {
				e.setCurrentAction(EscapeActionClearScreen)
				e.terminated = true
			} else {
				res = append(res, newEscapeParserResult(EscapeActionNone).WithRaw(i))
				continue mainLool
			}
		default:
			res = append(res, newEscapeParserResult(EscapeActionNone).WithRaw(i))
			brokenSequence = true
			continue mainLool
		}

		if e.terminated {
			res = append(res, e.currentResult)
			e.currentResult = nil
		}
	}

	return
}

func NewEscapeSequenceParser() escapeParser {
	return escapeParser{terminated: true}
}
