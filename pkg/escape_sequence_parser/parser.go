package escape_sequence_parser

type escapeParser struct {
	terminated    bool
	currentResult *escapeParserResult
}

func newEscapeParserResult(action EscapeAction) *escapeParserResult {
	return &escapeParserResult{
		action: action,
	}
}

type escapeParserResult struct {
	action EscapeAction
	args   [][]byte
}

func (e *escapeParserResult) GetAction() EscapeAction {
	switch e.action {

	case escapeActionEraseRightLeftLine:
		if len(e.args) == 0 {
			return EscapeActionEraseRight
		} else if len(e.args[0]) > 0 {
			if e.args[0][0] == 0x31 { // 1
				return EscapeActionEraseLeft
			} else if e.args[0][0] == 0x32 { // 2
				return EscapeActionEraseLine
			}
		}
		return EscapeActionNone

	case escapeActionEraseDownUpScreen:
		if len(e.args) == 0 {
			return EscapeActionEraseDown
		} else if len(e.args[0]) > 0 {
			if e.args[0][0] == 0x31 { // 1
				return EscapeActionEraseUp
			} else if e.args[0][0] == 0x32 { // 2
				return EscapeActionEraseScreen
			}
		}
		return EscapeActionNone

	}
	return e.action
}

func (e *escapeParserResult) GetArgs() []string {
	panic("not implemented") // TODO: Implement
}

func (e *escapeParserResult) WithRaw(b byte) *escapeParserResult {
	e.args = append(e.args, []byte{b})
	return e
}

func (e *escapeParserResult) addToLastArg(b byte) {
	if len(e.args) == 0 {
		e.args = append(e.args, []byte{b})
	} else {
		e.args[len(e.args)-1] = append(e.args[len(e.args)-1], b)
	}
}

func (e *escapeParserResult) addEmptyArg() {
	e.args = append(e.args, []byte{})
}

const (
	eTypeSequenceHeader            = 0x1b // '/e'
	eTypeControlSequenceIntroducer = 0x5b // '['
	eTypeStringTerminator          = 0x5c // '\'
	eTypeSemicolonDelimiter        = 0x3b // ';'

)

func (e *escapeParser) setCurrentAction(a EscapeAction) {
	if e.currentResult == nil {
		e.currentResult = newEscapeParserResult(a)
	} else {
		e.currentResult.action = a
	}
}

// return same if its digit 0-9
func isDigit(i byte) byte {
	if i >= 0x30 && i <= 0x39 {
		return i
	} else {
		return 0x00
	}
}
