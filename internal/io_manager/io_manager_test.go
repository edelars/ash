package io_manager

import (
	"ash/internal/colors_adapter"
	"ash/internal/configuration"
	"ash/pkg/escape_sequence_parser"
	"ash/pkg/termbox"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type pmImpl struct{}

func (pmimpl *pmImpl) DeleteLastSymbolFromCurrentBuffer() error {
	panic("not implemented") // TODO: Implement
}

func Test_inputManager_rollScreenUp(t *testing.T) {
	pm := pmImpl{}

	h := NewInputManager(
		&pm,
		nil,
		configuration.CmdRemoveLeftSymbol,
		colors_adapter.NewColorsAdapter(configuration.Colors{}),
		13,
		10,
	)

	// y,  x
	screen := [][]termbox.Cell{
		{{Ch: 1, Fg: 1, Bg: 1}, {Ch: 1, Fg: 1, Bg: 1}, {Ch: 1, Fg: 1, Bg: 1}, {Ch: 1, Fg: 1, Bg: 1}},
		{{Ch: 2, Fg: 2, Bg: 2}, {Ch: 2, Fg: 2, Bg: 2}, {Ch: 2, Fg: 2, Bg: 2}, {Ch: 2, Fg: 2, Bg: 2}},
		{{Ch: 3, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}},
	}
	wantScreen := [][]termbox.Cell{
		{{Ch: 2, Fg: 2, Bg: 2}, {Ch: 2, Fg: 2, Bg: 2}, {Ch: 2, Fg: 2, Bg: 2}, {Ch: 2, Fg: 2, Bg: 2}},
		{{Ch: 3, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}},
		{
			{Ch: constEmptyRune, Fg: h.defaultForegroundColor, Bg: h.defaultBackgroundColor},
			{Ch: constEmptyRune, Fg: h.defaultForegroundColor, Bg: h.defaultBackgroundColor},
			{Ch: constEmptyRune, Fg: h.defaultForegroundColor, Bg: h.defaultBackgroundColor},
			{Ch: constEmptyRune, Fg: h.defaultForegroundColor, Bg: h.defaultBackgroundColor},
		},
	}
	gotScreen := [][]termbox.Cell{
		{{Ch: 0, Fg: 0, Bg: 0}, {Ch: 0, Fg: 0, Bg: 0}, {Ch: 0, Fg: 0, Bg: 0}, {Ch: 0, Fg: 0, Bg: 0}},
		{{Ch: 0, Fg: 0, Bg: 0}, {Ch: 0, Fg: 0, Bg: 0}, {Ch: 0, Fg: 0, Bg: 0}, {Ch: 0, Fg: 0, Bg: 0}},
		{{Ch: 0, Fg: 0, Bg: 0}, {Ch: 0, Fg: 0, Bg: 0}, {Ch: 0, Fg: 0, Bg: 0}, {Ch: 0, Fg: 0, Bg: 0}},
	}
	get := func(x, y int) termbox.Cell {
		return screen[y][x]
	}
	set := func(x, y int, ch rune, fg termbox.Attribute, bg termbox.Attribute) {
		gotScreen[y][x].Ch = ch
		gotScreen[y][x].Fg = fg
		gotScreen[y][x].Bg = bg
	}

	h.rollScreenUp(1, 4, 3, get, set)
	assert.Equal(t, wantScreen, gotScreen)

	wantScreen2 := [][]termbox.Cell{
		{{Ch: 3, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}},
		{
			{Ch: constEmptyRune, Fg: h.defaultForegroundColor, Bg: h.defaultBackgroundColor},
			{Ch: constEmptyRune, Fg: h.defaultForegroundColor, Bg: h.defaultBackgroundColor},
			{Ch: constEmptyRune, Fg: h.defaultForegroundColor, Bg: h.defaultBackgroundColor},
			{Ch: constEmptyRune, Fg: h.defaultForegroundColor, Bg: h.defaultBackgroundColor},
		},
		{
			{Ch: constEmptyRune, Fg: h.defaultForegroundColor, Bg: h.defaultBackgroundColor},
			{Ch: constEmptyRune, Fg: h.defaultForegroundColor, Bg: h.defaultBackgroundColor},
			{Ch: constEmptyRune, Fg: h.defaultForegroundColor, Bg: h.defaultBackgroundColor},
			{Ch: constEmptyRune, Fg: h.defaultForegroundColor, Bg: h.defaultBackgroundColor},
		},
	}

	h.rollScreenUp(2, 4, 3, get, set)
	assert.Equal(t, wantScreen2, gotScreen)
}

func Test_inputManager_Read(t *testing.T) {
	h := inputManager{
		inputEventChan: make(chan termbox.Event),
	}

	var result []byte
	var echoExecCounter int

	h.echoPrintFunc = func(p []byte) {
		echoExecCounter++
	}

	// 1 test
	str := "1234567890"
	h.enterKey = 10
	go func() {
		for _, v := range str {
			h.inputEventChan <- termbox.Event{Ch: v}
		}
		h.inputEventChan <- termbox.Event{Key: 10}
	}()

	for {
		a := make([]byte, 1)
		n, err := h.Read(a)
		assert.NoError(t, err)
		if n == 0 {
			continue
		}
		result = append(result, a[:n]...)
		if result[len(result)-1] == 0x0A {
			break
		}
	}
	assert.Equal(t, str+"\n", string(result))
	assert.Equal(t, 0, len(h.inputBuffer))
	assert.Equal(t, 11, echoExecCounter)
	echoExecCounter = 0

	result = nil

	// 2 test
	str = ""
	h.enterKey = 13
	go func() {
		for _, v := range str {
			h.inputEventChan <- termbox.Event{Ch: v}
		}
		h.inputEventChan <- termbox.Event{Key: 13}
	}()

	for {
		a := make([]byte, 1)
		n, err := h.Read(a)
		assert.NoError(t, err)
		if n == 0 {
			continue
		}
		result = append(result, a[:n]...)
		if result[len(result)-1] == 0x0A {
			break
		}
	}
	assert.Equal(t, str+"\n", string(result))
	assert.Equal(t, 0, len(h.inputBuffer))

	result = nil

	// 3 test
	str = "456"
	h.enterKey = 13
	go func() {
		for _, v := range str {
			h.inputEventChan <- termbox.Event{Ch: v}
		}
		h.inputEventChan <- termbox.Event{Key: 13}
	}()

	for {
		a := make([]byte, 10)
		n, err := h.Read(a)
		assert.NoError(t, err)
		if n == 0 {
			continue
		}
		result = append(result, a[:n]...)
		if result[len(result)-1] == 0x0A {
			break
		}
	}
	assert.Equal(t, str+"\n", string(result))
	assert.Equal(t, 0, len(h.inputBuffer))

	result = nil

	// 4 test
	str = `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nullam eget quam ac leo tempor sollicitudin. Phasellus lorem lorem, hendrerit at urna nec, sagittis sodales leo. Proin consequat orci massa, nec ultrices ligula euismod vitae. Suspendisse mollis quam non convallis bibendum. Duis vel est hendrerit, tincidunt ipsum sed, posuere mauris. Praesent quam ex, bibendum sed arcu sed, vulputate ultrices arcu. Sed at turpis a est ultricies vestibulum. Praesent eu auctor odio. Nam porta id neque vel sollicitudin. Aliquam erat volutpat. Cras nec finibus elit. Aliquam est enim, ornare quis velit a, vehicula ornare ante. Quisque aliquet metus a erat consequat, sit amet commodo neque blandit.Nam porta viverra nisi, quis elementum dolor lacinia nec. Donec condimentum nulla odio, quis varius libero dignissim eget. Duis scelerisque purus id pretium viverra. Aliquam erat volutpat. Phasellus eleifend tincidunt porta. Nam facilisis hendrerit diam sit amet blandit. Quisque ac eros ultrices, tincidunt turpis in, auctor quam. Integer malesuada erat felis, non consectetur libero porta ac. Sed in nisi ut sem interdum vulputate in non nisi. Proin lorem eros, posuere vitae fringilla quis, dapibus id mi. Praesent commodo risus sit amet augue fermentum tincidunt ut vel orci. Mauris tincidunt mattis tellus ut malesuada. Donec ut nisi fermentum, suscipit risus at, pretium eros. Fusce fermentum consequat nibh a tincidunt. Etiam ornare massa et risus ornare maximus. Sed sollicitudin congue interdum fusce.`
	h.enterKey = 13
	go func() {
		for _, v := range str {
			h.inputEventChan <- termbox.Event{Ch: v}
		}
		h.inputEventChan <- termbox.Event{Key: 13}
	}()

	for {
		a := make([]byte, 10)
		n, err := h.Read(a)
		assert.NoError(t, err)
		if n == 0 {
			continue
		}
		result = append(result, a[:n]...)
		if result[len(result)-1] == 0x0A {
			break
		}
	}
	assert.Equal(t, str+"\n", string(result))
	assert.Equal(t, 0, len(h.inputBuffer))
}

func Test_inputManager_fillScreenSquareByXYWithChar(t *testing.T) {
	pm := pmImpl{}

	h := NewInputManager(
		&pm,
		nil,
		configuration.CmdRemoveLeftSymbol,
		colors_adapter.NewColorsAdapter(configuration.Colors{}),
		13,
		10,
	)

	// y,  x
	gotScreen := [][]termbox.Cell{
		{{Ch: 0, Fg: 0, Bg: 0}, {Ch: 0, Fg: 0, Bg: 0}, {Ch: 0, Fg: 0, Bg: 0}, {Ch: 0, Fg: 0, Bg: 0}},
		{{Ch: 0, Fg: 0, Bg: 0}, {Ch: 0, Fg: 0, Bg: 0}, {Ch: 0, Fg: 0, Bg: 0}, {Ch: 0, Fg: 0, Bg: 0}},
		{{Ch: 0, Fg: 0, Bg: 0}, {Ch: 0, Fg: 0, Bg: 0}, {Ch: 0, Fg: 0, Bg: 0}, {Ch: 0, Fg: 0, Bg: 0}},
	}
	wantScreen := [][]termbox.Cell{
		{
			{Ch: 0, Fg: 0, Bg: 0},
			{Ch: 66, Fg: termbox.ColorRed, Bg: termbox.ColorGreen},
			{Ch: 66, Fg: termbox.ColorRed, Bg: termbox.ColorGreen},
			{Ch: 0, Fg: 0, Bg: 0},
		},
		{
			{Ch: 0, Fg: 0, Bg: 0},
			{Ch: 66, Fg: termbox.ColorRed, Bg: termbox.ColorGreen},
			{Ch: 66, Fg: termbox.ColorRed, Bg: termbox.ColorGreen},
			{Ch: 0, Fg: 0, Bg: 0},
		},
		{
			{Ch: 0, Fg: 0, Bg: 0},
			{Ch: 66, Fg: termbox.ColorRed, Bg: termbox.ColorGreen},
			{Ch: 66, Fg: termbox.ColorRed, Bg: termbox.ColorGreen},
			{Ch: 0, Fg: 0, Bg: 0},
		},
	}

	set := func(x, y int, ch rune, fg termbox.Attribute, bg termbox.Attribute) {
		gotScreen[y][x].Ch = ch
		gotScreen[y][x].Fg = fg
		gotScreen[y][x].Bg = bg
	}

	h.fillScreenSquareByXYWithChar(1, 2, 0, 2, 66, termbox.ColorRed, termbox.ColorGreen, set)
	assert.Equal(t, wantScreen, gotScreen)
}

func Test_is256Color(t *testing.T) {
	type args struct {
		i escape_sequence_parser.EscapeColor
	}
	tests := []struct {
		name string
		args args
		want escape_sequence_parser.EscapeColor
	}{
		{
			name: "127",
			args: args{
				i: 127,
			},
			want: 0,
		},
		{
			name: "333",
			args: args{
				i: 333,
			},
			want: 333,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := is256Color(tt.args.i); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("is256Color() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isRGBColor(t *testing.T) {
	type args struct {
		i escape_sequence_parser.EscapeColor
	}
	tests := []struct {
		name string
		args args
		want escape_sequence_parser.EscapeColor
	}{
		{
			name: "1024",
			args: args{
				i: 1024,
			},
			want: 1024,
		},
		{
			name: "22",
			args: args{
				i: 22,
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isRGBColor(tt.args.i); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("isRGBColor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_inputManager_insertEmptyLines(t *testing.T) {
	pm := pmImpl{}

	h := NewInputManager(
		&pm,
		nil,
		configuration.CmdRemoveLeftSymbol,
		colors_adapter.NewColorsAdapter(configuration.Colors{}),
		13,
		10,
	)
	h.cursorX = 1
	h.cursorY = 2
	// y,  x
	screen := [][]termbox.Cell{
		{{Ch: 1, Fg: 1, Bg: 1}, {Ch: 1, Fg: 1, Bg: 1}, {Ch: 1, Fg: 1, Bg: 1}, {Ch: 1, Fg: 1, Bg: 1}},
		{{Ch: 2, Fg: 2, Bg: 2}, {Ch: 2, Fg: 2, Bg: 2}, {Ch: 2, Fg: 2, Bg: 2}, {Ch: 2, Fg: 2, Bg: 2}},
		{{Ch: 3, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}},
	}
	wantScreen := [][]termbox.Cell{
		{{Ch: 2, Fg: 2, Bg: 2}, {Ch: 2, Fg: 2, Bg: 2}, {Ch: 2, Fg: 2, Bg: 2}, {Ch: 2, Fg: 2, Bg: 2}},
		{
			{Ch: constEmptyRune, Fg: 0, Bg: 0},
			{Ch: constEmptyRune, Fg: 0, Bg: 0},
			{Ch: constEmptyRune, Fg: 0, Bg: 0},
			{Ch: constEmptyRune, Fg: 0, Bg: 0},
		},
		{{Ch: 3, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}},
	}

	get := func(x, y int) termbox.Cell {
		return screen[y][x]
	}
	set := func(x, y int, ch rune, fg termbox.Attribute, bg termbox.Attribute) {
		screen[y][x].Ch = ch
		screen[y][x].Fg = fg
		screen[y][x].Bg = bg
	}

	h.insertEmptyLines(1, 4, 3, get, set)
	assert.Equal(t, wantScreen, screen)

	// 2 test

	h.cursorY = 4
	// y,  x
	screen2 := [][]termbox.Cell{
		{{Ch: 1, Fg: 1, Bg: 1}, {Ch: 1, Fg: 1, Bg: 1}, {Ch: 1, Fg: 1, Bg: 1}, {Ch: 1, Fg: 1, Bg: 1}},
		{{Ch: 2, Fg: 2, Bg: 2}, {Ch: 2, Fg: 2, Bg: 2}, {Ch: 2, Fg: 2, Bg: 2}, {Ch: 2, Fg: 2, Bg: 2}},
		{{Ch: 3, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}},
		{{Ch: 4, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}},
		{{Ch: 5, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}}, // cursorY
		{{Ch: 6, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}},
		{{Ch: 7, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}},
		{{Ch: 8, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}},
	}
	wantScreen2 := [][]termbox.Cell{
		{{Ch: 3, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}},
		{{Ch: 4, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}},
		{
			{Ch: constEmptyRune, Fg: 0, Bg: 0},
			{Ch: constEmptyRune, Fg: 0, Bg: 0},
			{Ch: constEmptyRune, Fg: 0, Bg: 0},
			{Ch: constEmptyRune, Fg: 0, Bg: 0},
		},
		{
			{Ch: constEmptyRune, Fg: 0, Bg: 0},
			{Ch: constEmptyRune, Fg: 0, Bg: 0},
			{Ch: constEmptyRune, Fg: 0, Bg: 0},
			{Ch: constEmptyRune, Fg: 0, Bg: 0},
		},
		{{Ch: 5, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}}, // cursorY
		{{Ch: 6, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}},
		{{Ch: 7, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}},
		{{Ch: 8, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}},
	}

	get2 := func(x, y int) termbox.Cell {
		return screen2[y][x]
	}
	set2 := func(x, y int, ch rune, fg termbox.Attribute, bg termbox.Attribute) {
		screen2[y][x].Ch = ch
		screen2[y][x].Fg = fg
		screen2[y][x].Bg = bg
	}
	h.insertEmptyLines(2, 4, 8, get2, set2)
	assert.Equal(t, wantScreen2, screen2)

	// 3 test
	h.cursorY = 1
	// y,  x
	screen3 := [][]termbox.Cell{
		{{Ch: 1, Fg: 1, Bg: 1}, {Ch: 1, Fg: 1, Bg: 1}, {Ch: 1, Fg: 1, Bg: 1}, {Ch: 1, Fg: 1, Bg: 1}},
		{{Ch: 2, Fg: 2, Bg: 2}, {Ch: 2, Fg: 2, Bg: 2}, {Ch: 2, Fg: 2, Bg: 2}, {Ch: 2, Fg: 2, Bg: 2}}, // cursorY
		{{Ch: 3, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}},
		{{Ch: 4, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}},
		{{Ch: 5, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}},
		{{Ch: 6, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}},
		{{Ch: 7, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}},
		{{Ch: 8, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}},
	}
	wantScreen3 := [][]termbox.Cell{
		{
			{Ch: constEmptyRune, Fg: 0, Bg: 0},
			{Ch: constEmptyRune, Fg: 0, Bg: 0},
			{Ch: constEmptyRune, Fg: 0, Bg: 0},
			{Ch: constEmptyRune, Fg: 0, Bg: 0},
		}, // cursorY
		{{Ch: 2, Fg: 2, Bg: 2}, {Ch: 2, Fg: 2, Bg: 2}, {Ch: 2, Fg: 2, Bg: 2}, {Ch: 2, Fg: 2, Bg: 2}},
		{{Ch: 3, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}},
		{{Ch: 4, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}},
		{{Ch: 5, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}},
		{{Ch: 6, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}},
		{{Ch: 7, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}},
		{{Ch: 8, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}},
	}

	get3 := func(x, y int) termbox.Cell {
		return screen3[y][x]
	}
	set3 := func(x, y int, ch rune, fg termbox.Attribute, bg termbox.Attribute) {
		screen3[y][x].Ch = ch
		screen3[y][x].Fg = fg
		screen3[y][x].Bg = bg
	}
	h.insertEmptyLines(5, 4, 8, get3, set3)
	assert.Equal(t, wantScreen3, screen3)
}

func Test_inputManager_deleteLines(t *testing.T) {
	pm := pmImpl{}

	h := NewInputManager(
		&pm,
		nil,
		configuration.CmdRemoveLeftSymbol,
		colors_adapter.NewColorsAdapter(configuration.Colors{}),
		13,
		10,
	)
	h.cursorX = 1
	h.cursorY = 2
	// y,  x
	screen := [][]termbox.Cell{
		{{Ch: 1, Fg: 1, Bg: 1}, {Ch: 1, Fg: 1, Bg: 1}, {Ch: 1, Fg: 1, Bg: 1}, {Ch: 1, Fg: 1, Bg: 1}},
		{{Ch: 2, Fg: 2, Bg: 2}, {Ch: 2, Fg: 2, Bg: 2}, {Ch: 2, Fg: 2, Bg: 2}, {Ch: 2, Fg: 2, Bg: 2}},
		{{Ch: 3, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}}, // cursorY
		{{Ch: 4, Fg: 4, Bg: 4}, {Ch: 4, Fg: 4, Bg: 4}, {Ch: 4, Fg: 4, Bg: 4}, {Ch: 4, Fg: 4, Bg: 4}},
		{{Ch: 5, Fg: 5, Bg: 5}, {Ch: 5, Fg: 5, Bg: 5}, {Ch: 5, Fg: 5, Bg: 5}, {Ch: 5, Fg: 5, Bg: 5}},
	}
	wantScreen := [][]termbox.Cell{
		{{Ch: 1, Fg: 1, Bg: 1}, {Ch: 1, Fg: 1, Bg: 1}, {Ch: 1, Fg: 1, Bg: 1}, {Ch: 1, Fg: 1, Bg: 1}},
		{{Ch: 2, Fg: 2, Bg: 2}, {Ch: 2, Fg: 2, Bg: 2}, {Ch: 2, Fg: 2, Bg: 2}, {Ch: 2, Fg: 2, Bg: 2}},
		{
			{Ch: constEmptyRune, Fg: 0, Bg: 0},
			{Ch: constEmptyRune, Fg: 0, Bg: 0},
			{Ch: constEmptyRune, Fg: 0, Bg: 0},
			{Ch: constEmptyRune, Fg: 0, Bg: 0},
		}, // cursorY
		{
			{Ch: constEmptyRune, Fg: 0, Bg: 0},
			{Ch: constEmptyRune, Fg: 0, Bg: 0},
			{Ch: constEmptyRune, Fg: 0, Bg: 0},
			{Ch: constEmptyRune, Fg: 0, Bg: 0},
		},
		{
			{Ch: constEmptyRune, Fg: 0, Bg: 0},
			{Ch: constEmptyRune, Fg: 0, Bg: 0},
			{Ch: constEmptyRune, Fg: 0, Bg: 0},
			{Ch: constEmptyRune, Fg: 0, Bg: 0},
		},
	}

	set := func(x, y int, ch rune, fg termbox.Attribute, bg termbox.Attribute) {
		screen[y][x].Ch = ch
		screen[y][x].Fg = fg
		screen[y][x].Bg = bg
	}

	h.deleteLines(4, 5, set)
	assert.Equal(t, wantScreen, screen)

	// 2 test
	h.cursorY = 4
	// y,  x
	screen2 := [][]termbox.Cell{
		{{Ch: 1, Fg: 1, Bg: 1}, {Ch: 1, Fg: 1, Bg: 1}, {Ch: 1, Fg: 1, Bg: 1}, {Ch: 1, Fg: 1, Bg: 1}},
		{{Ch: 2, Fg: 2, Bg: 2}, {Ch: 2, Fg: 2, Bg: 2}, {Ch: 2, Fg: 2, Bg: 2}, {Ch: 2, Fg: 2, Bg: 2}},
		{{Ch: 3, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}},
		{{Ch: 4, Fg: 4, Bg: 4}, {Ch: 4, Fg: 4, Bg: 4}, {Ch: 4, Fg: 4, Bg: 4}, {Ch: 4, Fg: 4, Bg: 4}},
		{{Ch: 5, Fg: 5, Bg: 5}, {Ch: 5, Fg: 5, Bg: 5}, {Ch: 5, Fg: 5, Bg: 5}, {Ch: 5, Fg: 5, Bg: 5}}, // cursorY
	}
	wantScreen2 := [][]termbox.Cell{
		{{Ch: 1, Fg: 1, Bg: 1}, {Ch: 1, Fg: 1, Bg: 1}, {Ch: 1, Fg: 1, Bg: 1}, {Ch: 1, Fg: 1, Bg: 1}},
		{{Ch: 2, Fg: 2, Bg: 2}, {Ch: 2, Fg: 2, Bg: 2}, {Ch: 2, Fg: 2, Bg: 2}, {Ch: 2, Fg: 2, Bg: 2}},
		{{Ch: 3, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}},
		{{Ch: 4, Fg: 4, Bg: 4}, {Ch: 4, Fg: 4, Bg: 4}, {Ch: 4, Fg: 4, Bg: 4}, {Ch: 4, Fg: 4, Bg: 4}},
		{
			{Ch: constEmptyRune, Fg: 0, Bg: 0},
			{Ch: constEmptyRune, Fg: 0, Bg: 0},
			{Ch: constEmptyRune, Fg: 0, Bg: 0},
			{Ch: constEmptyRune, Fg: 0, Bg: 0},
		}, // cursorY
	}

	set2 := func(x, y int, ch rune, fg termbox.Attribute, bg termbox.Attribute) {
		screen2[y][x].Ch = ch
		screen2[y][x].Fg = fg
		screen2[y][x].Bg = bg
	}

	h.deleteLines(4, 5, set2)
	assert.Equal(t, wantScreen2, screen2)

	// 2 test
	h.cursorY = 0
	// y,  x
	screen3 := [][]termbox.Cell{
		{{Ch: 1, Fg: 1, Bg: 1}, {Ch: 1, Fg: 1, Bg: 1}, {Ch: 1, Fg: 1, Bg: 1}, {Ch: 1, Fg: 1, Bg: 1}}, // cursorY
		{{Ch: 2, Fg: 2, Bg: 2}, {Ch: 2, Fg: 2, Bg: 2}, {Ch: 2, Fg: 2, Bg: 2}, {Ch: 2, Fg: 2, Bg: 2}},
		{{Ch: 3, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}},
		{{Ch: 4, Fg: 4, Bg: 4}, {Ch: 4, Fg: 4, Bg: 4}, {Ch: 4, Fg: 4, Bg: 4}, {Ch: 4, Fg: 4, Bg: 4}},
		{{Ch: 5, Fg: 5, Bg: 5}, {Ch: 5, Fg: 5, Bg: 5}, {Ch: 5, Fg: 5, Bg: 5}, {Ch: 5, Fg: 5, Bg: 5}},
	}
	wantScreen3 := [][]termbox.Cell{
		{
			{Ch: constEmptyRune, Fg: 0, Bg: 0},
			{Ch: constEmptyRune, Fg: 0, Bg: 0},
			{Ch: constEmptyRune, Fg: 0, Bg: 0},
			{Ch: constEmptyRune, Fg: 0, Bg: 0},
		}, // cursorY
		{
			{Ch: constEmptyRune, Fg: 0, Bg: 0},
			{Ch: constEmptyRune, Fg: 0, Bg: 0},
			{Ch: constEmptyRune, Fg: 0, Bg: 0},
			{Ch: constEmptyRune, Fg: 0, Bg: 0},
		},
		{
			{Ch: constEmptyRune, Fg: 0, Bg: 0},
			{Ch: constEmptyRune, Fg: 0, Bg: 0},
			{Ch: constEmptyRune, Fg: 0, Bg: 0},
			{Ch: constEmptyRune, Fg: 0, Bg: 0},
		},
		{
			{Ch: constEmptyRune, Fg: 0, Bg: 0},
			{Ch: constEmptyRune, Fg: 0, Bg: 0},
			{Ch: constEmptyRune, Fg: 0, Bg: 0},
			{Ch: constEmptyRune, Fg: 0, Bg: 0},
		},
		{
			{Ch: constEmptyRune, Fg: 0, Bg: 0},
			{Ch: constEmptyRune, Fg: 0, Bg: 0},
			{Ch: constEmptyRune, Fg: 0, Bg: 0},
			{Ch: constEmptyRune, Fg: 0, Bg: 0},
		},
	}

	set3 := func(x, y int, ch rune, fg termbox.Attribute, bg termbox.Attribute) {
		screen3[y][x].Ch = ch
		screen3[y][x].Fg = fg
		screen3[y][x].Bg = bg
	}

	h.deleteLines(4, 5, set3)
	assert.Equal(t, wantScreen3, screen3)
}

func Test_inputManager_moveCursorAndCleanBetwen(t *testing.T) {
	pm := pmImpl{}

	h := NewInputManager(
		&pm,
		nil,
		configuration.CmdRemoveLeftSymbol,
		colors_adapter.NewColorsAdapter(configuration.Colors{}),
		13,
		10,
	)
	h.cursorX = 3
	h.cursorY = 4
	// y,  x
	screen := [][]termbox.Cell{
		{{Ch: 1, Fg: 1, Bg: 1}, {Ch: 1, Fg: 1, Bg: 1}, {Ch: 1, Fg: 1, Bg: 1}, {Ch: 1, Fg: 1, Bg: 1}},
		{{Ch: 2, Fg: 2, Bg: 2}, {Ch: 2, Fg: 2, Bg: 2}, {Ch: 2, Fg: 2, Bg: 2}, {Ch: 2, Fg: 2, Bg: 2}},
		{{Ch: 3, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}},
		{{Ch: 4, Fg: 4, Bg: 4}, {Ch: 4, Fg: 4, Bg: 4}, {Ch: 4, Fg: 4, Bg: 4}, {Ch: 4, Fg: 4, Bg: 4}},
		{{Ch: 5, Fg: 5, Bg: 5}, {Ch: 5, Fg: 5, Bg: 5}, {Ch: 5, Fg: 5, Bg: 5}, {Ch: 5, Fg: 5, Bg: 5}},
	}
	wantScreen := [][]termbox.Cell{
		{{Ch: 1, Fg: 1, Bg: 1}, {Ch: 1, Fg: 1, Bg: 1}, {Ch: 1, Fg: 1, Bg: 1}, {Ch: 1, Fg: 1, Bg: 1}},
		{{Ch: 2, Fg: 2, Bg: 2}, {Ch: 2, Fg: 2, Bg: 2}, {Ch: constEmptyRune, Fg: 0, Bg: 0}, {Ch: constEmptyRune, Fg: 0, Bg: 0}},
		{
			{Ch: constEmptyRune, Fg: 0, Bg: 0},
			{Ch: constEmptyRune, Fg: 0, Bg: 0},
			{Ch: constEmptyRune, Fg: 0, Bg: 0},
			{Ch: constEmptyRune, Fg: 0, Bg: 0},
		},
		{
			{Ch: constEmptyRune, Fg: 0, Bg: 0},
			{Ch: constEmptyRune, Fg: 0, Bg: 0},
			{Ch: constEmptyRune, Fg: 0, Bg: 0},
			{Ch: constEmptyRune, Fg: 0, Bg: 0},
		},
		{
			{Ch: constEmptyRune, Fg: 0, Bg: 0},
			{Ch: constEmptyRune, Fg: 0, Bg: 0},
			{Ch: constEmptyRune, Fg: 0, Bg: 0},
			{Ch: constEmptyRune, Fg: 0, Bg: 0},
		},
	}

	set := func(x, y int, ch rune, fg termbox.Attribute, bg termbox.Attribute) {
		screen[y][x].Ch = ch
		screen[y][x].Fg = fg
		screen[y][x].Bg = bg
	}

	h.moveCursorAndCleanBetwen(2, 1, 4, 5, set)
	assert.Equal(t, wantScreen, screen)
	assert.Equal(t, 2, h.cursorX)
	assert.Equal(t, 1, h.cursorY)

	// 2 test
	h.cursorX = 1
	h.cursorY = 0
	// y,  x
	screen2 := [][]termbox.Cell{
		{{Ch: 1, Fg: 1, Bg: 1}, {Ch: 1, Fg: 1, Bg: 1}, {Ch: 1, Fg: 1, Bg: 1}, {Ch: 1, Fg: 1, Bg: 1}},
		{{Ch: 2, Fg: 2, Bg: 2}, {Ch: 2, Fg: 2, Bg: 2}, {Ch: 2, Fg: 2, Bg: 2}, {Ch: 2, Fg: 2, Bg: 2}},
		{{Ch: 3, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}},
		{{Ch: 4, Fg: 4, Bg: 4}, {Ch: 4, Fg: 4, Bg: 4}, {Ch: 4, Fg: 4, Bg: 4}, {Ch: 4, Fg: 4, Bg: 4}},
		{{Ch: 5, Fg: 5, Bg: 5}, {Ch: 5, Fg: 5, Bg: 5}, {Ch: 5, Fg: 5, Bg: 5}, {Ch: 5, Fg: 5, Bg: 5}},
	}
	wantScreen2 := [][]termbox.Cell{
		{
			{
				Ch: 1,
				Fg: 1,
				Bg: 1,
			}, {Ch: constEmptyRune, Fg: 0, Bg: 0}, {Ch: constEmptyRune, Fg: 0, Bg: 0}, {Ch: constEmptyRune, Fg: 0, Bg: 0},
		},
		{
			{Ch: constEmptyRune, Fg: 0, Bg: 0},
			{Ch: constEmptyRune, Fg: 0, Bg: 0},
			{Ch: constEmptyRune, Fg: 0, Bg: 0},
			{Ch: constEmptyRune, Fg: 0, Bg: 0},
		},

		{
			{Ch: constEmptyRune, Fg: 0, Bg: 0},
			{Ch: constEmptyRune, Fg: 0, Bg: 0},
			{Ch: constEmptyRune, Fg: 0, Bg: 0},
			{Ch: constEmptyRune, Fg: 0, Bg: 0},
		},
		{
			{Ch: constEmptyRune, Fg: 0, Bg: 0},
			{Ch: constEmptyRune, Fg: 0, Bg: 0},
			{Ch: constEmptyRune, Fg: 0, Bg: 0},
			{Ch: constEmptyRune, Fg: 0, Bg: 0},
		},
		{
			{Ch: constEmptyRune, Fg: 0, Bg: 0},
			{Ch: constEmptyRune, Fg: 0, Bg: 0},
			{Ch: constEmptyRune, Fg: 0, Bg: 0},
			{Ch: constEmptyRune, Fg: 0, Bg: 0},
		},
	}

	set2 := func(x, y int, ch rune, fg termbox.Attribute, bg termbox.Attribute) {
		screen2[y][x].Ch = ch
		screen2[y][x].Fg = fg
		screen2[y][x].Bg = bg
	}

	h.moveCursorAndCleanBetwen(0, 4, 4, 5, set2)
	assert.Equal(t, wantScreen2, screen2)
	assert.Equal(t, 0, h.cursorX)
	assert.Equal(t, 4, h.cursorY)
}

func Test_inputManager_generateEscapeSeqCursorPosition(t *testing.T) {
	type fields struct {
		cursorX int
		cursorY int
	}
	tests := []struct {
		name   string
		fields fields
		want   []rune
	}{
		{
			name: "39",
			fields: fields{
				cursorX: 2,
				cursorY: 8,
			},
			want: []rune{0x1b, 0x5b, 0x39, 0x3B, 0x33, 0x52},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &inputManager{
				cursorX: tt.fields.cursorX,
				cursorY: tt.fields.cursorY,
			}
			if got := i.generateEscapeSeqCursorPosition(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("inputManager.generateEscapeSeqCursorPosition() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_inputManager_generateEscapeSeqDeviceAttr(t *testing.T) {
	tests := []struct {
		name string
		want []rune
	}{
		{
			name: "1",
			want: []rune{0x1b, 0x5b, 0x3F, 0x31, 0x3B, 0x32, 0x63},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &inputManager{}
			if got := i.generateEscapeSeqDeviceAttr(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("inputManager.generateEscapeSeqDeviceAttr() = %v, want %v", got, tt.want)
			}
		})
	}
}
