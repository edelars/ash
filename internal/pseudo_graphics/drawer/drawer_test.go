package drawer

import (
	"testing"

	"github.com/nsf/termbox-go"
	"github.com/stretchr/testify/assert"
)

func TestDrawer_saveScreenState(t *testing.T) {
	h := Drawer{}
	f := func(x, y int) termbox.Cell {
		if x > 10 || y > 7 {
			t.Error("fail 10")
		}
		return termbox.Cell{Fg: 1, Bg: 3, Ch: 'i'}
	}
	h.saveScreenState(10, 7, f)

	assert.Equal(t, 10, len(h.screenState))
	for _, v := range h.screenState {
		assert.Equal(t, 7, len(v))
		for _, s := range v {
			assert.Equal(t, 'i', s.Ch)
			assert.Equal(t, termbox.Attribute(1), s.Fg)
			assert.Equal(t, termbox.Attribute(3), s.Bg)
		}
	}
}
