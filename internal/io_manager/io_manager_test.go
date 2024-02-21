package io_manager

import (
	"testing"

	"ash/internal/colors_adapter"
	"ash/internal/configuration"

	"github.com/nsf/termbox-go"
	"github.com/stretchr/testify/assert"
)

type pmImpl struct{}

func (pmimpl *pmImpl) DeleteLastSymbolFromCurrentBuffer() error {
	panic("not implemented") // TODO: Implement
}

func Test_inputManager_rollScreenUp(t *testing.T) {
	pm := pmImpl{}
	h := NewInputManager(&pm, configuration.CmdRemoveLeftSymbol, colors_adapter.NewColorsAdapter(configuration.Colors{}))

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
