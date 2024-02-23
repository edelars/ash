package io_manager

import (
	"testing"

	"ash/internal/colors_adapter"
	"ash/internal/configuration"
	"ash/pkg/termbox"

	"github.com/stretchr/testify/assert"
)

type pmImpl struct{}

func (pmimpl *pmImpl) DeleteLastSymbolFromCurrentBuffer() error {
	panic("not implemented") // TODO: Implement
}

func Test_inputManager_rollScreenUp(t *testing.T) {
	pm := pmImpl{}

	h := NewInputManager(&pm, configuration.CmdRemoveLeftSymbol, colors_adapter.NewColorsAdapter(configuration.Colors{}), 13)

	// y,  x
	screen := [][]termbox.Cell{
		{{Ch: 1, Fg: 1, Bg: 1}, {Ch: 1, Fg: 1, Bg: 1}, {Ch: 1, Fg: 1, Bg: 1}, {Ch: 1, Fg: 1, Bg: 1}},
		{{Ch: 2, Fg: 2, Bg: 2}, {Ch: 2, Fg: 2, Bg: 2}, {Ch: 2, Fg: 2, Bg: 2}, {Ch: 2, Fg: 2, Bg: 2}},
		{{Ch: 3, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}},
	}
	wantScreen := [][]termbox.Cell{
		{{Ch: 2, Fg: 2, Bg: 2}, {Ch: 2, Fg: 2, Bg: 2}, {Ch: 2, Fg: 2, Bg: 2}, {Ch: 2, Fg: 2, Bg: 2}},
		{{Ch: 3, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}, {Ch: 3, Fg: 3, Bg: 3}},
		{{Ch: constEmptyRune, Fg: h.defaultForegroundColor, Bg: h.defaultBackgroundColor}, {Ch: constEmptyRune, Fg: h.defaultForegroundColor, Bg: h.defaultBackgroundColor}, {Ch: constEmptyRune, Fg: h.defaultForegroundColor, Bg: h.defaultBackgroundColor}, {Ch: constEmptyRune, Fg: h.defaultForegroundColor, Bg: h.defaultBackgroundColor}},
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
		{{Ch: constEmptyRune, Fg: h.defaultForegroundColor, Bg: h.defaultBackgroundColor}, {Ch: constEmptyRune, Fg: h.defaultForegroundColor, Bg: h.defaultBackgroundColor}, {Ch: constEmptyRune, Fg: h.defaultForegroundColor, Bg: h.defaultBackgroundColor}, {Ch: constEmptyRune, Fg: h.defaultForegroundColor, Bg: h.defaultBackgroundColor}},
		{{Ch: constEmptyRune, Fg: h.defaultForegroundColor, Bg: h.defaultBackgroundColor}, {Ch: constEmptyRune, Fg: h.defaultForegroundColor, Bg: h.defaultBackgroundColor}, {Ch: constEmptyRune, Fg: h.defaultForegroundColor, Bg: h.defaultBackgroundColor}, {Ch: constEmptyRune, Fg: h.defaultForegroundColor, Bg: h.defaultBackgroundColor}},
	}

	h.rollScreenUp(2, 4, 3, get, set)
	assert.Equal(t, wantScreen2, gotScreen)
}

func Test_inputManager_Read(t *testing.T) {
	h := inputManager{
		inputEventChan: make(chan termbox.Event),
	}

	var result []byte

	// 1 test
	str := "1234567890"
	go func() {
		for _, v := range str {
			h.inputEventChan <- termbox.Event{Ch: v}
		}
		h.inputEventChan <- termbox.Event{Key: '\n'}
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
	assert.Equal(t, str+string('\n'), string(result))
	assert.Equal(t, 0, len(h.inputBuffer))

	result = nil

	// 2 test
	str = ""
	go func() {
		for _, v := range str {
			h.inputEventChan <- termbox.Event{Ch: v}
		}
		h.inputEventChan <- termbox.Event{Key: '\n'}
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
	assert.Equal(t, str+string('\n'), string(result))
	assert.Equal(t, 0, len(h.inputBuffer))

	result = nil

	// 3 test
	str = "456"
	go func() {
		for _, v := range str {
			h.inputEventChan <- termbox.Event{Ch: v}
		}
		h.inputEventChan <- termbox.Event{Key: '\n'}
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
	assert.Equal(t, str+string('\n'), string(result))
	assert.Equal(t, 0, len(h.inputBuffer))

	result = nil

	// 4 test
	str = `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nullam eget quam ac leo tempor sollicitudin. Phasellus lorem lorem, hendrerit at urna nec, sagittis sodales leo. Proin consequat orci massa, nec ultrices ligula euismod vitae. Suspendisse mollis quam non convallis bibendum. Duis vel est hendrerit, tincidunt ipsum sed, posuere mauris. Praesent quam ex, bibendum sed arcu sed, vulputate ultrices arcu. Sed at turpis a est ultricies vestibulum. Praesent eu auctor odio. Nam porta id neque vel sollicitudin. Aliquam erat volutpat. Cras nec finibus elit. Aliquam est enim, ornare quis velit a, vehicula ornare ante. Quisque aliquet metus a erat consequat, sit amet commodo neque blandit.Nam porta viverra nisi, quis elementum dolor lacinia nec. Donec condimentum nulla odio, quis varius libero dignissim eget. Duis scelerisque purus id pretium viverra. Aliquam erat volutpat. Phasellus eleifend tincidunt porta. Nam facilisis hendrerit diam sit amet blandit. Quisque ac eros ultrices, tincidunt turpis in, auctor quam. Integer malesuada erat felis, non consectetur libero porta ac. Sed in nisi ut sem interdum vulputate in non nisi. Proin lorem eros, posuere vitae fringilla quis, dapibus id mi. Praesent commodo risus sit amet augue fermentum tincidunt ut vel orci. Mauris tincidunt mattis tellus ut malesuada. Donec ut nisi fermentum, suscipit risus at, pretium eros. Fusce fermentum consequat nibh a tincidunt. Etiam ornare massa et risus ornare maximus. Sed sollicitudin congue interdum fusce.`
	go func() {
		for _, v := range str {
			h.inputEventChan <- termbox.Event{Ch: v}
		}
		h.inputEventChan <- termbox.Event{Key: '\n'}
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
	assert.Equal(t, str+string('\n'), string(result))
	assert.Equal(t, 0, len(h.inputBuffer))
}
