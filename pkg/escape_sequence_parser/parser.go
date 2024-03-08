package escape_sequence_parser

import (
	"strconv"
)

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
			return EscapeActionEraseRightLine
		} else if len(e.args[0]) > 0 {
			if e.args[0][0] == 0x31 { // 1
				return EscapeActionEraseLeftLine
			} else if e.args[0][0] == 0x32 { // 2
				return EscapeActionEraseLine
			}
		}
		return EscapeActionNone

	case escapeActionEraseRightLeftScreen:
		if len(e.args) == 0 {
			return EscapeActionEraseRightScreen
		} else if len(e.args[0]) > 0 {
			if e.args[0][0] == 0x31 { // 1
				return EscapeActionEraseLeftScreen
			} else if e.args[0][0] == 0x32 { // 2
				return EscapeActionEraseScreen
			}
		}
		return EscapeActionNone

	case EscapeActionCursorHide:
		if len(e.args) == 1 && len(e.args[0]) == 2 && e.args[0][0] == 0x32 && e.args[0][1] == 0x35 {
			return EscapeActionCursorHide
		}
		return escapeActionPrivateControlSequence

	case EscapeActionCursorShow:
		if len(e.args) == 1 && len(e.args[0]) == 2 && e.args[0][0] == 0x32 && e.args[0][1] == 0x35 {
			return EscapeActionCursorShow
		}
		return escapeActionPrivateControlSequence

	}
	return e.action
}

func (e *escapeParserResult) GetArgs() [][]byte {
	return e.args
}

// trying to parse args and get x,y from it. If fail or empty args - result will be 1,1
// not valid for action == EscapeActionNone
func (e *escapeParserResult) GetIntsFromArgs() (x, y int) {
	x, y = 1, 1
	if len(e.args) < 1 || len(e.args[0]) == 0 {
		return
	}

	var s string
	for _, v := range e.args[0] {
		s = s + string(v)
	}
	value, err := strconv.ParseInt(s, 10, 64)
	if err == nil {
		x = int(value)
	}

	if len(e.args) < 2 || len(e.args[1]) == 0 {
		return
	}

	s = ""
	for _, v := range e.args[1] {
		s = s + string(v)
	}
	value, err = strconv.ParseInt(s, 10, 64)
	if err == nil {
		y = int(value)
	}

	return
}

func (e *escapeParserResult) WithRaw(b byte) *escapeParserResult {
	e.args = append(e.args, []byte{b})
	return e
}

// valid only if action == EscapeActionNone
func (e *escapeParserResult) GetRaw() []byte {
	if len(e.args) == 0 || len(e.args[0]) == 0 {
		return nil
	}
	return e.args[0]
}

func (e *escapeParser) setUpdateCurrentInputWithRaw(action EscapeAction, i byte) {
	if e.currentResult == nil {
		e.currentResult = newEscapeParserResult(action)
	}
	e.currentResult = e.currentResult.WithRaw(i)
}

// Valid only if Action == EscapeActionSetColor, bool true - background, false - foreground
// EscapeColorDefault - set foreground and background to the default color
// if return color > 256 and < 513 -256 color result
// if return color > 513  - RGB color result. termbox.Attribute color format
func (e *escapeParserResult) GetColorFormat() (EscapeColor, bool) {
	switch len(e.args) {
	case 1:
		if len(e.args[0]) == 1 && e.args[0][0] == 0x30 { // 0
			return EscapeColorDefault, false
		}
		if len(e.args[0]) > 2 || len(e.args[0]) == 0 {
			return EscapeColorDefault, false
		}
		var isBack bool
		if e.args[0][0] == 0x34 { // 4
			isBack = true
		} else if e.args[0][0] != 0x33 { // 3
			return EscapeColorDefault, false
		}

		switch e.args[0][1] {
		case 0x30:
			return EscapeColorDefault, isBack
		case 0x31:
			return EscapeColorRed, isBack
		case 0x32:
			return EscapeColorGreen, isBack
		case 0x33:
			return EscapeColorYellow, isBack
		case 0x34:
			return EscapeColorBlue, isBack
		case 0x35:
			return EscapeColorMagenta, isBack
		case 0x36:
			return EscapeColorCyan, isBack
		case 0x37:
			return EscapeColorWhite, isBack
		default:
			return EscapeColorDefault, false
		}
	case 2:
		if len(e.args[0]) != 2 || len(e.args[1]) != 1 || e.args[1][0] != 0x31 { // 1
			return EscapeColorDefault, false
		}
		var isBack bool
		if e.args[0][0] == 0x34 { // 4
			isBack = true
		} else if e.args[0][0] != 0x33 { // 3
			return EscapeColorDefault, false
		}

		switch e.args[0][1] {
		case 0x30:
			return EscapeColorBrightBlack, isBack
		case 0x31:
			return EscapeColorBrightRed, isBack
		case 0x32:
			return EscapeColorBrightGreen, isBack
		case 0x33:
			return EscapeColorBrightYellow, isBack
		case 0x34:
			return EscapeColorBrightBlue, isBack
		case 0x35:
			return EscapeColorBrightMagenta, isBack
		case 0x36:
			return EscapeColorBrightCyan, isBack
		case 0x37:
			return EscapeColorBrightWhite, isBack
		default:
			return EscapeColorDefault, false
		}
	case 3: // 256 color
		if len(e.args[0]) != 2 || len(e.args[1]) != 1 || len(e.args[2]) > 3 || len(e.args[2]) > 8 || len(e.args[2]) == 0 || e.args[1][0] != 0x35 { // 5
			return EscapeColorDefault, false
		}
		var isBack bool
		if e.args[0][0] == 0x34 { // 4
			isBack = true
		} else if e.args[0][0] != 0x33 { // 3
			return EscapeColorDefault, false
		}

		var s string
		for _, v := range e.args[2] {
			s = s + string(v)
		}
		value, _ := strconv.ParseInt(s, 10, 64)
		return EscapeColor(value) + 256, isBack

	case 5: // RGB color
		if len(e.args[0]) != 2 || len(e.args[1]) != 1 || len(e.args[2]) > 3 || len(e.args[2]) == 0 || e.args[1][0] != 0x32 { // 2
			return EscapeColorDefault, false
		}
		var isBack bool
		if e.args[0][0] == 0x34 { // 4
			isBack = true
		} else if e.args[0][0] != 0x33 { // 3
			return EscapeColorDefault, false
		}

		var r string
		for _, v := range e.args[2] {
			r = r + string(v)
		}
		valueR, _ := strconv.ParseInt(r, 10, 8)

		var g string
		for _, v := range e.args[3] {
			g = g + string(v)
		}
		valueG, _ := strconv.ParseInt(g, 10, 8)

		var b string
		for _, v := range e.args[4] {
			b = b + string(v)
		}
		valueB, _ := strconv.ParseInt(b, 10, 8)

		value := rGBToAttribute(uint8(valueR), uint8(valueG), uint8(valueB))
		return EscapeColor(value), isBack
	default:
		return EscapeColorDefault, false
	}
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

// RGBToAttribute is used to convert an rgb triplet into a termbox attribute.
// This attribute can only be applied when termbox is in Full RGB mode,
// otherwise it'll be ignored and no color will be drawn.
// R, G, B have to be in the range of 0 and 255.
func rGBToAttribute(r uint8, g uint8, b uint8) uint64 {
	var color uint64 = uint64(b)
	color += uint64(g) << 8
	color += uint64(r) << 16
	color += 1 << 25
	color = color * uint64(escapeColorMaxAttr)
	// Left-shift back to the place where rgb is stored.
	return color
}
