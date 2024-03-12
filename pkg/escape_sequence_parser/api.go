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
	EscapeActionNone             = 0x00
	EscapeActionCursorPosition   = 0x48 // H
	EscapeActionCursorUp         = 0x41 // A
	EscapeActionCursorDown       = 0x42 // B
	EscapeActionCursorForward    = 0x43 // C
	EscapeActionCursorBackward   = 0x44 // D
	EscapeActionCursorNextLine   = 0x45 // E
	EscapeActionCursorPrevLine   = 0x46 // F
	EscapeActionCursorLeft       = 0x47 // G
	EscapeActionCursorMoveToLine = 0x64 // d

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

	EscapeActionSequenceHeader            = 0x1b // '/e'
	EscapeActionControlSequenceIntroducer = 0x5b // '['
	escapeActionStringTerminator          = 0x5c // '\'
	escapeActionSemicolonDelimiter        = 0x3b // ';'

	escapeActionSetScreenMode = 0x3d // '='

	EscapeActionSaveCursorPositionInMemory    = 0x37 // '7'
	EscapeActionRestoreCursorPositionInMemory = 0x38 // '8'
	EscapeActionSetReverseIndex               = 0x39 // '9'

	escapeActionHVCursorPosition = 0x66 // f
	escapeActionSP               = 0x20 // SP is a literal space character (0x20)
	EscapeActionSetCursorOption  = 0x71 // q

	EscapeActionReportCursorPosition = 0x6e // n
	EscapeActionReportDeviceAttr     = 0x30 // 0

	EscapeActionSetSrollingRegion        = 0x72 // r
	EscapeActionSetKeypadApplicationMode = 0x31 // 1

	EscapeActionSetEnablesASCIIMode          = 0x32 // 2
	EscapeActionSetEnablesDECLineDrawingMode = 0x30 // 0
	escapeActionBeginCharacterSet            = 0x28 // (

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
	var controlSequence, brokenSequence, spSequence, beginCharacterSet bool

mainLool:
	for _, i := range b {
		if e.currentResult != nil && e.terminated {
			res = append(res, e.currentResult)
			e.currentResult = nil
			e.terminated = false
		}

		if brokenSequence && i != EscapeActionSequenceHeader {
			e.setUpdateCurrentInputWithRaw(EscapeActionNone, i)
			continue mainLool
		}

		switch i {
		case EscapeActionSequenceHeader:
			e.terminated = false
			controlSequence, brokenSequence, spSequence, beginCharacterSet = false, false, false, false
			e.currentResult = newEscapeParserResult(EscapeActionNone)
			continue mainLool
		case EscapeActionControlSequenceIntroducer:
			controlSequence = true
			continue mainLool
		case escapeActionStringTerminator:
			e.terminated = true
			continue mainLool
		case EscapeActionCursorPosition,
			EscapeActionCursorUp,
			EscapeActionCursorForward,
			EscapeActionCursorBackward,
			EscapeActionCursorNextLine,
			EscapeActionCursorPrevLine,
			EscapeActionCursorLeft,
			EscapeActionCursorMoveToLine,
			escapeActionEraseRightLeftLine,
			escapeActionEraseRightLeftScreen,
			EscapeActionCursorHide,
			EscapeActionTextInsertChar,
			EscapeActionTextDeleteChar,
			EscapeActionTextEraseChar,
			EscapeActionTextInsertLine,
			EscapeActionScrollUp,
			EscapeActionScrollDown,
			EscapeActionSetColor,
			EscapeActionSetSrollingRegion:
			e.setCurrentAction(EscapeAction(i))
			e.terminated = true
		case escapeActionHVCursorPosition:
			e.setCurrentAction(EscapeActionCursorPosition)
			e.terminated = true
		case EscapeActionCursorDown:
			if beginCharacterSet {
				e.setCurrentAction(EscapeActionSetEnablesASCIIMode)
			} else {
				e.setCurrentAction(EscapeActionCursorDown)
			}
			e.terminated = true
		case EscapeActionTextDeleteLine:
			if controlSequence == false && e.currentResult != nil {
				e.setCurrentAction(EscapeActionSetReverseIndex)
			} else {
				e.setCurrentAction(EscapeAction(i))
			}
			e.terminated = true

		case isDigit(i):
			if e.currentResult == nil {
				// res = append(res, newEscapeParserResult(EscapeActionNone).WithRaw(i))
				e.setUpdateCurrentInputWithRaw(EscapeActionNone, i)

				continue mainLool
			}
			if !controlSequence && e.currentResult != nil {
				switch i {
				case 0x37: // 7
					e.setCurrentAction(EscapeActionSaveCursorPositionInMemory)
					e.terminated = true
					continue mainLool
				case 0x38: // 7
					e.setCurrentAction(EscapeActionRestoreCursorPositionInMemory)
					e.terminated = true
					continue mainLool
				}
			}
			if !controlSequence && beginCharacterSet && i == 0x30 {
				e.setCurrentAction(EscapeActionSetEnablesDECLineDrawingMode)
				e.terminated = true
				continue mainLool
			}

			e.currentResult.addToLastArg(i)
		case EscapeActionCursorShow: // h
			// if e.currentResult == nil || (e.currentResult != nil && e.currentResult.action == escapeActionPrivateControlSequence) {
			if e.currentResult == nil {
				e.setCurrentAction(EscapeActionCursorShow)
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
				continue mainLool
			}
			if e.currentResult != nil && len(e.currentResult.GetRaw()) == 1 && e.currentResult.GetRaw()[0] == 0x30 {
				e.setCurrentAction(EscapeActionReportDeviceAttr)
				e.terminated = true
				continue mainLool
			}
			e.setUpdateCurrentInputWithRaw(EscapeActionNone, i)
			continue mainLool
		case escapeActionPrivateControlSequence:
			e.setCurrentAction(escapeActionPrivateControlSequence)
			e.terminated = false
		case escapeActionSetScreenMode:
			if controlSequence {
				e.setCurrentAction(escapeActionSetScreenMode)
				e.terminated = false
			} else if !controlSequence && !brokenSequence {
				e.setCurrentAction(EscapeActionSetKeypadApplicationMode)
				e.terminated = false
			} else {
				brokenSequence = true
				e.setUpdateCurrentInputWithRaw(EscapeActionNone, i)
			}
		case escapeActionSP:
			spSequence = true
		case EscapeActionSetCursorOption:
			if !spSequence {
				e.setUpdateCurrentInputWithRaw(EscapeActionNone, i)
				brokenSequence = true
				continue mainLool
			}
			e.setCurrentAction(EscapeActionSetCursorOption)
			e.terminated = true
		case EscapeActionReportCursorPosition:
			if !controlSequence || e.currentResult == nil || len(e.currentResult.GetRaw()) != 1 || e.currentResult.GetRaw()[0] != 0x36 {
				e.setUpdateCurrentInputWithRaw(EscapeActionNone, i)
				brokenSequence = true
				continue mainLool
			}
			e.setCurrentAction(EscapeActionReportCursorPosition)
			e.terminated = true
		case escapeActionBeginCharacterSet:
			beginCharacterSet = true
		default:
			if e.currentResult != nil {
				e.currentResult.addToLastArg(i)
			} else {
				e.setUpdateCurrentInputWithRaw(EscapeActionNone, i)
			}
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
